package main

import (
	"dojo-concurrency/client"
	"time"
)

func getPages(c client.CursorClient) []client.Cursor {
	res := make([]client.Cursor, 0)
	pos := ""
	for {
		resp := c.Call(pos)
		res = append(res, resp)
		next, err := resp.Next()
		if err != nil {
			break
		}
		pos = next
	}

	return res
}

func process(cCursor client.CursorClient, c client.Client) []client.Response {
	pages := getPages(cCursor)

	result := make([]client.Response, 0)
	for _, p := range pages {
		for _, id := range p.Content {
			result = append(result, c.Call(id))

		}
	}

	return result
}

func main() {
	client.NewRunner4().RunProcessCursor(process, 130*time.Millisecond)
}
