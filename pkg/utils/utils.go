package utils

import (
	"fmt"
	"os"

	"github.com/openshift/origin/third_party/github.com/docker/libcompose/project"
)

func getInput(prompt string) (response string) {
	fmt.Fprintf(os.Stdout, prompt)
	fmt.Scanf("%s", &response)
	return
}

func AskForData(config map[string]*project.ServiceConfig) {
	for svc := range config {

		// If ports not specified
		if len(config[svc].Ports) == 0 {
			msg := fmt.Sprintf("[%s] No ports defined to send traffic to, no provider service will be created. Do you want to create a service? y/[n]: ", svc)
			if resp := getInput(msg); resp == "y" {
				msg = fmt.Sprintf("[%s] Enter ports: ", svc)
				config[svc].Ports = append(config[svc].Ports, getInput(msg))
			}
		}

		// handle build
		if len(config[svc].Build) != 0 {
			msg := fmt.Sprintf("[%s] Do you want to use %q's git origin as source? [y]/n: ", svc, config[svc].Build)
			if resp := getInput(msg); resp == "n" {
				msg = fmt.Sprintf("[%s] Enter application image name: ", svc)
				config[svc].Image = getInput(msg)
				config[svc].Build = ""
			}
		}
	}
}
