package pool

import "google.golang.org/grpc"

var (

	//DialOptions ...
	dialOptions = []grpc.DialOption{}

	//FactoryGRPC ...
	FactoryGRPC = func() (interface{}, error) {
		return grpc.Dial(address, dialOptions...)
	}

	//CloseGRPC ...
	CloseGRPC = func(v interface{}) error { return v.(*grpc.ClientConn).Close() }
)
