package main

import (
	"fmt"
	"os"
)

func remove(build string) {
	dirname := fmt.Sprintf("builds/%s", build)

	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		validation("Build", build, "does not exist")
	} else if err != nil {
		validation("Invalid build argument", err.Error())
	}

	if err := os.RemoveAll(dirname); err != nil {
		validation("Unable to remove", dirname, ":", err.Error())
	}

	done("Build removed")
}
