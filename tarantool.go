// Package tarantool provides a client for tarantool 1.5.
package tarantool15

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/lomik/go-tnt"
)

var ErrEmptyTuple = errors.New("empty tuple")

type ExecutionError struct {
	Space interface{}
	Index int
	Err   error
}

func (e ExecutionError) Error() string {
	return fmt.Sprintf("could not execute tarantool request, space %v, index %d: %v", e.Space, e.Index, e.Err)
}

type Config struct {
	MaxReconnects     int
	ReconnectInterval time.Duration
	ConnectTimeout    time.Duration
	QueryTimeout      time.Duration
	// DebugfFunc could be set to debug connection creation.
	DebugfFunc func(string, ...interface{})
}

func Connect(addr string, conf Config) (*Conn, error) {
	db := &Conn{
		addr:           addr,
		maxReconnects:  conf.MaxReconnects,
		reconnInterval: conf.ReconnectInterval,
		connTimeout:    conf.ConnectTimeout,
		queryTimeout:   conf.QueryTimeout,
		logf:           conf.DebugfFunc,
	}
	err := db.connectTnt(false)
	return db, err
}

type Conn struct {
	addr           string
	maxReconnects  int
	reconnInterval time.Duration
	connTimeout    time.Duration
	queryTimeout   time.Duration

	logf func(string, ...interface{})
	mu   sync.RWMutex
	conn *tnt.Connection
}

func (c *Conn) connectTnt(reconnect bool) error {
	opts := &tnt.Options{
		ConnectTimeout: c.connTimeout,
		QueryTimeout:   c.queryTimeout,
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	var attempts int
	for c.conn == nil || c.conn.IsClosed() {
		tntConn, err := tnt.Connect(c.addr, opts)
		if err == nil || !reconnect {
			c.conn = tntConn
			return err
		}

		if c.logf != nil {
			c.logf("connection to %s (%d of %d) failed: %v", c.addr, attempts, c.maxReconnects, err)
		}

		if c.maxReconnects > 0 && attempts > c.maxReconnects {
			return err
		}

		attempts++

		c.mu.Unlock()
		time.Sleep(c.reconnInterval)
		c.mu.Lock()
	}

	return nil
}

func (c *Conn) Close() error {
	c.conn.Close()
	return nil
}

func (c *Conn) Execute(q Query) (result []tnt.Tuple, err error) {
	var closed bool
	c.mu.RLock()
	closed = c.conn.IsClosed()
	c.mu.RUnlock()

	if closed {
		err := c.connectTnt(true)
		if err != nil {
			return nil, err
		}
	}

	return c.conn.Execute(q)
}
