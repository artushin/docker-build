package main

import (
	"os"
)

const dockerComposeFile = "docker-compose.yml"

func clean() {
	if _, err := os.Stat(dockerComposeFile); os.IsNotExist(err) {
		done()
	}

	if err := os.Remove(dockerComposeFile); err != nil {
		validation("Unable to remove", dockerComposeFile, ":", err.Error())
	}

	done()
}
