package main

import (
	"fmt"
	"log"

	"github.com/fsouza/go-dockerclient"
)

var client *docker.Client

func main() {

	fmt.Println("Starting devctl-proxy")

	var err error
	endpoint := "unix:///var/run/docker.sock"
	client, err = docker.NewClient(endpoint)
	handleError(err, true)

	containers := make(map[string]*container)

	activeContainers, _ := client.ListContainers(docker.ListContainersOptions{})
	for _, con := range activeContainers {
		id := con.ID
		container, _ := client.InspectContainer(id)
		labels := container.Config.Labels
		if _, ok := labels["devctl"]; ok {
			addContainer(containers, id, labels)
			containers[id].SetIP(getContainerIP(id))
			containers[id].WriteConfig()
		}
	}

	eventListener := make(chan *docker.APIEvents)

	client.AddEventListener(eventListener)

	for {
		select {
		case event := <-eventListener:
			if event != nil {
				if event.Action == "create" {
					id := event.Actor.ID
					attributes := event.Actor.Attributes
					addContainer(containers, id, attributes)
					fmt.Println(containers)
				}
				if event.Action == "connect" {
					// event.Actor.ID is network id
					attributes := event.Actor.Attributes
					id := attributes["container"]
					if c, ok := containers[id]; ok {
						c.SetIP(getContainerIP(id))
						c.WriteConfig()
					}
				}
				if event.Action == "disconnect" {
					attributes := event.Actor.Attributes
					id := attributes["container"]
					if c, ok := containers[id]; ok {
						c.DeleteConfig()
						delete(containers, id)
					}
				}
			}
		}
	}
}

func getContainerIP(id string) (ip string) {
	container, _ := client.InspectContainer(id)
	return container.NetworkSettings.IPAddress
}

func handleError(err error, fatal ...bool) {
	if err != nil {
		if len(fatal) > 0 && fatal[0] {
			log.Fatal(err)
		} else {
			log.Print(err)
		}
	}
}

func addContainer(containers map[string]*container, id string, labels map[string]string) {
	if hostname, ok := labels["devctl"]; ok {
		port := "80"
		if devctlPort, ok := labels["devctl-port"]; ok {
			port = devctlPort
		}
		hostname += ".devctl"
		containers[id] = &container{
			ID:       id,
			Hostname: hostname,
			Port:     port,
		}
	}
}
