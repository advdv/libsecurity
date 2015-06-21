package main

import (
    "fmt"
    "os"
    "github.com/fsouza/go-dockerclient"
)

type CheckInfect struct {
	images []docker.APIImages
  containers []docker.APIContainers
}

func check (client *docker.Client, infect_id string ) *CheckInfect {
  infect := &CheckInfect{[]docker.APIImages{},[]docker.APIContainers{}}
  imgs, _ := client.ListImages(docker.ListImagesOptions{All: false})
  for _, img := range imgs {
    himgs, _ := client.ImageHistory(img.ID)
    for _, himg := range himgs {
      if himg.ID == infect_id {
        infect.images = append(infect.images,img)
      }
    }
  }
  containers, _ := client.ListContainers(docker.ListContainersOptions{All: false})
  for _, container := range containers {
    himgs, _ := client.ImageHistory(container.Image)
    for _, himg := range himgs {
      if himg.ID == infect_id {
        infect.containers = append(infect.containers,container)
      }
    }
  }

  return infect
}

func main() {
    endpoint := os.Getenv("DOCKER_HOST")
    path := os.Getenv("DOCKER_CERT_PATH")
    ca := fmt.Sprintf("%s/ca.pem", path)
    cert := fmt.Sprintf("%s/cert.pem", path)
    key := fmt.Sprintf("%s/key.pem", path)
    client, _ := docker.NewTLSClient(endpoint, cert, key, ca)

    infect_id := os.Getenv("INFECT_ID")
    fmt.Println(check(client,infect_id))
}
