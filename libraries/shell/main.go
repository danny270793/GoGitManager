package shell

import (
	"bytes"
	"os/exec"
)

type Shell struct {
	ShellToUse string
}

func (s *Shell) Execute(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(s.ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func NewBashShell() Shell {
	s := Shell{ShellToUse: "bash"}
	return s
}
