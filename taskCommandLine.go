package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
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
