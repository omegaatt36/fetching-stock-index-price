package influxdb

import (
	"log"
	"sync"
	"time"

	// no idea why.
	_ "github.com/influxdata/influxdb1-client"
	inf "github.com/influxdata/influxdb1-client/v2"
	"github.com/pkg/errors"
)

// Client wraps influxdb client.
type Client struct {
	sync.Mutex

	inf.Client
	config inf.BatchPointsConfig
	buffer inf.BatchPoints
	dbName string
}

// NewClient creates client.
func NewClient(addr, db, user, pass string) (*Client, error) {
	c, err := inf.NewHTTPClient(inf.HTTPConfig{
		Addr:     addr,
		Username: user,
		Password: pass,
	})
	if err != nil {
		return nil, err
	}

	_, _, err = c.Ping(5 * time.Second)
	if err != nil {
		log.Fatal(err)
	}

	return &Client{
		Client: c,
		config: inf.BatchPointsConfig{
			Database: db,
		},
		dbName: db,
	}, nil
}

// AsyncWrite adds data to queue.
func (c *Client) AsyncWrite(measurement string, tags map[string]string, fields map[string]interface{}, timestamp time.Time) {
	c.Lock()
	defer c.Unlock()

	if c.buffer == nil {
		c.buffer, _ = inf.NewBatchPoints(c.config)
	}
	point, err := inf.NewPoint(measurement, tags, fields, timestamp)
	if err != nil {
		log.Panicln(err)
	}
	c.buffer.AddPoint(point)
}

// Flush sends data.
func (c *Client) Flush() error {
	c.Lock()
	defer c.Unlock()
	if c.buffer == nil {
		return nil
	}

	err := c.Write(c.buffer)
	c.buffer = nil
	return err
}

func checkResponse(response *inf.Response) error {
	if err := response.Error(); err != nil {
		return err
	}

	if len(response.Results) == 0 {
		return errors.New("empty results")
	}

	if len(response.Results[0].Series) == 0 {
		return errors.New("empty series")
	}

	return nil
}
