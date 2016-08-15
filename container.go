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
	tmpl, err = template.New("nginx.tmpl").Parse(`
		upstream {{.Hostname}} {
			server {{.IP}}:{{.Port}};
		}

		server {
			listen 80;
			server_name {{.Hostname}};

			access_log /var/log/nginx/{{.Hostname}}.access.json.log;
			error_log /var/log/nginx/{{.Hostname}}.access.log;

			# include /etc/nginx/ssl-public.conf;

			# require headers for http proxy
			proxy_set_header Client-IP         $remote_addr;
			proxy_set_header X-Real-IP         $remote_addr;
			proxy_set_header X-Forwarded-For   $remote_addr;
			proxy_set_header Host              $http_host;
			proxy_set_header X-Forwarded-Proto $scheme;
			proxy_set_header X-Forwarded-Port  $server_port;
			proxy_set_header Upgrade           $http_upgrade;
			proxy_set_header Connection        $http_connection;

			proxy_http_version 1.1;
			proxy_redirect off;
			proxy_next_upstream off;
			proxy_read_timeout 100s;

			# error_page 502 /502.devctl.html;
			# location /502.devctl.html {
			#  return 502 'Nothing is running at port {{.Port}} on your host. For this to work, you need to check your the relevant project and start the corresponding server on OS X.';
			# }

			location / {
				proxy_pass http://{{.Hostname}};
			}
		}
	`)

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
		err := exec.Command("nginx", "-s", "reload").Run()
		handleError(err)
	}
}

func (c *container) DeleteConfig() {
	fileName := c.fileName()
	fmt.Println("deleting conf " + fileName)
	err := os.Remove(fileName)
	handleError(err)

	if nginxReload {
		err := exec.Command("nginx", "-s", "reload").Run()
		handleError(err)
	}
}

func (c *container) fileName() string {
	return path.Join(nginxConfLocation, c.ID) + ".conf"
}
