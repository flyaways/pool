# pool
[![Go Report Card](https://goreportcard.com/badge/github.com/flyaways/pool?style=flat-square)](https://goreportcard.com/report/github.com/flyaways/pool)
[![Build Status Travis](https://travis-ci.org/flyaways/pool.svg?branch=master)](https://travis-ci.org/flyaways/pool)
[![Build Status Semaphore](https://semaphoreci.com/api/v1/flyaways/pool/branches/master/shields_badge.svg)](https://semaphoreci.com/flyaways/pool)
[![LICENSE](https://img.shields.io/badge/licence-Apache%202.0-brightgreen.svg?style=flat-square)](https://github.com/flyaways/pool/blob/master/LICENSE)

pool is Used to manage and reuse connections.
thread safe connection pool for tcp,rpc,grpc. 
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
	"google.golang.org/grpc"
)

func main() { 
	//创建连接池
	p, err := pool.NewChannelPool(&PoolConfig{
		InitialCap:  5,
		MaxCap:      30,
		Factory:     pool.FactoryGRPC,
		Close:       pool.CloseGRPC,
		IdleTimeout: time.Second*5,
		DialTimeout:time.Second*5,
		 DialOptions:[]grpc.DialOption{
         	grpc.WithInsecure()
		 },
	})

	if err != nil {
		log.Printf("%#v\n", err)
		return
	}

	if p == nil {
		log.Printf("p= %#v\n", p)
		return
	}

	//释放连接池
	defer p.Release()

	//获取一个连接
	v, err := p.Get()
	if err != nil {
		log.Printf("%#v\n", err)
		return
	}

	//todo
	//conn=v.(*grpc.ClientConn)

	//用完放回连接
	if p.Put(v) != nil {
		log.Printf("%#v\n", err)
		return
	}

	//打印连接池大小
	log.Printf("len=%d\n", p.Len())
}
```

## Reference
 * [https://github.com/fatih/pool](https://github.com/fatih/pool)
 * [https://github.com/silenceper/pool]( https://github.com/silenceper/pool)
 * [https://github.com/daizuozhuo/rpc-example]( https://github.com/daizuozhuo/rpc-example)

## License
* The MIT License (MIT) - see LICENSE for more details

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bhttps%3A%2F%2Fgithub.com%2Fflyaways%2Fpool.svg?type=large)](https://app.fossa.io/projects/git%2Bhttps%3A%2F%2Fgithub.com%2Fflyaways%2Fpool?ref=badge_large)
