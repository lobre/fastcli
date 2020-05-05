package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

const (
	configFile = ".fastcli"
	fastlyURL  = "https://api.fastly.com"
)

type Env struct {
	Name  string
	ID    string
	Token string
}

type Config []Env

func (c Config) getEnv(name string) (Env, error) {
	for _, env := range c {
		if name == env.Name {
			return env, nil
		}
	}
	return Env{}, errors.New("environment not found")
}

func main() {
	if !cmdExists("http") {
		fmt.Print("httpie is not installed")
		os.Exit(1)
	}

	userHome, err := os.UserHomeDir()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	confPath := path.Join(userHome, configFile)
	if !fileExists(confPath) {
		fmt.Printf("no config file at %s", confPath)
		os.Exit(1)
	}

	f, err := ioutil.ReadFile(confPath)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	var conf Config
	if err := json.Unmarshal(f, &conf); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	var envName string
	flag.StringVar(&envName, "e", "", "environment")
	flag.Parse()

	if envName == "" {
		flag.Usage()
		os.Exit(1)
	}

	env, err := conf.getEnv(envName)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	if len(flag.Args()) == 0 {
		fmt.Println("missing httpie arguments")
		os.Exit(1)
	}

	var args []string
	args = append(args, fmt.Sprintf("%s/service/%s/%s", fastlyURL, env.ID, flag.Args()[0]))
	args = append(args, flag.Args()[1:]...)
	args = append(args, fmt.Sprintf("Fastly-Key:%s", env.Token))
	cmd := exec.Command("http", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	fmt.Print(string(output))
}

func cmdExists(cmd string) bool {
	if _, err := exec.LookPath(cmd); err != nil {
		return false
	}
	return true
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
