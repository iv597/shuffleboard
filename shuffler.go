package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func nextRunner() int {
	if len(config.tasks) == 1 {
		return 0
	}

	if config.taskSwitch == TSM_RANDOMIZED {
		// the logic here is [0, n) so this is safe
		return rand.Intn(len(config.tasks))
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

	resp, err := http.DefaultTransport.RoundTrip(req)

	if err != nil {
		if resp.StatusCode != 302 {
			panic(fmt.Sprintf("Error making request: %v", err))
		}
	}

	for key, val := range resp.Header {
		w.Header().Set(key, val[0])
	}

	w.WriteHeader(resp.StatusCode)

	if resp.Body != nil {
		go func() {
			io.Copy(w, resp.Body)
			resp.Body.Close()
		}()
	}

	log.Printf("[child %d, delayed %dms, upstream %dms] %v\n", runner, (delay / time.Millisecond), time.Since(reqSentTime)/time.Millisecond, p.ByName("path"))
}
