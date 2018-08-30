package pool

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

func Test_channelPool_Get(t *testing.T) {
	type fields struct {
		mu          sync.Mutex
		conns       chan *idleConn
		factory     func() (interface{}, error)
		close       func(interface{}) error
		idleTimeout time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &channelPool{
				mu:          tt.fields.mu,
				conns:       tt.fields.conns,
				factory:     tt.fields.factory,
				close:       tt.fields.close,
				idleTimeout: tt.fields.idleTimeout,
			}
			got, err := c.Get()
			if (err != nil) != tt.wantErr {
				t.Errorf("channelPool.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("channelPool.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_channelPool_Put(t *testing.T) {
	type fields struct {
		mu          sync.Mutex
		conns       chan *idleConn
		factory     func() (interface{}, error)
		close       func(interface{}) error
		idleTimeout time.Duration
	}
	type args struct {
		conn interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &channelPool{
				mu:          tt.fields.mu,
				conns:       tt.fields.conns,
				factory:     tt.fields.factory,
				close:       tt.fields.close,
				idleTimeout: tt.fields.idleTimeout,
			}
			if err := c.Put(tt.args.conn); (err != nil) != tt.wantErr {
				t.Errorf("channelPool.Put() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_channelPool_Close(t *testing.T) {
	type fields struct {
		mu          sync.Mutex
		conns       chan *idleConn
		factory     func() (interface{}, error)
		close       func(interface{}) error
		idleTimeout time.Duration
	}
	type args struct {
		conn interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &channelPool{
				mu:          tt.fields.mu,
				conns:       tt.fields.conns,
				factory:     tt.fields.factory,
				close:       tt.fields.close,
				idleTimeout: tt.fields.idleTimeout,
			}
			if err := c.Close(tt.args.conn); (err != nil) != tt.wantErr {
				t.Errorf("channelPool.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_channelPool_Release(t *testing.T) {
	type fields struct {
		mu          sync.Mutex
		conns       chan *idleConn
		factory     func() (interface{}, error)
		close       func(interface{}) error
		idleTimeout time.Duration
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &channelPool{
				mu:          tt.fields.mu,
				conns:       tt.fields.conns,
				factory:     tt.fields.factory,
				close:       tt.fields.close,
				idleTimeout: tt.fields.idleTimeout,
			}
			c.Release()
		})
	}
}

func Test_channelPool_Len(t *testing.T) {
	type fields struct {
		mu          sync.Mutex
		conns       chan *idleConn
		factory     func() (interface{}, error)
		close       func(interface{}) error
		idleTimeout time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &channelPool{
				mu:          tt.fields.mu,
				conns:       tt.fields.conns,
				factory:     tt.fields.factory,
				close:       tt.fields.close,
				idleTimeout: tt.fields.idleTimeout,
			}
			if got := c.Len(); got != tt.want {
				t.Errorf("channelPool.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}
