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
	"time"
)

const (
	TSM_SEQUENTIAL TaskSwitchMethod = iota
	TSM_RANDOMIZED
)

var config Config

func main() {
	runCountHelp := "number of parallel executions - if your application is asynchronous, the default of 1 is safe"
	runCount := kingpin.Flag("count", runCountHelp).Default("1").Short('c').Int()

	bindTo := kingpin.Flag("bind", "bind address (IP/hostname)").Default("localhost").Short('b').String()

	portNum := kingpin.Flag("port", "port to listen on").Default("8005").Short('p').Int()

	taskPortsHelp := "comma-separated list (length of `count`) of ports to use for spawned processes"
	taskPortsRaw := kingpin.Flag("innerPorts", taskPortsHelp).Short('P').String()

	taskAddress := kingpin.Flag("taskAddress", "address the spawned tasks are listening on").Default("localhost").Short('a').String()

	tsrHelp := fmt.Sprintf("logic to use for selecting which spawned process should receive the request: %d for sequential, %d for random (NOT IMPLEMENTED)", int(TSM_SEQUENTIAL), int(TSM_RANDOMIZED))
	taskSwitchRaw := kingpin.Flag("taskSwitchLogic", tsrHelp).Default(TSM_RANDOMIZED.String()).Short('s').Int()

	minWait := kingpin.Flag("minWait", "the shortest (in ms) a request should be delayed").Default("0").Short('w').Int()
	maxWait := kingpin.Flag("maxWait", "the longest (in ms) a request should be delayed").Default("2500").Short('W').Int()

	task := getTaskCmd(kingpin.Arg("command", "task to shuffle"))

	kingpin.Version("0.0")
	kingpin.Parse()

	config.http.address = *bindTo
	config.http.port = *portNum
	config.taskSwitch = TaskSwitchMethod(*taskSwitchRaw)

	var taskPortsSlice []string

	if len(*taskPortsRaw) > 0 {
		taskPortsSlice = strings.Split(*taskPortsRaw, ",")

		if len(taskPortsSlice) > *runCount {
			log.Fatal("More ports defined than task runners allowed by runCount")
			os.Exit(1)
		}
	} else {
		for i := 1; i <= *runCount; i++ {
			taskPortsSlice = append(taskPortsSlice, fmt.Sprintf("%d", i+config.http.port))
		}
	}

	for _, port := range taskPortsSlice {
		if len(port) >= 1 {
			newVal, err := strconv.Atoi(port)

			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}

			task := TaskRunner{0, HttpConfig{*taskAddress, newVal}, *task, time.Duration(*minWait) * time.Millisecond, time.Duration(*maxWait) * time.Millisecond}
			config.tasks = append(config.tasks, task)
		}
	}

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
