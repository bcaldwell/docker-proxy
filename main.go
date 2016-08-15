package main

import (
	"log"

	"github.com/fsouza/go-dockerclient"
)

var client *docker.Client

func main() {

	var err error
	endpoint := "unix:///var/run/docker.sock"
	client, err = docker.NewClient(endpoint)
	handleError(err, true)

	containers := make(map[string]*container)

	eventListener := make(chan *docker.APIEvents)

	client.AddEventListener(eventListener)

	for {
		select {
		case event := <-eventListener:
			if event != nil {
				if event.Action == "create" {
					id := event.Actor.ID
					attributes := event.Actor.Attributes
					if hostname, ok := attributes["devctl"]; ok {
						hostname += ".devctl"
						containers[id] = &container{
							ID:       id,
							Hostname: hostname,
							Port:     80,
						}
					}
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
