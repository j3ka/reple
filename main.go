package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/creack/pty"
	"golang.org/x/term"
)

const USAGE = "  reple spawn {command}\n  reple eval < {path_to_file}\n  {stdin} | reple eval"
const EXAMPLES = "  reple spawn 'bash --noprofile --norc -i'\n  reple eval < test.py\n  echo 'ls -la' | reple eval"

var FIFO_PATH string = strings.Join([]string{os.TempDir(), "reple"}, string(os.PathSeparator))

func main() {
	command := ""
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
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

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }()

	// Get current terminal size
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		// Fallback to reasonable defaults if we can't get terminal size
		width, height = 80, 24
	}

	cmd := exec.Command(command, args...)

	// Create PTY with proper terminal dimensions
	winsize := &pty.Winsize{
		Rows: uint16(height),
		Cols: uint16(width),
		X:    0,
		Y:    0,
	}

	ptmx, err := pty.StartWithSize(cmd, winsize)
	processError(err)
	defer func() { _ = ptmx.Close() }()

	pipe, err := os.OpenFile(FIFO_PATH, os.O_RDWR, 0640)
	processError(err)
	defer pipe.Close()

	pipeReader := bufio.NewReader(pipe)

	// Handle terminal resize signals
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				log.Printf("error resizing pty: %s", err)
			}
		}
	}()

	go func() {
		_, err = io.Copy(os.Stdout, ptmx)
		if err != nil {
			term.Restore(int(os.Stdin.Fd()), oldState)
			ptmx.Close()
			os.Exit(0)
		}
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
	if len(os.Args) < num+1 {
		usage()
		os.Exit(1)
	}
}
