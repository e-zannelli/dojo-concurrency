package main

import "dojo-concurrency/client"

func NewCache() Cache {
	return Cache{data: make(map[string]client.Response)}
}

type Cache struct {
	data map[string]client.Response
}

func (c *Cache) Get(key string, fetch func(string) client.Response) client.Response {
	if v, ok := c.data[key]; ok {
		return v
	}

	v := fetch(key)
	c.data[key] = v

	return v
}
