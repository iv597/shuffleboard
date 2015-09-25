package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"math/rand"
	"net/http"
	"time"
)

func nextRunner() int {
	if len(config.tasks) == 1 {
		return 0
	}

	if config.taskSwitch == TSM_RANDOMIZED {
		return round(rand.Float64() * float64(len(config.tasks)))
	}

	return 0 // for now, not implemented TODO
}

func Shuffler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rand.Seed(time.Now().UnixNano())

	runner := nextRunner()

	minDelay := int(config.tasks[runner].minWait)
	maxDelay := int(config.tasks[runner].maxWait)

	delay := time.Duration(rand.Intn(maxDelay-minDelay) + minDelay)

	upUrl := fmt.Sprintf("http://%s:%d%s", config.tasks[runner].http.address, config.tasks[runner].http.port, p.ByName("path"))

	time.Sleep(delay)

	req, err := http.NewRequest(r.Method, upUrl, r.Body)

	if err != nil {
		panic("Could not generate the upstream request")
	}

	req.ContentLength = r.ContentLength
	req.Header = r.Header

	for _, cookie := range r.Cookies() {
		req.AddCookie(cookie)
	}

	reqSentTime := time.Now()

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		panic(fmt.Sprintf("Error making request: %v", err))
	}

	defer resp.Body.Close()

	io.Copy(w, resp.Body)

	fmt.Printf("[child %d, delayed %dms, upstream %dms] %v\n", runner, (delay / time.Millisecond), time.Since(reqSentTime)/time.Millisecond, p.ByName("path"))
}
