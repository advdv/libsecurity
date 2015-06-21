package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	"github.com/advanderveer/docksec/twitter"
	"github.com/fsouza/go-dockerclient"
)

//
// Consider the example that busybox has a vulnerability
//
// 1. consumersa are running a countainer as such: `docker run -it -p 8080:80 jerbi/apache:1.0`
// 2. and have the docksec container running: `make build && make`
// 3. sycoso (re)tweets: 'CVE-2014-6271 in 9e1ed860cc088ae4b68ce28fb8888739652729e1107054f58dff90979f7dc935'
// 4. first time: running container should restart with latest pulled version

var client *docker.Client

const (
	FILENAME = "images-cve.dat"
	FORMAT   = "%s@%s\r\n"
)

func markAsFixed(imageName string, cveName string) error {
	f, err := os.OpenFile(FILENAME, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		return err
	}
	defer f.Close()
	imageAndCve := fmt.Sprintf(FORMAT, cveName, imageName)

	_, err = f.WriteString(imageAndCve)
	return err
}

func alreadyFixed(imageName string, cveName string) (bool, error) {
	imageAndCve := fmt.Sprintf(FORMAT, cveName, imageName)

	f, err := os.OpenFile(FILENAME, os.O_RDONLY|os.O_RDWR|os.O_CREATE, 0)
	if err != nil {
		fmt.Printf("failed opening %s", FILENAME)
		return false, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	found := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == strings.TrimSpace(imageAndCve) {
			fmt.Printf("found!")
			found = true
			break
		}
	}
	return found, nil
}

func fix(vul *twitter.Vulnerable) error {
	auth, err := docker.NewAuthConfigurations(strings.NewReader("{}"))
	if err != nil {
		return err
	}

	//fix by pulling latest and marking as such afterwards
	for _, img := range vul.Images {
		res, err := alreadyFixed(img.ID, vul.CVE)
		if err != nil {
			log.Printf("Error checking if the image %s is already fixed: %s", img.ID, err)
		} else if res {
			log.Printf("Already pulled latest tags for image '%s': %s", img.ID, err)
			continue
		}

		for _, tag := range img.RepoTags {
			//pull each latest for each tag
			repo := tag[:strings.Index(tag, ":")]
			opts := docker.PullImageOptions{
				Repository: repo,
				Tag:        "latest",
			}

			log.Printf("Pulling '%s:latest'...", repo)
			err := client.PullImage(opts, auth.Configs["https://index.docker.io/v1/"])
			if err != nil {
				log.Printf("Error pulling latest image: %s", err)
			}
			log.Printf("Done")
		}

		err = markAsFixed(img.ID, vul.CVE)
		if err != nil {
			log.Printf("Error marking img '%s' and cve '%s' as fixed: %s", img.ID, vul.CVE, err)
		}
	}

	//restart each container with its newly pulled image and existing configs
	//@TODO handle container names?
	for _, apic := range vul.Containers {

		c, err := client.InspectContainer(apic.ID)
		if err != nil {
			log.Printf("Error inspecting container '%s': %s", apic.ID, err)
		}

		hostConfig := c.HostConfig

		if !c.State.Running {
			continue
		}

		//@todo we assume we are looking to run the latest
		c.Config.Image = c.Config.Image[:strings.Index(c.Config.Image, ":")] + ":latest"

		newc, err := client.CreateContainer(docker.CreateContainerOptions{
			Config: c.Config,
		})

		if err != nil {
			log.Printf("Error creating new container of '%s': %s", c.ID, err)
			break
		}

		//@todo this port change is actually domain logic
		hostConfig.PortBindings["80/tcp"] = []docker.PortBinding{
			docker.PortBinding{
				HostPort: "8081",
			},
		}

		err = client.StartContainer(newc.ID, hostConfig)
		if err != nil {
			log.Printf("Error start new container '%s': %s", newc.ID, err)
		}

		log.Printf("Created and started container of '%s': '%s', stopping old...", c.ID, newc.ID)

		err = client.StopContainer(c.ID, 1)
		if err != nil {
			log.Printf("Failed to stop old container '%s': %s", c.ID, err)
		}
	}

	return nil
}

func scan(infect_id string) (*twitter.Vulnerable, error) {
	infect := &twitter.Vulnerable{"", []docker.APIContainers{}, []docker.APIImages{}}
	imgs, err := client.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		return infect, err
	}

	for _, img := range imgs {
		himgs, _ := client.ImageHistory(img.ID)
		for _, himg := range himgs {
			if himg.ID == infect_id {
				infect.Images = append(infect.Images, img)
			}
		}
	}
	containers, err := client.ListContainers(docker.ListContainersOptions{All: false})
	if err != nil {
		return infect, err
	}

	for _, container := range containers {
		himgs, _ := client.ImageHistory(container.Image)
		for _, himg := range himgs {
			if himg.ID == infect_id {
				infect.Containers = append(infect.Containers, container)
			}
		}
	}

	return infect, nil
}

func main() {
	f, err := os.Create(FILENAME)
	if err != nil {
		log.Fatalf("failed to create file: %s", err)
	}
	f.Close()

	client, err = docker.NewClient("unix:///var/run/docker.sock")
	if err != nil {
		log.Fatalf("Failed to connect to the docker daemon: %s", err)
	}

	tw, err := twitter.NewStream("advanderveer")
	if err != nil {
		if apiErr, ok := err.(*anaconda.ApiError); ok {
			log.Printf("Api Error, headers: %s", apiErr.Header)

		}

		log.Fatalf("Failed to create twitter stream: %s", err)
	}

	defer tw.Close()

	//handle signals
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		log.Println("Received interrupt signal, quitting twitter stream...")
		tw.Stop()
	}()

	//listen for tweets
	log.Printf("Starting twitter stream...")
	evs := tw.Start()
	for ev := range evs {
		if ev.Type == twitter.EventNewVulnerability {

			// run check images, eg 8c2e06607696bd4afb3d03b687e361cc43cf8ec1a4a725bc96e39f05ba97dd55
			log.Printf("Scanning Daemon for image '%s' (%s)...", ev.Image, ev.CVE)
			res, err := scan(ev.Image)
			if err != nil {
				log.Fatalf("Failed to scan host: %s", err)
			}
			log.Printf("Done, Found vulnerabilities in %d images and %d containers!", len(res.Images), len(res.Containers))

			//require: hostname, image name and container id
			hostname, err := os.Hostname()
			if err != nil {
				log.Printf("Failed to determine hostname: %s", err)
			}

			if len(res.Images) > 0 || len(res.Containers) > 0 {
				res.CVE = ev.CVE
				err = tw.ReplyVulnerable(ev, hostname, res)
				if err != nil {
					log.Printf("Error replying vulnerable: %s", err)
				}
			}

			log.Println("Sent reply")

		} else if ev.Type == twitter.EventFixVulnerability {
			vul := ev.Vulnerable
			err := fix(vul)
			if err != nil {
				log.Printf("Failed to fix cve '%s' for %d images and %d containers", vul.CVE, len(vul.Images), len(vul.Containers))
			}

			log.Printf("Fixed cve '%s' for %d images and %d containers", vul.CVE, len(vul.Images), len(vul.Containers))
		}
	}
}
