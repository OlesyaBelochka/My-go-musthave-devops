package main

import (
	"fmt"
	"net/http"
)

func HandleMetrics(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Redacted())
}
