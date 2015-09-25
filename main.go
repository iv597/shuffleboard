package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/alecthomas/kingpin.v2"
	//"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	taskPortsRaw := kingpin.Flag("innerPorts", taskPortsHelp).Short('P').String()

	task := getTaskCmd(kingpin.Arg("command", "task to shuffle"))

	kingpin.Version("0.0.0.0")
	kingpin.Parse()

	taskPorts := []int{}

	taskPortsSlice := strings.Split(*taskPortsRaw, ",")

	for _, port := range taskPortsSlice {
		if len(port) >= 1 {
			newVal, err := strconv.Atoi(port)

			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}

			taskPorts = append(taskPorts, newVal)
		}
	}

	fmt.Printf("%v%v%v%v  %v\n", *runCount, *bindTo, *portNum, taskPorts, *task)

	r := httprouter.New()
	r.GET("/*path", Shuffler)
	r.HEAD("/*path", Shuffler)
	r.OPTIONS("/*path", Shuffler)
	r.POST("/*path", Shuffler)
	r.PUT("/*path", Shuffler)
	r.PATCH("/*path", Shuffler)
	r.DELETE("/*path", Shuffler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", *bindTo, *portNum), r))
}

func Shuffler(_ http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	fmt.Printf("[child 1, 200ms] %v\n", p.ByName("path"))
}
