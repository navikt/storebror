package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	flag "github.com/spf13/pflag"
)

func checkout(repository, destination string) error {
	cmd := exec.Command("git", "clone", "--depth", "1", repository, destination)
	log.Printf("Running '%s'...\n", strings.Join(cmd.Args, " "))
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		lines := strings.Split(string(bytes), "\n")
		for _, line := range lines {
			log.Printf("git: %s", line)
		}
	}
	return err
}

func run(repository string) error {
	if len(repository) == 0 {
		return fmt.Errorf("missing repository URL.")
	}

	dir, err := ioutil.TempDir("", "storebror")
	if err != nil {
		return fmt.Errorf("while creating temporary directory at '%s': %s", dir, err)
	}
	defer os.RemoveAll(dir) // clean up

	log.Printf("Checking out '%s' to '%s'...\n", repository, dir)
	err = checkout(repository, dir)
	if err != nil {
		return fmt.Errorf("while checking out repository: %s", err)
	}

	return nil
}

func main() {
	flag.Parse()
	repository := flag.Arg(0)
	err := run(repository)
	if err == nil {
		return
	}
	log.Printf("Error: %s\n", err)
	os.Exit(1)
}
