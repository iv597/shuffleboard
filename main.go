package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	//"io"
	//"net/http"
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

func main() {
	runCountHelp := "number of parallel executions - if your application is asynchronous, the default of 1 is safe"
	runCount := kingpin.Flag("count", runCountHelp).Default("1").Short('c').Int()
	bindTo := kingpin.Flag("bind", "bind address (IP/hostname)").Default("localhost").Short('b').String()
	portNum := kingpin.Flag("port", "port to listen on").Default("8005").Short('p').Int()
	taskPortsHelp := "comma-separated list (length of `count`) of ports to use for spawned processes"
	taskPorts := kingpin.Flag("innerPorts", taskPortsHelp).Short('P').String()

	task := getTaskCmd(kingpin.Arg("command", "task to shuffle"))

	kingpin.Version("0.0.0")
	kingpin.Parse()

	fmt.Printf("%v%v%v%v  %v", *runCount, *bindTo, *portNum, *taskPorts, *task)
}
