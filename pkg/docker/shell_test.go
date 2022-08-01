package docker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseCmd(t *testing.T) {
	t.Run("should split cmd from params in cmd string", func(t *testing.T) {
		exe, _ := parseCmd("docker info")
		assert.Equal(t, "docker", exe)
	})
	t.Run("should append arg in cmd string at the start of params array", func(t *testing.T) {
		_, params := parseCmd("docker info", "hello")
		assert.Equal(t, []string{"info", "hello"}, params)
	})
}
