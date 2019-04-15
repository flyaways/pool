package pool

import (
	"reflect"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
)

func TestNewGRPCPool(t *testing.T) {
	type args struct {
		o           *Options
		dialOptions []grpc.DialOption
	}
	tests := []struct {
		name    string
		args    args
		want    *GRPCPool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGRPCPool(tt.args.o, tt.args.dialOptions...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGRPCPool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGRPCPool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGRPCPool_Get(t *testing.T) {
	type fields struct {
		Mu          sync.Mutex
		IdleTimeout time.Duration
		conns       chan *grpcIdleConn
		factory     func() (*grpc.ClientConn, error)
		close       func(*grpc.ClientConn) error
	}
	tests := []struct {
		name    string
		fields  fields
		want    *grpc.ClientConn
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &GRPCPool{
				Mu:          tt.fields.Mu,
				IdleTimeout: tt.fields.IdleTimeout,
				conns:       tt.fields.conns,
				factory:     tt.fields.factory,
				close:       tt.fields.close,
			}
			got, err := c.Get()
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCPool.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCPool.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
