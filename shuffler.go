package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Shuffler(_ http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	fmt.Printf("[child 1, 200ms] %v\n", p.ByName("path"))
}
