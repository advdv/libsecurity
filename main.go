package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/advanderveer/docksec/twitter"
)

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
