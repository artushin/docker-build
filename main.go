package main

import (
	"errors"
	"flag"
	"github.com/juju/deputy"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

var (
	configFile string
)

func main() {
	log.SetPrefix("docker-build: ")
	flag.StringVar(&configFile, "file", "docker-build.yaml", "Specify an alternate build file (default: docker-build.yml)")
	flag.StringVar(&configFile, "f", "docker-build.yaml", "Specify an alternate build file (default: docker-build.yml)")
	flag.Parse()
	config, err := readConfig()
	if err != nil {
		os.Exit(1)
	}

	log.Println(config.variables)

	// b, err := config.build()
	// if err != nil {
	// 	log.Println(err)
	// 	return nil, err
	// }

	// log.Println(string(b))
}

type Variable struct {
	Cmd   string
	Value string
}

type Brick struct {
	Links []string
	File  string
}

type Config struct {
	Variables map[string]*Variable
	variables map[string]string
	Bricks    map[string]*Brick
}

func readConfig() (*Config, error) {
	var config *Config

	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = yaml.Unmarshal(b, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if err := validateConfig(config); err != nil {
		log.Println(err)
		return nil, err
	}

	return config, nil
}

var (
	ErrCmdAndValue  = errors.New("A variable can not contain both a cmd and value entry")
	ErrNoCmdOrValue = errors.New("A variable must have a command or a value")
)

func validateConfig(config *Config) error {
	if len(config.Variables) > 0 {
		config.variables = map[string]string{}
		for name, v := range config.Variables {
			if len(v.Cmd) > 0 && len(v.Value) > 0 {
				return ErrCmdAndValue
			}

			switch {
			case len(v.Cmd) > 0:
				value, err := execute(v.Cmd)
				if err != nil {
					return err
				}
				config.variables[name] = value
			case len(v.Value) > 0:
				config.variables[name] = v.Value
			default:
				return ErrNoCmdOrValue
			}
		}
	}
	return nil
}

var ErrNoCommand = errors.New("No command to execute")

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

// func (config *Config) build() (string, error) {
// 	b := bytes.NewBuffer(nil)

// 	for name, brick := range config.Bricks {
// 		brick.File
// 	}
// }
