package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

func getFileEnvMap() map[string]string {
	fileEnvMap := make(map[string]string)

	for _, env := range os.Environ() {
		nameVal := strings.SplitN(env, "=", 2)
		if strings.HasSuffix(nameVal[0], "_FILE") {
			fileEnvMap[nameVal[0]] = nameVal[1]
		}
	}

	return fileEnvMap
}

func readFileEnvs(fileEnvMap map[string]string) map[string]string {
	readFileEnvMap := make(map[string]string)

	for fileEnvName, filePath := range fileEnvMap {
		envValue, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatal(err)
		} else {
			envName := strings.TrimSuffix(fileEnvName, "_FILE")
			readFileEnvMap[envName] = string(envValue)
		}
	}

	return readFileEnvMap
}

func buildEnvs(envMap map[string]string) []string {
	var envs []string

	for name, value := range envMap {
		envs = append(envs, name+"="+value)
	}

	return envs
}

func pipeSignals(cmd *exec.Cmd, c chan os.Signal) {
	for sig := range c {
		cmd.Process.Signal(sig)
	}
}

func main() {
	fileEnvMap := getFileEnvMap()
	log.Println(fileEnvMap)

	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	cmd.Env = append(os.Environ(), buildEnvs(readFileEnvs(fileEnvMap))...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	c := make(chan os.Signal, 1)
	defer close(c)
	signal.Notify(c)
	go pipeSignals(cmd, c)

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
