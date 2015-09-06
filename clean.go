package main

import (
	"log"
	"os"
)

const dockerComposeFile = "docker-compose.yml"

func clean() {
	if _, err := os.Stat(dockerComposeFile); os.IsNotExist(err) {
		log.Println(err)
		done()
	}

	if err := os.Remove(dockerComposeFile); err != nil {
		validation("Unable to remove", dockerComposeFile, ":", err.Error())
	}

	done()
}
