package main

import (
	"log"

	"github.com/advanderveer/docksec/twitter"
)

func main() {
	tw, err := twitter.NewStream()
	if err != nil {
		log.Printf("Failed to create twitter stream: %s", err)
	}

	//listen for tweets
	log.Printf("Starting twitter stream...")
	ch := tw.Start()
	for {
		select {
		case <-tw.Quit():
			log.Println("Error in twitter stream, exiting!")
			return
		case msg := <-ch:
			log.Printf("Message: %T", msg)
		}
	}

}
