package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/advanderveer/docksec/twitter"
	"github.com/fsouza/go-dockerclient"
)

var dock *docker.Client

type CheckInfect struct {
	images     []docker.APIImages
	containers []docker.APIContainers
}

func scan(client *docker.Client, infect_id string) (*CheckInfect, error) {
	infect := &CheckInfect{[]docker.APIImages{}, []docker.APIContainers{}}
	imgs, err := client.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		return infect, err
	}

	for _, img := range imgs {
		himgs, _ := client.ImageHistory(img.ID)
		for _, himg := range himgs {
			if himg.ID == infect_id {
				infect.images = append(infect.images, img)
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
				infect.containers = append(infect.containers, container)
			}
		}
	}

	return infect, nil
}

func main() {
	tw, err := twitter.NewStream("advanderveer")
	if err != nil {
		log.Fatalf("Failed to create twitter stream: %s", err)
	}

	//handle signals
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		log.Println("Received interrupt signal, quitting twitter stream...")
		tw.Stop()
	}()

	//run checkinfect
	// log.Println("Running checkinfect.sh...")
	// err = RunCheckInfect("8c2e06607696bd4afb3d03b687e361cc43cf8ec1a4a725bc96e39f05ba97dd55")
	// if err != nil {
	// 	log.Printf("Failed to run checkinfect.sh: %s", err)
	// }
	// log.Println("Ran succesfully!")

	//
	client, err := docker.NewClient("unix:///var/run/docker.sock")
	if err != nil {
		log.Fatalf("Failed to connect to the docker daemon: %s", err)
	}

	// run check images
	log.Printf("Scanning Daemon...")
	res, err := scan(client, "8c2e06607696bd4afb3d03b687e361cc43cf8ec1a4a725bc96e39f05ba97dd55")
	if err != nil {
		log.Fatalf("Failed to scan host: %s", err)
	}
	log.Printf("Found vulnerabilities in %d images and %d containers!", len(res.images), len(res.containers))

	//listen for tweets
	log.Printf("Starting twitter stream...")
	evs := tw.Start()
	for ev := range evs {
		if ev.Type == twitter.EventNewVulnerability {

			//@todo handle checking of vulnerability
			log.Printf("Reported '%s' in '%s'...", ev.CVE, ev.Image)

			//require: hostname, image name and container id
			hostname := "myhostname"
			image := "image"
			cid := "container id"

			err := tw.ReplyVulnerable(ev.Tweet, hostname, image, cid)
			if err != nil {
				log.Printf("Error replying vulnerable: %s", err)
			}

			log.Println("Sent reply")

		} else if ev.Type == twitter.EventFixVulnerability {
			//@todo handle fixing vulnerability
			log.Printf("Fixing '%s' with image '%s'... ", ev.Selector, ev.Image)
		}
	}
}
