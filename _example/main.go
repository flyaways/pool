package main

import (
	"fmt"
	"time"

	"github.com/flyaways/pool"
	"google.golang.org/grpc"
)

func main() { 
	//初始化创建连接池
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

	//异常处理
	if err != nil {
		log.Printf("%#v\n", err)
		return
	}

	//异常处理
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

	//todo 执行业务逻辑
	//conn=v.(*grpc.ClientConn)

	//放回连接池
	defer func() {
		if p.Put(v) != nil {
			log.Printf("%#v\n", err)
			return
		}
	}()

	//打印连接池大小
	log.Printf("len=%d\n", p.Len())
}