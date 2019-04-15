package pool

import (
	"net"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestTCPPool_Get(t *testing.T) {
	type fields struct {
		Mu          sync.Mutex
		IdleTimeout time.Duration
		conns       chan *tcpIdleConn
		factory     func() (net.Conn, error)
		close       func(net.Conn) error
	}
	tests := []struct {
		name    string
		fields  fields
		want    net.Conn
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &TCPPool{
				Mu:          tt.fields.Mu,
				IdleTimeout: tt.fields.IdleTimeout,
				conns:       tt.fields.conns,
				factory:     tt.fields.factory,
				close:       tt.fields.close,
			}
			got, err := c.Get()
			if (err != nil) != tt.wantErr {
				t.Errorf("TCPPool.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TCPPool.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTCPPool(t *testing.T) {
	type args struct {
		c *Config
	}
	tests := []struct {
		name    string
		args    args
		want    *TCPPool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTCPPool(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTCPPool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTCPPool() = %v, want %v", got, tt.want)
			}
		})
	}
}
