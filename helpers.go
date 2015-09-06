package main

import (
	"errors"
	"fmt"
	"github.com/juju/deputy"
	"log"
	"os"
	"os/exec"
)

func validation(str ...interface{}) {
	log.SetPrefix("docker-build: ")
	log.Println(str...)
	os.Exit(1)
}

func internalError(err error) {
	log.Println(err)
	os.Exit(1)
}

func done(str ...interface{}) {
	if len(str) > 0 {
		fmt.Println(str...)
	}
	os.Exit(0)
}

type Variable struct {
	Cmd   string `yaml:"cmd,omitempty"`
	Value string `yaml:"value,omitempty"`
}

var (
	ErrCmdAndValue  = errors.New("A variable can not contain both a cmd and value entry")
	ErrNoCmdOrValue = errors.New("A variable must have a command or a value")
	ErrNoCommand    = errors.New("No command to execute")
)

func parseVariables(variables map[string]*Variable) (map[string]string, error) {
	parsed := map[string]string{}
	for name, v := range variables {
		if len(v.Cmd) > 0 && len(v.Value) > 0 {
			return nil, ErrCmdAndValue
		}

		switch {
		case len(v.Cmd) > 0:
			value, err := execute(v.Cmd)
			if err != nil {
				return nil, err
			}
			parsed[name] = value
		case len(v.Value) > 0:
			parsed[name] = v.Value
		default:
			return nil, ErrNoCmdOrValue
		}
	}
	return parsed, nil
}

func execute(cmd string) (string, error) {
	if len(cmd) == 0 {
		return "", ErrNoCommand
	}

	var value string
	d := deputy.Deputy{
		Errors: deputy.FromStderr,
		StdoutLog: func(b []byte) {
			value = string(b)
		},
	}
	if err := d.Run(exec.Command("sh", "-c", cmd)); err != nil {
		return "", err
	}
	return value, nil
}

func nodeFile(buildDir, node string) string {
	var nodeFile string
	if node == "" {
		nodeFile = fmt.Sprintf("%s/docker-build.yml", buildDir)
	} else {
		nodeFile = fmt.Sprintf("%s/docker-build-%s.yml", buildDir, node)
	}
	return nodeFile
}
