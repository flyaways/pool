package pool

import (
	"bufio"
	"encoding/gob"
	"errors"
	"net"
	"net/rpc"
	"time"
)

var (
	//ErrClosed 连接池已经关闭Error
	ErrClosed = errors.New("pool is closed")
)

//Pool 基本方法
type Pool interface {
	Get() (interface{}, error)

	Put(interface{}) error

	Close(interface{}) error

	Release()

	Len() int
}

func New(network, address string, initialCap, maxCap int, dialTimeout, idleTimeout time.Duration) (Pool, error) {
	factory := func() (interface{}, error) {
		conn, err := net.DialTimeout("tcp", address, dialTimeout)
		if err != nil {
			return nil, err
		}

		if network == "rpc" {
			encBuf := bufio.NewWriter(conn)
			c := rpc.NewClientWithCodec(&GobCoDec{
				Closer:  conn,
				Decoder: gob.NewDecoder(conn),
				Encoder: gob.NewEncoder(encBuf),
				EncBuf:  encBuf,
				Timeout: dialTimeout,
			})

			return c, err
		}
		return conn, err
	}

	close := func(v interface{}) error {
		if network == "rpc" {
			return v.(*rpc.Client).Close()
		}
		return v.(net.Conn).Close()
	}

	return NewChannelPool(&PoolConfig{
		InitialCap:  initialCap,
		MaxCap:      maxCap,
		Factory:     factory,
		Close:       close,
		IdleTimeout: idleTimeout,
	})
}
