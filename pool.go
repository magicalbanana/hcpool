package hcpool

import (
	"net/http"
	"time"
)

// Pool ...
type Pool struct {
	clients chan *http.Client
}

// Options ...
type Options struct {
	DisableKeepAlives   bool
	MaxIdleConnsPerHost int
	RoundTripTimeout    time.Duration
}

// NewPool ...
func NewPool(size int, opts Options) *Pool {
	pool := new(Pool)
	pool.clients = make(chan *http.Client, size)

	tr := &TimeoutTransport{
		Transport: http.Transport{
			DisableKeepAlives:   opts.DisableKeepAlives,
			MaxIdleConnsPerHost: opts.MaxIdleConnsPerHost,
		},
		RoundTripTimeout: opts.RoundTripTimeout,
	}

	for i := 0; i < size; i++ {
		c := &http.Client{Transport: tr}
		pool.clients <- c
	}

	return pool
}

// getClient ...
func (p *Pool) getClient() *http.Client {
	return <-p.clients
}

// returnClient ...
func (p *Pool) returnClient(client *http.Client) {
	p.clients <- client
}

// Do ...
func (p *Pool) Do(req *http.Request) (*http.Response, error) {
	client := p.getClient()
	defer p.returnClient(client)

	resp, doErr := client.Do(req)
	if doErr != nil {
		return nil, doErr
	}

	return resp, nil
}

// Close all tehe channels
func (p *Pool) Close() {
	close(p.clients)
}
