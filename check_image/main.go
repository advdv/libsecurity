package main

import (
	"fmt"
	"os"

	"github.com/fsouza/go-dockerclient"
)

func main() {
	args := os.Args[1:]

//endpoint := os.Getenv("DOCKER_HOST")
        endpoint := "unix:///var/run/docker.sock"
	//path := os.Getenv("DOCKER_CERT_PATH")
	//ca := fmt.Sprintf("%s/ca.pem", path)
	//cert := fmt.Sprintf("%s/cert.pem", path)
	//key := fmt.Sprintf("%s/key.pem", path)
	//client, _ := docker.NewTLSClient(endpoint, cert, key, ca)
        client, _ := docker.NewClient(endpoint)

	imgs, _ := client.ListImages(docker.ListImagesOptions{All: false})

	for _, img := range imgs {
		for i := range args {

			config := &docker.Config{Image: img.ID, Entrypoint:[]string{"/bin/sh", "-c"}, Cmd: []string{"/script/shell.sh"}}

			// Create the container from an id, specify the command, mount CVE volume
			container, err := client.CreateContainer(docker.CreateContainerOptions{Config: config})
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// Get the working directory
			pwd, err := os.Getwd()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			binds := []string{pwd + "/cve/" + args[i] + "/:/script"}

			hostConfig := &docker.HostConfig{Binds: binds}
			err = client.StartContainer(container.ID, hostConfig)

			code, err := client.WaitContainer(container.ID)

			if code > 0 {
				fmt.Println(args[i] + " in " + img.ID)
			}

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}
}
