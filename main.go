package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/creack/pty"
)

const USAGE = "  reple spawn {command}\n  reple eval < {path_to_file}\n  {stdin} | reple eval"
const EXAMPLES = "  reple spawn 'bash --noprofile --norc -i'\n  reple eval < test.py\n  echo 'ls -la' | reple eval"

var FIFO_PATH string = strings.Join([]string{os.TempDir(), "reple"}, string(os.PathSeparator))

func main() {
	switch os.Args[1] {
	case "spawn":
		checkArgsNumber(2)
		syscall.Mkfifo(FIFO_PATH, 0640)
		defer os.Remove(FIFO_PATH)
		spawn(parseSpawnCommand(os.Args[2]))
	case "eval":
		eval()
	default:
		usage()
		os.Exit(1)
	}
}

func spawn(command string, args []string) {
	cmd := exec.Command(command, args...)
	ptmx, err := pty.Start(cmd)
	processError(err)
	defer func() { _ = ptmx.Close() }()

	pipe, err := os.OpenFile(FIFO_PATH, os.O_RDWR, 0640)
	processError(err)
	defer pipe.Close()

	pipeReader := bufio.NewReader(pipe)

	go func() {
		_, _ = io.Copy(os.Stdout, ptmx)
	}()
	go func() {
		_, _ = io.Copy(ptmx, pipeReader)
	}()
	_, _ = io.Copy(ptmx, os.Stdin)
}

func eval() {
	f, err := os.OpenFile(FIFO_PATH, os.O_WRONLY, 0600)
	processError(err)
	defer f.Close()

	stdin, err := io.ReadAll(os.Stdin)
	processError(err)
	str := string(stdin) + "\n"

	_, err = f.WriteString(str)
	processError(err)
}

func usage() {
	fmt.Println("\nUSAGE:\n" + USAGE + "\n\nEXAMPLES:\n" + EXAMPLES + "\n")
}

func processError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func parseSpawnCommand(rawCommand string) (string, []string) {
	splited := strings.Split(rawCommand, " ")
	return splited[0], splited[1:]
}

func checkArgsNumber(num int) {
	if len(os.Args) < num + 1 {
		usage()
		os.Exit(1)
	}
}

