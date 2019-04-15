package pool

import (
	"net"
	"sync"
	"time"
)

//TCPPool pool info
type TCPPool struct {
	Mu          sync.Mutex
	IdleTimeout time.Duration
	conns       chan *tcpIdleConn
	factory     func() (net.Conn, error)
	close       func(net.Conn) error
}

type tcpIdleConn struct {
	conn net.Conn
	t    time.Time
}

//Get get from pool
func (c *TCPPool) Get() (net.Conn, error) {
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
func (c *TCPPool) Put(conn net.Conn) error {
	if conn == nil {
		return errRejected
	}

	c.Mu.Lock()
	defer c.Mu.Unlock()

	if c.conns == nil {
		return c.close(conn)
	}

	select {
	case c.conns <- &tcpIdleConn{conn: conn, t: time.Now()}:
		return nil
	default:
		//连接池已满，直接关闭该链接
		return c.close(conn)
	}
}

//Close close all connection
func (c *TCPPool) Close() {
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
func (c *TCPPool) IdleCount() int {
	c.Mu.Lock()
	conns := c.conns
	c.Mu.Unlock()
	return len(conns)
}

//NewTCPPool init tcp pool
func NewTCPPool(o *Options) (*TCPPool, error) {
	if err := o.validate(); err != nil {
		return nil, err
	}

	//init pool
	pool := &TCPPool{
		conns: make(chan *tcpIdleConn, o.MaxCap),
		factory: func() (net.Conn, error) {
			target := o.nextTarget()
			if target == "" {
				return nil, errTargets
			}

			return net.DialTimeout("tcp", target, o.DialTimeout)
		},
		close:       func(v net.Conn) error { return v.Close() },
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
		pool.conns <- &tcpIdleConn{conn: conn, t: time.Now()}
	}

	return pool, nil
}
