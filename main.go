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
		log.Println(ev)
	}

}
