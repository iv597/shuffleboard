package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

func Shuffler(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	fmt.Printf("[child 1, 200ms] %v\n", p.ByName("path"))
	time.Sleep(200 * time.Millisecond)
	fmt.Fprintf(w, "%v\n", config)
}
