package main

import (
	"http"
)

// Do sends an HTTP request with additional logging, instrumentation,
// authorization, fault tolerance and load balancing.
func (c *Client) Do(r *http.Request) (resp *http.Response, err error) {
	r.Header.Add("Authorization", c.Token)
	for i := 0; i <= c.Tolerance; i++ {
		r.URL.Host = c.Backends[atomic.AddUint64(&c.robin, 1)%uint64(len(c.Backends))]
		log.Print("%s: %s %s", r.UserAgent, r.Method, r.URL)
		start := time.Now()
		resp, err = http.DefaultClient.Do(r)
		c.latency.Observe(time.Since(start))
		c.requests.Add(1)
		if err != nil {
			time.Sleep(time.Duration(i) * c.Backoff)
			continue
		}
		break
	}
	return res, err
}
