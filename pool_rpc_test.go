package pool

import (
	"log"
	"net/rpc"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestNewRPCPool(t *testing.T) {
	type args struct {
		o *Options
	}
	tests := []struct {
		name    string
		args    args
		want    *RPCPool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRPCPool(tt.args.o)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRPCPool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRPCPool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRPCPool_Get(t *testing.T) {
	type fields struct {
		Mu          sync.Mutex
		IdleTimeout time.Duration
		conns       chan *rpcIdleConn
		factory     func() (*rpc.Client, error)
		close       func(*rpc.Client) error
	}
	tests := []struct {
		name    string
		fields  fields
		want    *rpc.Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RPCPool{
				Mu:          tt.fields.Mu,
				IdleTimeout: tt.fields.IdleTimeout,
				conns:       tt.fields.conns,
				factory:     tt.fields.factory,
				close:       tt.fields.close,
			}
			got, err := c.Get()
			if (err != nil) != tt.wantErr {
				t.Errorf("RPCPool.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RPCPool.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleRPCPool() {
	options := &Options{
		InitTargets:  []string{"127.0.0.1:8080"},
		InitCap:      5,
		MaxCap:       30,
		DialTimeout:  time.Second * 5,
		IdleTimeout:  time.Second * 60,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	}

	p, err := NewRPCPool(options)

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
