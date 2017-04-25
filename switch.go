package main

import (
	"os"
	"fmt"
	"os/exec"
	"flag"
	"bufio"
)

const (
	headerRegexp = `^([\w-]+):\s*(.+)`
	authRegexp   = `^(.+):([^\s].+)`
)

var (
	gopath = os.Getenv("GOPATH")
	headerslice headerSlice

	f = flag.String("folder", "/usr/bin", "")
	v = flag.String("version", "", "")
	usage = `Usage: php-switch [--options]
		-f Folder where php versions are instaleld. By default it will be /usr/bin
		-v Version you want to switch to. Is Empty by default`
)

//region HeadSlice
type headerSlice []string
func (h *headerSlice) String() string {
	return fmt.Sprintf("%s", *h)
}
func (h *headerSlice) Set(value string) error {
	*h = append(*h, value)
	return nil
}
//endregion

func handleCmd (name string, args []string, successMessage string) {
	cmd := exec.Command(name, args...)
	reader, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(reader)
	go func() {
		fmt.Println(successMessage)
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

	fmt.Println("\r\n")
}

func main() {
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

	//region Current PHP version
	cmdCurrentVersionName := "php"
	cmdCurrentversionArgs := []string{"php -r \\@phpinfo\\(\\)\\; | grep 'PHP Version' -m 1"}

	handleCmd(cmdCurrentVersionName, cmdCurrentversionArgs, "Currently used version:")
	//endregion


	//region Locate PHP versions
	cmdLocateName := "locate"
	cmdLocateArgs := []string{"-r '^" + root + "/php[0-9]'"}

	handleCmd(cmdLocateName, cmdLocateArgs, "Available versions:")
	//endregion


}
