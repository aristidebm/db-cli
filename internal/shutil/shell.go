package shutil

import (
	"os"
	"os/exec"
	"syscall"
	"time"
	"io"
)

type cmdOption struct {
	stdin io.Reader
	stdout io.Writer
	stderr io.Reader
	cwd string
}

type cmdOptionFunc func (*cmdOption)

func WithStdin(stdin io.Reader) cmdOptionFunc {
	return func (o *cmdOption)  {
		o.stdin = stdin		
	}
}

func WithStdout(stdout io.Writer) cmdOptionFunc {
	return func (o *cmdOption)  {
		o.stdout = stdout 
	}
}

func WithStderr(stderr io.Reader) cmdOptionFunc {
	return func (o *cmdOption)  {
		o.stderr = stderr 
	}
}

func WithCwd(cwd string)  cmdOptionFunc {
	return func (o *cmdOption)  {
		o.cwd = cwd 
	}
}

func defaultCmdOption(o *cmdOption, options ...cmdOptionFunc) {
	for _, fn := range options {
		fn(o)	
	} 

	if o.stdin == nil {
		o.stdin = os.Stdin 
	}

	if o.stdout == nil {
		o.stdout = os.Stdout 
	}

	if o.stderr == nil {
		o.stderr = os.Stderr
	}
}

// Command execution utilities
func RunCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func RunInteractiveCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	
	// For interactive commands, we want to replace the current process
	// This mimics the behavior of exec in shell
	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			}
		}
		return err
	}
	return nil
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func CreateDirIfNotExists(dir string) error {
	if !FileExists(dir) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

func GetCurrentTimestamp() string {
	return time.Now().Format("20060102_150405")
}
