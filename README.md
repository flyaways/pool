# Pool

[![GoDoc](https://godoc.org/github.com/flyaways/pool?status.svg)](https://godoc.org/github.com/flyaways/pool)
[![Go Report Card](https://goreportcard.com/badge/github.com/flyaways/pool?style=flat-square)](https://goreportcard.com/report/github.com/flyaways/pool)
[![Build Status Travis](https://travis-ci.org/flyaways/pool.svg?branch=master)](https://travis-ci.org/flyaways/pool)
[![Build Status Semaphore](https://semaphoreci.com/api/v1/flyaways/pool/branches/master/shields_badge.svg)](https://semaphoreci.com/flyaways/pool)
[![Sourcegraph](https://sourcegraph.com/github.com/flyaways/pool/-/badge.svg)](https://sourcegraph.com/github.com/flyaways/pool?badge)
[![Open Source Helpers](https://www.codetriage.com/flyaways/pool/badges/users.svg)](https://www.codetriage.com/flyaways/pool)
[![LICENSE](https://img.shields.io/badge/licence-Apache%202.0-brightgreen.svg?style=flat-square)](https://github.com/flyaways/pool/blob/master/LICENSE)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/flyaways/pool.svg)
[![Release](https://img.shields.io/github/release/flyaways/pool.svg?style=flat-square)](https://github.com/flyaways/pool/releases)
[![LICENSE](https://img.shields.io/badge/license-Anti%20996-blue.svg)](https://github.com/996icu/996.ICU/blob/master/LICENSE)


Pool is Used to manage and reuse client connections to service cluster.

Pool provides several key features:

* **General Purpose** - Pool for GRPC,RPC,TCP.support RPC timeout.

* **Support Cluster** - Connet to Cluster.

* **Danamic Update** - Danamic update targets.

Pool runs on Linux, Mac OS X, and Windows.

**Note**: Random to pick a target to get one connection for loadbalance.

## Install
```sh
go get -u gopkg.in/flyaways/pool.v1
```
## Usage

```go
import "gopkg.in/flyaways/pool.v1"
```

## Example
```go
package main

import (
	"log"
	"time"

	"gopkg.in/flyaways/pool.v1"
	"google.golang.org/grpc"
)

func main() {
	options := &pool.Options{
		InitTargets:  []string{"127.0.0.1:8080"},
		InitCap:      5,
		MaxCap:       30,
		DialTimeout:  time.Second * 5,
		IdleTimeout:  time.Second * 60,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	}

	
	p, err := pool.NewGRPCPool(options, grpc.WithInsecure())//for grpc
	//p, err := pool.NewRPCPool(options) 			//for rpc
	//p, err := pool.NewTCPPool(options)			//for tcp

	if err != nil {
		log.Printf("%#v\n", err)
		return
	}

	if p == nil {
		log.Printf("p= %#v\n", p)
		return
	}

	defer p.Close()

	//todo
	//danamic update targets
	//options.Input()<-&[]string{}

	conn, err := p.Get()
	if err != nil {
		log.Printf("%#v\n", err)
		return
	}

	defer p.Put(conn)

	//todo
	//conn.DoSomething()

	log.Printf("len=%d\n", p.IdleCount())
}

```

## Reference
 * [https://github.com/fatih/pool](https://github.com/fatih/pool)
 * [https://github.com/silenceper/pool]( https://github.com/silenceper/pool)
 * [https://github.com/daizuozhuo/rpc-example]( https://github.com/daizuozhuo/rpc-example)
 
## Contribution Welcomed !

Contributors

* [flyaways](https://github.com/flyaways) 

## License
* The MIT License (MIT) - see LICENSE for more details

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bhttps%3A%2F%2Fgithub.com%2Fflyaways%2Fpool.svg?type=large)](https://app.fossa.io/projects/git%2Bhttps%3A%2F%2Fgithub.com%2Fflyaways%2Fpool?ref=badge_large)
