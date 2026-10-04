package pool

import (
	"errors"
	"io"
	"net"
	"net/rpc"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

func lateReturnOptions(target string) *Options {
	o := NewOptions()
	o.InitTargets = []string{target}
	o.InitCap = 1
	o.MaxCap = 1
	o.input = make(chan *[]string)
	return o
}

func assertPutAfterClose(t *testing.T, closePool func(), putNil, putConn func() error, want error) {
	t.Helper()
	closePool()
	closePool()
	if err := putNil(); err != errRejected {
		t.Fatalf("Put(nil) error = %v, want %v", err, errRejected)
	}
	if err := putConn(); err != want {
		t.Fatalf("Put after Close error = %v, want %v", err, want)
	}
}

func TestTCPPoolPutAfterClose(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	opts := lateReturnOptions(listener.Addr().String())
	defer close(opts.input)
	p, err := NewTCPPool(opts)
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	peer, err := listener.Accept()
	if err != nil {
		t.Fatal(err)
	}
	defer peer.Close()
	conn, err := p.Get()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	assertPutAfterClose(t, p.Close, func() error { return p.Put(nil) }, func() error { return p.Put(conn) }, nil)
	if _, err := conn.Write([]byte("x")); err == nil {
		t.Fatal("returned TCP connection is still open")
	}
}

func TestRPCPoolPutAfterClose(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	opts := lateReturnOptions(listener.Addr().String())
	defer close(opts.input)
	p, err := NewRPCPool(opts)
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	peer, err := listener.Accept()
	if err != nil {
		t.Fatal(err)
	}
	defer peer.Close()
	conn, err := p.Get()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	assertPutAfterClose(t, p.Close, func() error { return p.Put(nil) }, func() error { return p.Put(conn) }, nil)
	call := conn.Go("Unused.Method", 0, new(int), make(chan *rpc.Call, 1))
	select {
	case result := <-call.Done:
		if result.Error != rpc.ErrShutdown {
			t.Fatalf("returned RPC client error = %v, want %v", result.Error, rpc.ErrShutdown)
		}
	case <-time.After(time.Second):
		t.Fatal("returned RPC client is still open")
	}
	if err := peer.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	if _, err := peer.Read(make([]byte, 1)); err != io.EOF {
		t.Fatalf("returned RPC connection read error = %v, want EOF", err)
	}
}

func TestGRPCPoolPutAfterClose(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	server := grpc.NewServer()
	defer server.Stop()
	go server.Serve(listener)
	opts := lateReturnOptions(listener.Addr().String())
	defer close(opts.input)
	p, err := NewGRPCPool(opts, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	conn, err := p.Get()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	assertPutAfterClose(t, p.Close, func() error { return p.Put(nil) }, func() error { return p.Put(conn) }, nil)
	if state := conn.GetState(); state != connectivity.Shutdown {
		t.Fatalf("returned gRPC connection state = %v, want Shutdown", state)
	}
}

func TestPoolPutAfterCloseError(t *testing.T) {
	want := errors.New("close failed")
	t.Run("TCP", func(t *testing.T) {
		conn, peer := net.Pipe()
		defer conn.Close()
		defer peer.Close()
		calls := 0
		p := &TCPPool{
			conns: make(chan *tcpIdleConn, 1),
			close: func(got net.Conn) error {
				calls++
				if got != conn {
					t.Errorf("close received %v, want returned connection", got)
				}
				return want
			},
		}
		assertPutAfterClose(t, p.Close, func() error { return p.Put(nil) }, func() error { return p.Put(conn) }, want)
		if calls != 1 {
			t.Fatalf("close calls = %d, want 1", calls)
		}
	})
	t.Run("RPC", func(t *testing.T) {
		conn := new(rpc.Client)
		calls := 0
		p := &RPCPool{
			conns: make(chan *rpcIdleConn, 1),
			close: func(got *rpc.Client) error {
				calls++
				if got != conn {
					t.Errorf("close received %v, want returned client", got)
				}
				return want
			},
		}
		assertPutAfterClose(t, p.Close, func() error { return p.Put(nil) }, func() error { return p.Put(conn) }, want)
		if calls != 1 {
			t.Fatalf("close calls = %d, want 1", calls)
		}
	})
	t.Run("GRPC", func(t *testing.T) {
		conn := new(grpc.ClientConn)
		calls := 0
		p := &GRPCPool{
			conns: make(chan *grpcIdleConn, 1),
			close: func(got *grpc.ClientConn) error {
				calls++
				if got != conn {
					t.Errorf("close received %v, want returned connection", got)
				}
				return want
			},
		}
		assertPutAfterClose(t, p.Close, func() error { return p.Put(nil) }, func() error { return p.Put(conn) }, want)
		if calls != 1 {
			t.Fatalf("close calls = %d, want 1", calls)
		}
	})
}
