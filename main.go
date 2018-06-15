package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/navikt/storebror/scanner"
	flag "github.com/spf13/pflag"
)

type Author struct {
	Name  string
	Email string
}

var (
	noCleanTemporaryDirectory *bool   = flag.Bool("no-clean", false, "Clean up temporary directories on finish")
	gitConfigName             *string = flag.String("name", "Storebror", "Name of Git author")
	gitConfigEmail            *string = flag.String("email", "storebror@example.com", "Email of Git author")
)

func (a Author) Format() string {
	return fmt.Sprintf("%s <%s>", a.Name, a.Email)
}

func run_dmc(command ...string) error {
	cmd := exec.Command(command[0], command[1:]...)
	log.Printf("Running '%s'...\n", strings.Join(cmd.Args, " "))
	bytes, err := cmd.CombinedOutput()
	output := string(bytes)
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		log.Printf("%s: %s", command[0], line)
	}
	if err == nil {
		log.Printf("Exit code 0\n")
	}
	return err
}

func checkout(repository, destination string) error {
	return run_dmc("git", "clone", "--depth", "1", repository, destination)
}

func commit(workdir string, author Author, message string) error {
	os.Chdir(workdir)
	err := run_dmc("git", "add", "--all")
	if err != nil {
		return err
	}
	return run_dmc("git", "commit", "--author", author.Format(), "--message", message)
}

func run(repository string) error {
	if len(repository) == 0 {
		return fmt.Errorf("missing repository URL.")
	}

	dir, err := ioutil.TempDir("", "storebror")
	if err != nil {
		return fmt.Errorf("while creating temporary directory at '%s': %s", dir, err)
	}
	if !*noCleanTemporaryDirectory {
		defer os.RemoveAll(dir)
	}

	log.Printf("Checking out '%s' to '%s'...\n", repository, dir)
	err = checkout(repository, dir)
	if err != nil {
		return fmt.Errorf("while checking out repository: %s", err)
	}

	os.Chdir(dir)
	results := scanner.Process(dir)
	for _, result := range results {
		if result.Err != nil {
			return fmt.Errorf("while processing repository: %s", err)
		}
	}

	if len(results) == 0 {
		return nil
	}

	run_dmc("git", "status")
	run_dmc("git", "diff")
	commit(dir, Author{Name: *gitConfigName, Email: *gitConfigEmail}, results.Description())

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
