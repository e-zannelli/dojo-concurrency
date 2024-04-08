package main

import (
	"dojo-concurrency/client"
	"time"
)

func process(c client.Client, generator <-chan string) []client.Response {
	for id := range generator {
		go c.Call(id)
	}

	return make([]client.Response, 0)
}

func main() {
	client.NewRunner0(false).RunProcess(process, 115*time.Millisecond)
}
