package docker

import (
	"os/exec"
	"regexp"
	"strings"
)

type shell interface {
	runShellCmd(cmd string, args ...string) (string, error)
}

type baseShell struct{}

func (sh *baseShell) runShellCmd(cmd string, args ...string) (string, error) {
	exe, params := parseCmd(cmd, args...)
	out, err := exec.Command(exe, params...).CombinedOutput()
	return string(out), err
}

func parseCmd(cmd string, args ...string) (string, []string) {
	re := regexp.MustCompile(`[\r\t\n\f ]+`)
	a := strings.Split(re.ReplaceAllString(cmd, " "), " ")

	params := args
	if len(a) > 1 {
		params = append(a[1:], args...)
	}
	exe := a[0]

	return exe, params
}
