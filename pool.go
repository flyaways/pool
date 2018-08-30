package pool

import (
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

var (

	//Address dial address
	address = "127.0.0.1:8080"

	//DialTimeout ...
	dialTimeout = time.Second * 5
)

//Pool 基本方法
type Pool interface {
	//Get 从pool中取一个连接
	Get() (interface{}, error)

	//Put 将连接放回pool中 例如：derfer p.Put(v)
	Put(interface{}) error

	//Close 关闭单条连接
	Close(interface{}) error

	//Release 释放连接池中所有链接
	Release()

	//Len 连接池中已有的连接数量
	Len() int
}

//Config  连接池相关配置
type Config struct {
	//远程地址
	Address string
	//连接池中拥有的最小连接数
	InitialCap int
	//连接池中拥有的最大的连接数
	MaxCap int
	//生成连接的方法
	Factory func() (interface{}, error)
	//关闭链接的方法
	Close func(interface{}) error
	//链接最大空闲时间，超过该事件则将失效
	IdleTimeout time.Duration
	//连接请求超时
	DialTimeout time.Duration
	//grpc DialOptions
	DialOptions []grpc.DialOption
}

//NewChannelPool 初始化链接
func NewChannelPool(poolConfig *Config) (Pool, error) {
	if poolConfig.InitialCap < 0 || poolConfig.MaxCap <= 0 || poolConfig.InitialCap > poolConfig.MaxCap {
		return nil, errors.New("invalid capacity settings")
	}

	address = poolConfig.Address
	dialTimeout = poolConfig.DialTimeout
	if poolConfig.DialOptions != nil && len(poolConfig.DialOptions) > 0 {
		dialOptions = poolConfig.DialOptions
	}

	c := &channelPool{
		conns:       make(chan *idleConn, poolConfig.MaxCap),
		factory:     poolConfig.Factory,
		close:       poolConfig.Close,
		idleTimeout: poolConfig.IdleTimeout,
	}

	for i := 0; i < poolConfig.InitialCap; i++ {
		conn, err := c.factory()
		if err != nil {
			c.Release()
			return nil, fmt.Errorf("factory is not able to fill the pool: %v", err)
		}
		c.conns <- &idleConn{conn: conn, t: time.Now()}
	}

	return c, nil
}
