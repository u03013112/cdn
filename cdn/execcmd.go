package cdn

import (
	"bytes"
	"errors"
	// "log"
	"os/exec"
	"time"
)

var (
	Timeout    = 3 * time.Second
	ErrTimeout = errors.New("Command timed out.")
)

type Invoker interface {
	Command(string, ...string) ([]byte, error)
}

type Invoke struct{}

func (i Invoke) Command(name string, arg ...string) ([]byte, error) {
	cmd := exec.Command(name, arg...)
	return CombinedOutputTimeout(cmd, Timeout)
}

func CombinedOutputTimeout(c *exec.Cmd, timeout time.Duration) ([]byte, error) {
	var b bytes.Buffer
	c.Stdout = &b
	c.Stderr = &b
	if err := c.Start(); err != nil {
		return nil, err
	}
	err := WaitTimeout(c, timeout)
	return b.Bytes(), err
}

func WaitTimeout(c *exec.Cmd, timeout time.Duration) error {
	timer := time.NewTimer(timeout)
	done := make(chan error)
	go func() { done <- c.Wait() }()
	select {
	case err := <-done:
		timer.Stop()
		return err
	case <-timer.C:
		if err := c.Process.Kill(); err != nil {
			log.Errorf("FATAL error killing process: %s", err)
			return err
		}
		<-done
		return ErrTimeout
	}
}
