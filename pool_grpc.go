package pool

import (
	"context"
	"sync"
	"time"

	"google.golang.org/grpc"
)

//GRPCPool pool info
type GRPCPool struct {
	Mu          sync.Mutex
	IdleTimeout time.Duration
	conns       chan *grpcIdleConn
	factory     func() (*grpc.ClientConn, error)
	close       func(*grpc.ClientConn) error
}

type grpcIdleConn struct {
	conn *grpc.ClientConn
	t    time.Time
}

//Get get from pool
func (c *GRPCPool) Get() (*grpc.ClientConn, error) {
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
func (c *GRPCPool) Put(conn *grpc.ClientConn) error {
	if conn == nil {
		return errRejected
	}

	c.Mu.Lock()
	defer c.Mu.Unlock()

	if c.conns == nil {
		return c.close(conn)
	}

	select {
	case c.conns <- &grpcIdleConn{conn: conn, t: time.Now()}:
		return nil
	default:
		//连接池已满，直接关闭该链接
		return c.close(conn)
	}
}

//Close close pool
func (c *GRPCPool) Close() {
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
func (c *GRPCPool) IdleCount() int {
	c.Mu.Lock()
	conns := c.conns
	c.Mu.Unlock()
	return len(conns)
}

//NewGRPCPool init grpc pool
func NewGRPCPool(o *Options, dialOptions ...grpc.DialOption) (*GRPCPool, error) {
	if err := o.validate(); err != nil {
		return nil, err
	}

	//init pool
	pool := &GRPCPool{
		conns: make(chan *grpcIdleConn, o.MaxCap),
		factory: func() (*grpc.ClientConn, error) {
			target := o.nextTarget()
			if target == "" {
				return nil, errTargets
			}

			ctx, cancel := context.WithTimeout(context.Background(), o.DialTimeout)
			defer cancel()

			return grpc.DialContext(ctx, target, dialOptions...)
		},
		close:       func(v *grpc.ClientConn) error { return v.Close() },
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
		pool.conns <- &grpcIdleConn{conn: conn, t: time.Now()}
	}

	return pool, nil
}
