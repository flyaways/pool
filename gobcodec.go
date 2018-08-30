package pool

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"net/rpc"
	"time"
)

//GobCoDec ...
type GobCoDec struct {
	Timeout time.Duration
	Closer  io.ReadWriteCloser
	Decoder *gob.Decoder
	Encoder *gob.Encoder
	EncBuf  *bufio.Writer
}

//WriteRequest ...
func (c *GobCoDec) WriteRequest(r *rpc.Request, body interface{}) (err error) {
	if err = c.timeoutCoder(r, "write request"); err != nil {
		return
	}

	if err = c.timeoutCoder(body, "write request body"); err != nil {
		return
	}

	return c.EncBuf.Flush()
}

//ReadResponseHeader ...
func (c *GobCoDec) ReadResponseHeader(r *rpc.Response) error {
	return c.Decoder.Decode(r)
}

//ReadResponseBody ...
func (c *GobCoDec) ReadResponseBody(body interface{}) error {
	return c.Decoder.Decode(body)
}

//Close ...
func (c *GobCoDec) Close() error {
	return c.Closer.Close()
}

func (c *GobCoDec) timeoutCoder(e interface{}, msg string) error {
	if c.Timeout < 0 {
		c.Timeout = time.Second * 5
	}

	echan := make(chan error, 1)
	go func() { echan <- c.Encoder.Encode(e) }()

	select {
	case e := <-echan:
		return e
	case <-time.After(c.Timeout):
		return fmt.Errorf("Timeout %s", msg)
	}
}
