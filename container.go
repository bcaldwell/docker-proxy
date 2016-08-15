package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"text/template"
)

var tmpl *template.Template
var nginxConfLocation string
var nginxReload = true

func init() {
	var err error
	tmpl, err = template.New("nginx.tmpl").ParseFiles("./nginx.tmpl")

	handleError(err, true)

	nginxConfLocation = os.Getenv("NGINXCONFLOCATION")
	if nginxConfLocation == "" {
		nginxConfLocation = "/etc/nginx/conf.d"
	}

	val := os.Getenv("NGINXRELOAD")
	if val == "false" {
		nginxReload = false
	}
}

type container struct {
	ID       string
	IP       string
	Hostname string
	Port     int
}

func (c *container) SetIP(ip string) {
	c.IP = ip
}

func (c *container) WriteConfig() {
	fileName := c.fileName()
	fmt.Println("Writing conf to " + fileName)

	f, err := os.Create(fileName)
	handleError(err)
	err = tmpl.Execute(f, c)
	handleError(err)

	if nginxReload {
		err := exec.Command("service", "nginx", "reload").Run()
		handleError(err)
	}
}

func (c *container) DeleteConfig() {
	fileName := c.fileName()
	fmt.Println("deleting conf " + fileName)
	err := os.Remove(fileName)
	handleError(err)
}

func (c *container) fileName() string {
	return path.Join(nginxConfLocation, c.ID) + ".conf"
}
