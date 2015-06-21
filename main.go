package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"

	"github.com/advanderveer/docksec/twitter"
)

func RunCheckInfect(image string) error {
	cmd := exec.Command("./checkinfect/infect.sh", image)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

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
	log.Println("Running checkinfect.sh...")
	err = RunCheckInfect("img")
	if err != nil {
		log.Printf("Failed to run checkinfect.sh: %s", err)
	}
	log.Println("Ran succesfully!")

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
