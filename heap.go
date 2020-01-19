package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"

	_ "net/http/pprof"
)

var store = [][]byte{}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		allocs := r.URL.Query().Get("alloc")
		alloc, err := strconv.Atoi(allocs)
		if err == nil && alloc > 0 {
			log.Printf("Allocating %d bytes", alloc)
			b := make([]byte, alloc)
			keep := r.URL.Query().Get("keep")
			if keep != "false" {
				store = append(store, b)
			}
		}
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		s := fmt.Sprintf("HeapAlloc:%v Sys:%v NextGC:%v NumGC:%v",
			m.HeapAlloc, m.Sys, m.NextGC, m.NumGC)
		log.Print(s)
		fmt.Fprintln(w, s)
	})

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
