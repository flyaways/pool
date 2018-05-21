# pool

Golang tcp,rpc pool,support rpc timeout

## usage
```go
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
		"tcp",//rpc
		"127.0.0.1:8080",//address
		5,//initialCap
		30,//maxCap
		time.Second*5,//dialTimeout
		time.Second*5,//idleTimeout
		time.Second*5,//rpcTimeout
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

	c, err := p.Get()
	if err != nil {
		fmt.Printf("%#v\n", err)
		return
	}

	//conn=c.(net.Conn)

	if p.Put(c) != nil {
		fmt.Printf("%#v\n", err)
		return
	}

	fmt.Printf("len=%d\n", p.Len())
}

```

#### reference
[https://github.com/fatih/pool](https://github.com/fatih/pool)
[https://github.com/silenceper/pool]( https://github.com/silenceper/pool)
[https://github.com/daizuozhuo/rpc-example]( https://github.com/daizuozhuo/rpc-example)

## License

The MIT License (MIT) - see LICENSE for more details
