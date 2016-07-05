package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/openshift/origin/third_party/github.com/docker/libcompose/project"
)

// CheckIfFileExists checks if all files in array exists and that they are files (not directories).
func CheckIfFileExists(files []string) error {
	for _, filename := range files {
		fileInfo, err := os.Stat(filename)
		if err != nil {
			return fmt.Errorf("file %q not found", filename)
		}
		if fileInfo.IsDir() {
			return fmt.Errorf("%q is a directory", filename)
		}
	}
	return nil
}

// getInput is a generic function that takes a prompt string and scans a line
// and always returns string, caller has to parse the value to the required type
func getInput(prompt string) (response string) {
	fmt.Fprintf(os.Stdout, prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	response = scanner.Text()
	return
}

// AskForData loops on the parsed docker-compose file's config and the data that is missing will be queried
// interactively to user
func AskForData(config map[string]*project.ServiceConfig) {
	for svc := range config {

		// If ports not specified
		if len(config[svc].Ports) == 0 {
			msg := fmt.Sprintf("[%s] No ports defined to send traffic to, no provider service will be created. Do you want to create a service? y/[n]: ", svc)
			if resp := getInput(msg); resp == "y" {
				msg = fmt.Sprintf("[%[1]s] Note: Ports should of the form '8080' or '13306:3306' or '18888:8888 3306'\n[%[1]s] Enter ports : ", svc)

				// user can provide multiple ports separated by space
				// for e.g. "3306 18080:8080 7878"
				for _, port := range strings.Split(getInput(msg), " ") {
					config[svc].Ports = append(config[svc].Ports, port)
				}
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
