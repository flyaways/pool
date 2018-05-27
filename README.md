# pool

pool is Used to manage and reuse connections.
thread safe connection pool for tcp,rpc. 
Support rpc timeout.

## Install
```sh
go get -u github.com/flyaways/pool
```

## Example
```go
package main

import (
	"fmt"
	"time"

	"github.com/flyaways/pool"
)

func main() {
	p, err := pool.New(
		"tcp",
		"127.0.0.1:8080",
		5,
		30,
		time.Second*5, //dialTimeout
		time.Second*5, //idleTimeout
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

	//todo
	//conn=v.(net.Conn)

	if p.Put(v) != nil {
		fmt.Printf("%#v\n", err)
		return
	}

	fmt.Printf("len=%d\n", p.Len())
}

```

## Reference
 * [https://github.com/fatih/pool](https://github.com/fatih/pool)
 * [https://github.com/silenceper/pool]( https://github.com/silenceper/pool)
 * [https://github.com/daizuozhuo/rpc-example]( https://github.com/daizuozhuo/rpc-example)

## License
* The MIT License (MIT) - see LICENSE for more details
