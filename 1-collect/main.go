package main

import (
	"dojo-concurrency/client"
	"time"
)

func process(c client.Client, generator <-chan string) []client.Response {
	result := make([]client.Response, 0)
	for id := range generator {
		result = append(result, c.Call(id))
	}

	return result
}

func main() {
	client.NewRunner0(true).RunProcess(process, 115*time.Millisecond)
}
