package client

import (
	"fmt"
	"golang.org/x/time/rate"
	"math"
	"math/rand"
	"slices"
	"strconv"
	"sync"
	"time"
)

type Client interface {
	Call(id string) Response
}

type CursorClient interface {
	Call(pos string) Cursor
}

type Cursor struct {
	Content []string
	next    string
}

func (c *Cursor) Next() (string, error) {
	if c.next == "" {
		return "", fmt.Errorf("no next page")
	}

	return c.next, nil
}

type Response struct {
	Body string
}

type SimpleClient struct{}

func (c SimpleClient) Call(id string) Response {
	return Response{
		Body: fmt.Sprintf("Res%s", id),
	}
}

type SlowClient struct {
	c     Client
	delay time.Duration
}

func NewSlowClient(delay time.Duration) *SlowClient {
	return &SlowClient{
		c:     SimpleClient{},
		delay: delay,
	}
}

func (c *SlowClient) Call(id string) Response {
	jitter := int(math.Round(float64(c.delay) * 0.1))
	delay := time.Duration(int(c.delay) - jitter/2 + rand.Intn(jitter))

	time.Sleep(delay)

	return c.c.Call(id)
}

type MyCursorClient struct {
	cursorKeys []string
}

func NewCursorClient() *MyCursorClient {
	return &MyCursorClient{
		cursorKeys: []string{"", "second", "foo", "bar", "last"},
	}
}

func (c *MyCursorClient) Call(pos string) Cursor {
	time.Sleep(time.Duration(8+rand.Intn(5)) * time.Millisecond)
	next := ""
	index := slices.Index(c.cursorKeys, pos)
	if index < 0 {
		panic(fmt.Sprintf("Unknown position %s", pos))
	}
	if len(c.cursorKeys) > index+1 {
		next = c.cursorKeys[index+1]
	}
	content := make([]string, 0, 20)
	startIndex := 0
	switch pos {
	case "":
		startIndex = 0
	case "second":
		startIndex = 20
	case "foo":
		startIndex = 40
	case "bar":
		startIndex = 60
	case "last":
		startIndex = 80
	}
	for i := startIndex; i < startIndex+20; i++ {
		content = append(content, strconv.Itoa(i))
	}

	return Cursor{
		Content: content,
		next:    next,
	}
}

type MaxConClient struct {
	c       Client
	delay   time.Duration
	maxCon  int
	currCon int
	mu      sync.Mutex
}

func NewMaxConClient(maxCon int, delay time.Duration) *MaxConClient {
	return &MaxConClient{
		c:      NewSlowClient(10 * time.Millisecond),
		delay:  delay,
		maxCon: maxCon,
	}
}

func (c *MaxConClient) take() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.currCon == c.maxCon {
		return fmt.Errorf("max connections reached")
	}
	c.currCon++
	return nil
}

func (c *MaxConClient) release() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.currCon--
}

func (c *MaxConClient) Call(id string) Response {
	for {
		err := c.take()
		if err == nil {
			break
		}
		time.Sleep(c.delay)
	}
	defer c.release()

	return c.c.Call(id)
}

type RateLimitedClient struct {
	c         Client
	connCount int
	limiter   rate.Limiter
}

func NewRateLimitedClient() *RateLimitedClient {
	return &RateLimitedClient{
		c:       NewSlowClient(10 * time.Millisecond),
		limiter: *rate.NewLimiter(rate.Every(990*time.Microsecond), 2),
	}
}

func (c *RateLimitedClient) Call(id string) Response {
	for {
		if c.limiter.Allow() {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	return c.c.Call(id)
}
