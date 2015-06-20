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
	for {
		select {
		case <-tw.Quit():
			log.Println("Quit in twitter stream, exiting!")
			return
		case msg := <-tw.Events():
			log.Printf("Message: %T", msg)
		}
	}

}
