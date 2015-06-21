package main

import (
    "fmt"
    "strings"
    "os"
    "github.com/fsouza/go-dockerclient"
    "encoding/json"
)

type CheckInfect struct {
  images     []docker.APIImages     `json:"images,omitempty"`
  containers []docker.APIContainers `json:"containers,omitempty"`
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
  containers, _ := client.ListContainers(docker.ListContainersOptions{All: true})
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
    var client *docker.Client

    endpoint := os.Getenv("DOCKER_HOST")
    if ( endpoint != "") {
      if (strings.HasPrefix(endpoint, "tcp://")) {
        path := os.Getenv("DOCKER_CERT_PATH")
        if( path != "") {
          ca := fmt.Sprintf("%s/ca.pem", path)
          cert := fmt.Sprintf("%s/cert.pem", path)
          key := fmt.Sprintf("%s/key.pem", path)
          client, _ = docker.NewTLSClient(endpoint, cert, key, ca)
        } else {
          client, _ = docker.NewClient(endpoint)
        }
      } else {
        client, _ = docker.NewClient(endpoint)
      }
    } else {
      endpoint := "unix:///var/run/docker.sock"
      client, _ = docker.NewClient(endpoint)
    }

    infect_id := os.Getenv("INFECT_ID")
    infect := check(client,infect_id)
    fmt.Println(infect)
    data, _ := json.Marshal(infect)
    fmt.Println(string(data))
}
