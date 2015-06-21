package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"

	"github.com/advanderveer/docksec/twitter"
	"github.com/fsouza/go-dockerclient"
)

var dock *docker.Client

func Scan(image string) ([]docker.APIContainers, []docker.APIImages, error) {
	cs := []docker.APIContainers{}
	imgs := []docker.APIImages{}

	return cs, imgs, nil

}

func RunCheckInfect(image string) error {
	buff := bytes.NewBuffer(nil)
	mw := io.MultiWriter(os.Stdout, buff)

	cmd := exec.Command("./checkinfect/infect.sh", image)
	cmd.Stderr = mw
	cmd.Stdout = mw

	log.Println("STDOUT:", buff.String())

	return cmd.Run()
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

	// run check images
	log.Printf("Scanning Daemon...")
	cs, imgs, err := Scan("8c2e06607696bd4afb3d03b687e361cc43cf8ec1a4a725bc96e39f05ba97dd55")
	if err != nil {
		log.Fatalf("Failed to scan host: %s", err)
	}
	log.Println("Found vulnerabilities in %s images and %d containers!", len(cs), len(imgs))

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
