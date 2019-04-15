package pool

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"net/rpc"
	"sync"
	"time"
)

//RPCPool pool info
type RPCPool struct {
	Mu          sync.Mutex
	IdleTimeout time.Duration
	conns       chan *rpcIdleConn
	factory     func() (*rpc.Client, error)
	close       func(*rpc.Client) error
}

type rpcIdleConn struct {
	conn *rpc.Client
	t    time.Time
}

//Get get from pool
func (c *RPCPool) Get() (*rpc.Client, error) {
	c.Mu.Lock()
	conns := c.conns
	c.Mu.Unlock()

	if conns == nil {
		return nil, errClosed
	}
	for {
		select {
		case wrapConn := <-conns:
			if wrapConn == nil {
				return nil, errClosed
			}
			//判断是否超时，超时则丢弃
			if timeout := c.IdleTimeout; timeout > 0 {
				if wrapConn.t.Add(timeout).Before(time.Now()) {
					//丢弃并关闭该链接
					c.close(wrapConn.conn)
					continue
				}
			}
			return wrapConn.conn, nil
		default:
			conn, err := c.factory()
			if err != nil {
				return nil, err
			}

			return conn, nil
		}
	}
}

//Put put back to pool
func (c *RPCPool) Put(conn *rpc.Client) error {
	if conn == nil {
		return errRejected
	}

	c.Mu.Lock()
	defer c.Mu.Unlock()

	if c.conns == nil {
		return c.close(conn)
	}

	select {
	case c.conns <- &rpcIdleConn{conn: conn, t: time.Now()}:
		return nil
	default:
		//连接池已满，直接关闭该链接
		return c.close(conn)
	}
}

//Close close all connection
func (c *RPCPool) Close() {
	c.Mu.Lock()
	conns := c.conns
	c.conns = nil
	c.factory = nil
	closeFun := c.close
	c.close = nil
	c.Mu.Unlock()

	if conns == nil {
		return
	}

	close(conns)
	for wrapConn := range conns {
		closeFun(wrapConn.conn)
	}
}

//IdleCount idle connection count
func (c *RPCPool) IdleCount() int {
	c.Mu.Lock()
	conns := c.conns
	c.Mu.Unlock()
	return len(conns)
}

//Codec ...
type Codec struct {
	Timeout time.Duration
	Closer  io.ReadWriteCloser
	Decoder *gob.Decoder
	Encoder *gob.Encoder
	EncBuf  *bufio.Writer
}

//WriteRequest ...
func (c *Codec) WriteRequest(r *rpc.Request, body interface{}) (err error) {
	if err = c.timeoutCoder(r, "write request"); err != nil {
		return
	}

	if err = c.timeoutCoder(body, "write request body"); err != nil {
		return
	}

	return c.EncBuf.Flush()
}

//ReadResponseHeader ...
func (c *Codec) ReadResponseHeader(r *rpc.Response) error {
	return c.Decoder.Decode(r)
}

//ReadResponseBody ...
func (c *Codec) ReadResponseBody(body interface{}) error {
	return c.Decoder.Decode(body)
}

//Close ...
func (c *Codec) Close() error {
	return c.Closer.Close()
}

func (c *Codec) timeoutCoder(e interface{}, msg string) error {
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

//NewRPCPool init rpc pool
func NewRPCPool(o *Options) (*RPCPool, error) {
	if err := o.validate(); err != nil {
		return nil, err
	}

	//init pool
	pool := &RPCPool{
		conns: make(chan *rpcIdleConn, o.MaxCap),
		factory: func() (*rpc.Client, error) {
			target := o.nextTarget()
			if target == "" {
				return nil, errTargets
			}

			conn, err := net.DialTimeout("tcp", target, o.DialTimeout)
			if err != nil {
				return nil, err
			}

			encBuf := bufio.NewWriter(conn)
			p := rpc.NewClientWithCodec(&Codec{
				Closer:  conn,
				Decoder: gob.NewDecoder(conn),
				Encoder: gob.NewEncoder(encBuf),
				EncBuf:  encBuf,
				Timeout: o.WriteTimeout,
			})

			return p, err
		},
		close:       func(v *rpc.Client) error { return v.Close() },
		IdleTimeout: o.IdleTimeout,
	}

	//danamic update targets
	o.update()

	//init make conns
	for i := 0; i < o.InitCap; i++ {
		conn, err := pool.factory()
		if err != nil {
			pool.Close()
			return nil, err
		}
		pool.conns <- &rpcIdleConn{conn: conn, t: time.Now()}
	}

	return pool, nil
}
