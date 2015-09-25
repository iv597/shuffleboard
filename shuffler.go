package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"math"
	"math/rand"
	"net/http"
	"time"
)

func nextRunner() int {
	if config.taskSwitch == TSM_RANDOMIZED {
		return round(rand.Float64() * float64(len(config.tasks)))
	}

	return 0 // for now, not implemented TODO
}

func Shuffler(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	runner := nextRunner()

	minDelay := int(config.tasks[runner].minWait)
	maxDelay := int(config.tasks[runner].maxWait)

	delay := time.Duration(rand.Intn(maxDelay-minDelay) + minDelay)

	fmt.Printf("[child %d, %dms] %v\n", runner, (delay / time.Millisecond), p.ByName("path"))
	time.Sleep(delay)
	fmt.Fprintf(w, "%v\n", config)
}
