package main

import (
	"os"
	"fmt"
	"log"
	"os/exec"
	"flag"
	"bufio"
)

const (
	headerRegexp = `^([\w-]+):\s*(.+)`
	authRegexp   = `^(.+):([^\s].+)`
)

type headerSlice []string

func (h *headerSlice) String() string {
	return fmt.Sprintf("%s", *h)
}

func (h *headerSlice) Set(value string) error {
	*h = append(*h, value)
	return nil
}

var (
	gopath = os.Getenv("GOPATH")
	headerslice headerSlice

	f = flag.String("folder", "/usr/bin", "")
)

var usage = `Usage: php-switch [--options] <version>

	-f Folder where php versions are instaleld. By default it will be /usr/bin
`

func visit(path string, f os.FileInfo, err error) error {
	fmt.Printf("Visited: %s\n", path)
	return nil
}

func main() {
	out, err := exec.Command("date").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The date is %s\n", out)

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintln(usage))
	}




	flag.Var(&headerslice, "H", "")

	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Println("No Arg !")
		flag.Usage()
		os.Exit(1)
	}

	root := *f

	version := flag.Args()[0]
	fmt.Println(version)

	cmdName := "locate"
	cmdArgs := []string{"-r '^" + root + "/php[0-9]'"}

	cmd := exec.Command(cmdName, cmdArgs...)
	cmdReader, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		fmt.Println("Installed versions \r\n %s")
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		os.Exit(1)
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		os.Exit(1)
	}
}
