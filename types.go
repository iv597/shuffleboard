package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os/exec"
	"time"
)

type taskCommandLine []string

func (s *taskCommandLine) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func (s *taskCommandLine) String() string {
	return ""
}

func (s *taskCommandLine) IsCumulative() bool {
	return true
}

func getTaskCmd(s kingpin.Settings) (target *[]string) {
	target = new([]string)
	s.SetValue((*taskCommandLine)(target))
	return
}

type TaskSwitchMethod int

func (t TaskSwitchMethod) String() string {
	return fmt.Sprintf("%d", int(t))
}

type Config struct {
	taskSwitch TaskSwitchMethod
	tasks      []TaskRunner
	http       HttpConfig
}

type TaskRunner struct {
	runs     int
	http     HttpConfig
	command  taskCommandLine
	instance *exec.Cmd
	minWait  time.Duration
	maxWait  time.Duration
}

type HttpConfig struct {
	address string
	port    int
}
