package client

import (
	"fmt"
	"slices"
	"sort"
	"strconv"
	"sync"
	"time"
)

type Runner struct {
	client         Client
	cursorClient   CursorClient
	generator      <-chan string
	expectedCalls  []string
	expectedResult []Response
}

func NewRunner0(collect bool) *Runner {
	expectedCalls := makeNExpectedCalls(100)

	r := &Runner{
		client:        NewSlowClient(100 * time.Millisecond),
		generator:     makeGeneratorForList(expectedCalls),
		expectedCalls: expectedCalls,
	}

	if collect {
		r.expectedResult = makeNExpectedResults(100)
	}

	return r
}

func NewRunner1() *Runner {
	expectedCalls := makeNExpectedCalls(100)

	return &Runner{
		client:         NewMaxConClient(10, 100*time.Millisecond),
		generator:      makeGeneratorForList(expectedCalls),
		expectedCalls:  expectedCalls,
		expectedResult: makeNExpectedResults(100),
	}
}

func NewRunner2() *Runner {
	expectedCalls := makeNExpectedCalls(100)

	return &Runner{
		client:         NewRateLimitedClient(),
		generator:      makeGeneratorForList(expectedCalls),
		expectedCalls:  expectedCalls,
		expectedResult: makeNExpectedResults(100),
	}

}
func NewRunner3() *Runner {
	expectedCalls := makeNExpectedCalls(100)

	return &Runner{
		client:        NewSlowClient(100 * time.Millisecond),
		generator:     makeGeneratorForList(expectedCalls),
		expectedCalls: expectedCalls,
	}
}

func NewRunner4() *Runner {
	expectedCalls := makeNExpectedCalls(100)
	return &Runner{
		client:         NewMaxConClient(10, time.Millisecond),
		cursorClient:   NewCursorClient(),
		expectedCalls:  expectedCalls,
		expectedResult: makeNExpectedResults(100),
	}
}

func NewRunner5() *Runner {
	expectedCalls := makeNExpectedCalls(100)

	generator := make(chan string)
	go func() {
		for i := 0; i < 2; i++ {
			for _, id := range expectedCalls {
				generator <- id
			}
		}
		close(generator)
	}()
	return &Runner{
		client:         NewSlowClient(100 * time.Millisecond),
		cursorClient:   NewCursorClient(),
		generator:      generator,
		expectedCalls:  expectedCalls,
		expectedResult: append(makeNExpectedResults(100), makeNExpectedResults(100)...),
	}
}

func (r *Runner) RunProcess(fn func(c Client, generator <-chan string) []Response, expectedDuration time.Duration) {
	c := &statsClient{
		client:  r.client,
		chCalls: make(chan string, 1000),
		wg:      sync.WaitGroup{},
	}

	startTime := time.Now()
	result := fn(c, r.generator)
	duration := time.Since(startTime)
	if duration > expectedDuration {
		fmt.Println("❌ Process took too long: ", duration, "expected <", expectedDuration)
	} else {
		fmt.Println("✔️ Process took: ", duration)
	}
	if r.expectedResult != nil {
		sort.Slice(result, func(i, j int) bool {
			return result[i].Body < result[j].Body
		})
		sort.Slice(r.expectedResult, func(i, j int) bool {
			return r.expectedResult[i].Body < r.expectedResult[j].Body
		})
		if !slices.Equal(result, r.expectedResult) {
			fmt.Println("❌️ Result is incorrect")
		}
	}

	c.assertCalls(r.expectedCalls)
}

func (r *Runner) RunProcessCursor(fn func(cursorClient CursorClient, client Client) []Response, expectedDuration time.Duration) {
	c := &statsClient{
		client:  r.client,
		chCalls: make(chan string, 1000),
		wg:      sync.WaitGroup{},
	}

	startTime := time.Now()
	result := fn(r.cursorClient, c)
	duration := time.Since(startTime)
	if duration > expectedDuration {
		fmt.Println("❌ Process took too long: ", duration, "expected <", expectedDuration)
	} else {
		fmt.Println("✔️ Process took: ", duration)
	}
	expectedCalls := make([]string, 0)
	for i := 0; i < 100; i++ {
		expectedCalls = append(expectedCalls, strconv.Itoa(i))
	}
	if r.expectedResult != nil {
		sort.Slice(result, func(i, j int) bool {
			return result[i].Body < result[j].Body
		})
		sort.Slice(r.expectedResult, func(i, j int) bool {
			return r.expectedResult[i].Body < r.expectedResult[j].Body
		})
		if !slices.Equal(result, r.expectedResult) {
			fmt.Println("❌️ Result is incorrect")
		}
	}

	c.assertCalls(expectedCalls)
}

type statsClient struct {
	client  Client
	chCalls chan string
	wg      sync.WaitGroup
}

func (c *statsClient) registerCall(id string) {
	if len(c.chCalls) == cap(c.chCalls) {
		panic("Client at capacity, too many calls")
	}
	c.chCalls <- id
	c.wg.Done()
}

func (c *statsClient) Call(id string) Response {
	c.wg.Add(1)
	defer c.registerCall(id)
	return c.client.Call(id)
}

func (c *statsClient) assertCalls(expectedCalls []string) {
	expected := make(map[string]bool)
	for _, id := range expectedCalls {
		expected[id] = false
	}
	//c.wg.Wait()
	close(c.chCalls)
	for id := range c.chCalls {
		called, ok := expected[id]
		if !ok {
			fmt.Println("❌ Unexpected call for id", id)
		}
		if called {
			fmt.Println("❌ Double calls for id", id)
		}
		expected[id] = true
	}

	missing := make([]string, 0)
	for id, called := range expected {
		if !called {
			missing = append(missing, id)
		}
	}
	if len(missing) > 0 {
		slices.Sort(missing)
		fmt.Printf("❌ Missing %d calls: %v\n", len(missing), missing)
	}
}

func makeNExpectedCalls(n int) []string {
	expectedCalls := make([]string, n)
	for i := 0; i < n; i++ {
		expectedCalls[i] = strconv.Itoa(i)
	}
	return expectedCalls
}

func makeNExpectedResults(n int) []Response {
	expectedResults := make([]Response, n)
	for i := 0; i < n; i++ {
		expectedResults[i] = Response{
			Body: fmt.Sprintf("Res%d", i),
		}
	}
	return expectedResults
}

func makeGeneratorForList(ids []string) <-chan string {
	generator := make(chan string)
	go func() {
		for _, id := range ids {
			generator <- id
		}
		close(generator)
	}()
	return generator
}
