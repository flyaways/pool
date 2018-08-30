package pool

import (
	"net"
)

var (
	//FactoryTCP ...
	FactoryTCP = func() (interface{}, error) { return net.DialTimeout("tcp", address, dialTimeout) }

	//CloseTCP ...
	CloseTCP = func(v interface{}) error { return v.(net.Conn).Close() }
)
