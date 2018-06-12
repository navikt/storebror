package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	flag "github.com/spf13/pflag"
)

var (
	repository *string = flag.String("repository", "", "foo")
)

func checkout(repository, destination string) error {
	cmd := exec.Command("git", "clone", "--depth", "1", repository, destination)
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf(string(bytes))
	}
	return err
}

func run() error {
	dir, err := ioutil.TempDir("", "storebror")
	log.Printf("Creating temporary directory at %s\n", dir)
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir) // clean up

	log.Printf("Checking out %s...\n", *repository)
	err = checkout(*repository, dir)

	return err
}

func main() {
	flag.Parse()
	err := run()
	if err == nil {
		return
	}
	fmt.Println(err)
	os.Exit(1)
}
