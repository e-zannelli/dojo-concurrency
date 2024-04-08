package main

import (
	"dojo-concurrency/client"
	"time"
)

func process(c client.Client, generator <-chan string) []client.Response {
	result := make([]client.Response, 0)
	cache := NewCache()

	for id := range generator {
		result = append(result, cache.Get(id, c.Call))
	}

	return result
}

func main() {
	client.NewRunner5().RunProcess(process, 110*time.Millisecond)
}
