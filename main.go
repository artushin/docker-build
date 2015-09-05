package main

import (
	"flag"
	"log"
	"os"
)

const helpText = `Usage: docker-build [COMMAND] [arg...]

docker-build is a helper for generating and managing docker-compose files for multiple nodes. It has NO association with Docker, Inc. and is provided open-source and as is.

Version: 0.1

Author:
  Alex Artushin - <https://github.com/artushin>

Commands:
  make			Generates build files to ./build/{build-name}
  rm			Removes the generated build
  checkout		Checkout a specified node in a build as docker-compose.yaml to the current directory
  clean		 	Removes the docker-compose.yaml file from the current directory`

var (
	configFile string
)

func main() {
	log.SetFlags(0)
	flag.Parse()

	arg := flag.Arg(0)
	switch arg {
	case "make":
		if build := flag.Arg(1); build != "" {
			make(build)
		}
	case "rm":
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
