package main

import (
	"dojo-concurrency/client"
	"time"
)

// Client is limited to 1 new call per millisecond, if you call more the call will block longer
func process(c client.Client, generator <-chan string) []client.Response {
	result := make([]client.Response, 0)
	for id := range generator {
		result = append(result, c.Call(id))
	}

	return result
}

func main() {
	client.NewRunner2().RunProcess(process, 130*time.Millisecond)
}
