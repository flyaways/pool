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

	p, err := pool.NewGRPCPool(options, grpc.WithInsecure())

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
