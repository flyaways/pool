package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"net/rpc"
	"time"

	"github.com/flyaways/pool"
)

func newPool(network, address string, initialCap, maxCap int, dialTimeout, idleTimeout, rpcTimeout time.Duration) (pool.Pool, error) {
	factory := func() (interface{}, error) {
		conn, err := net.DialTimeout("tcp", address, dialTimeout)
		if err != nil {
			return nil, err
		}

		if network == "rpc" {
			encBuf := bufio.NewWriter(conn)
			c := rpc.NewClientWithCodec(&pool.GobCoDec{
				Closer:  conn,
				Decoder: gob.NewDecoder(conn),
				Encoder: gob.NewEncoder(encBuf),
				EncBuf:  encBuf,
				Timeout: rpcTimeout,
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

	return pool.NewChannelPool(&pool.PoolConfig{
		InitialCap:  initialCap,
		MaxCap:      maxCap,
		Factory:     factory,
		Close:       close,
		IdleTimeout: idleTimeout,
	})
}

func main() {
	p, err := newPool(
		"tcp",
		"127.0.0.1:8080",
		5,
		30,
		time.Second*5,
		time.Second*5,
		time.Second*5,
	)

	if err != nil {
		fmt.Printf("%#v\n", err)
		return
	}

	if p == nil {
		fmt.Printf("p= %#v\n", p)
		return
	}

	defer p.Release()

	v, err := p.Get()
	if err != nil {
		fmt.Printf("%#v\n", err)
		return
	}

	//do something
	//conn=v.(net.Conn)

	if p.Put(v) != nil {
		fmt.Printf("%#v\n", err)
		return
	}

	fmt.Printf("len=%d\n", p.Len())
}
