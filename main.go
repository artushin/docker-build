package main

import (
	"flag"
	"log"
	"os"
)

const helpText = `
	docker-build is a helper for generating and managing docker-compose files for multiple nodes.
	The command takes an action as its first argument and optional second and third arguments which modifies the action.

	Actions:
		"make": Generates build files to ./build/{build-name}. The second argument should be the build name.
		"remove": Removes the entire generated build. The second argument should be the build name.
		"checkout": If a build has been generated, the build's docker-build-{node-name}.yml is generated and placed in the current directory as docker-compose.yaml. The second argument should be the build name. The third argument should be the node name.
		"clean": Removes the docker-compose.yaml file from the current directory.
`

var (
	configFile string
)

func main() {
	arg := flag.Arg(0)
	switch arg {
	case "make":
		if build := flag.Arg(1); build != "" {
			make(build)
		}
	case "remove":
		if build := flag.Arg(1); build != "" {
			remove(build)
		}
	case "checkout":
		if build := flag.Arg(1); build != "" {
			checkout(build, flag.Arg(2))
		}
	case "clean":
		clean()
	}
	log.Println(helpText)
	os.Exit(1)
}
