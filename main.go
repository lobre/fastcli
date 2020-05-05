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
		fmt.Println("httpie is not installed")
		os.Exit(1)
	}

	userHome, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	confPath := path.Join(userHome, configFile)
	if !fileExists(confPath) {
		fmt.Printf("no config file at %s\n", confPath)
		os.Exit(1)
	}

	f, err := ioutil.ReadFile(confPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var conf Config
	if err := json.Unmarshal(f, &conf); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var envName string
	var verbose bool
	var pretty string
	flag.StringVar(&envName, "e", "", "environment")
	flag.BoolVar(&verbose, "v", false, "verbose")
	flag.StringVar(&pretty, "p", "all", "pretty")
	flag.Parse()

	if envName == "" {
		flag.Usage()
		os.Exit(1)
	}

	env, err := conf.getEnv(envName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(flag.Args()) == 0 {
		fmt.Println("missing httpie arguments")
		os.Exit(1)
	}

	var args []string
	if verbose {
		args = append(args, "-v")
	}
	args = append(args, "--ignore-stdin")
	args = append(args, fmt.Sprintf("--pretty=%s", pretty))
	args = append(args, fmt.Sprintf("%s/service/%s/%s", fastlyURL, env.ID, flag.Args()[0]))
	args = append(args, flag.Args()[1:]...)
	args = append(args, fmt.Sprintf("Fastly-Key:%s", env.Token))
	cmd := exec.Command("http", args...)
	fmt.Printf("executing: %s\n", cmd.String())

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	fmt.Println(string(output))
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
