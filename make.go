package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type BuildConfig struct {
	Variables map[string]*Variable
	variables map[string]string
	Nodes     map[string]*NodeConfig
}

func make(build string) {
	filename := fmt.Sprintf("conf/%s.yml", build)
	dirName := fmt.Sprintf("builds/%s", build)

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		validation("Build", build, "does not exist")
	} else if err != nil {
		validation("Invalid build argument", err.Error())
	}

	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		validation("Build", build, "already exists. Run \"docker-build remove", build, "\" if first to remove it.")
	}

	config, err := readBuildConfig(filename)
	if err != nil {
		validation("Unable to read config file", err.Error())
	}

	if len(config.Nodes) == 0 {
		validation("Build must contain at least one node")
	}

	for name, node := range config.Nodes {
		if len(node.Containers) == 0 {
			validation("Node", name, "must have at least one container")
		}
	}

	var envB []byte
	if len(config.Variables) > 0 {
		var err error
		envB, err = config.buildEnv()
		if err != nil {
			validation("Unable to build env file", err.Error())
		}
	}

	nodeBs, err := config.buildNodeConfigs()
	if err != nil {
		validation("Unable to build node configs", err.Error())
	}

	if err := os.Mkdir(dirName, 0644); err != nil {
		validation("Unable to create build directory", err.Error())
	}

	envFile := fmt.Sprintf("%s/docker-build.env", dirName)
	if len(config.Variables) > 0 {
		if err := ioutil.WriteFile(envFile, envB, 0644); err != nil {
			validation("Unable to write docker-compose.yaml file:", err.Error())
		}
	}

	for name, b := range nodeBs {
		nodeFilename := nodeFile(dirName, name)
		if err := ioutil.WriteFile(nodeFilename, b, 0644); err != nil {
			validation("Unable to write ", nodeFilename, "file:", err.Error())
		}
	}

	done("Build", build, "created")
}

func readBuildConfig(configFile string) (*BuildConfig, error) {
	var config *BuildConfig

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

func (config *BuildConfig) parse() error {
	if len(config.Variables) > 0 {
		vars, err := parseVariables(config.Variables)
		if err != nil {
			return err
		}
		config.variables = vars
	}
	return nil
}

func (config *BuildConfig) buildEnv() ([]byte, error) {
	b, err := yaml.Marshal(config.variables)
	if err != nil {
		internalError(err)
	}
	return b, nil
}

func (config *BuildConfig) buildNodeConfigs() (map[string][]byte, error) {
	configs := map[string][]byte{}
	for name, node := range config.Nodes {
		b, err := yaml.Marshal(node)
		if err != nil {
			internalError(err)
		}
		configs[name] = b
	}
	return configs, nil
}
