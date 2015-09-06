package main

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"text/template"
)

type Container struct {
	Links     []string `yaml:"links,omitempty"`
	File      string   `yaml:"file,omitempty"`
	Variables []string `yaml:"variables,omitempty"`
}

type NodeConfig struct {
	Variables  map[string]*Variable `yaml:"variables,omitempty"`
	variables  map[string]string
	Containers map[string]*Container `yaml:"containers,omitempty"`
}

func checkout(build, node string) {
	buildDir := fmt.Sprintf("builds/%s", build)
	if _, err := os.Stat(buildDir); os.IsNotExist(err) {
		validation("Build", build, "does not exist")
	} else if err != nil {
		validation("Invalid build argument:", err.Error())
	}

	envFile := fmt.Sprintf("%s/docker-build.env", buildDir)
	_, err := os.Stat(envFile)
	if os.IsNotExist(err) {
		envFile = ""
	}

	nodeFilename := nodeFile(buildDir, node)
	if _, err := os.Stat(nodeFilename); os.IsNotExist(err) {
		validation("Node", node, "does not exist")
	} else if err != nil {
		validation("Invalid node argument:", err.Error())
	}

	config, err := readNodeConfig(nodeFilename)
	if err != nil {
		validation("Unable to read node config:", err.Error())
	}

	b, err := config.build(envFile)
	if err != nil {
		validation("Unable to build docker-compose.yml file:", err.Error())
	}

	if err := ioutil.WriteFile("docker-compose.yml", b, 0644); err != nil {
		validation("Unable to write docker-compose.yml file:", err.Error())
	}

	done("Wrote docker-compose.yml for build", build, "node", node)
}

func readNodeConfig(configFile string) (*NodeConfig, error) {
	var config *NodeConfig

	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}

	if err := config.parse(); err != nil {
		return nil, err
	}

	return config, nil
}

func (config *NodeConfig) parse() error {
	if len(config.Variables) > 0 {
		vars, err := parseVariables(config.Variables)
		if err != nil {
			return err
		}
		config.variables = vars
	}
	return nil
}

func (config *NodeConfig) build(envFile string) ([]byte, error) {
	b := bytes.NewBuffer(nil)

	for name, container := range config.Containers {
		f, err := ioutil.ReadFile(container.File)
		if err != nil {
			return nil, err
		}

		t, err := template.New(name).Parse(string(f))
		if err != nil {
			return nil, err
		}

		if err := t.Execute(b, config.variables); err != nil {
			return nil, err
		}

		if envFile != "" {
			b.WriteString("\n  env_file: ")
			b.WriteString(envFile)
		}

		if len(container.Links) > 0 {
			data := map[string]interface{}{
				"links": container.Links,
			}

			l, err := yaml.Marshal(data)
			if err != nil {
				return nil, err
			}
			l = bytes.Replace(l, []byte("\n"), []byte("\n   "), len(container.Links))

			if _, err := b.WriteString("\n  "); err != nil {
				return nil, err
			}
			if _, err := b.Write(l); err != nil {
				return nil, err
			}
		} else {
			if _, err := b.WriteString("\n"); err != nil {
				return nil, err
			}
		}
	}

	return b.Bytes(), nil
}
