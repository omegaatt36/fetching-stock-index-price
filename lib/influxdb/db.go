package influxdb

import (
	"log"
)

type connPool struct {
	conns   chan *Client
	host    string
	db      string
	user    string
	pass    string
	maxConn int
}

var pool connPool

// Releaser is fund to release connection.
type Releaser func()

func (p *connPool) init() {
	p.conns = make(chan *Client, p.maxConn)
	for i := 0; i < p.maxConn; i++ {
		c, err := NewClient(p.host, p.db, p.user, p.pass)
		if err != nil {
			log.Panicln(err)
		}

		p.conns <- c
	}
}

// Get returns an influxdb client. We implemented connection pool to limit
// concurrent querys.
func Get() (*Client, Releaser) {
	c := <-pool.conns
	return c, func() {
		pool.conns <- c
	}
}
