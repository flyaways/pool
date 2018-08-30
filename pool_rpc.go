package pool

import (
	"bufio"
	"encoding/gob"
	"net"
	"net/rpc"
)

var (
	//FactoryRPC ...
	FactoryRPC = func() (interface{}, error) {
		conn, err := net.DialTimeout("tcp", address, dialTimeout)
		if err != nil {
			return nil, err
		}

		encBuf := bufio.NewWriter(conn)
		c := rpc.NewClientWithCodec(&GobCoDec{
			Closer:  conn,
			Decoder: gob.NewDecoder(conn),
			Encoder: gob.NewEncoder(encBuf),
			EncBuf:  encBuf,
			Timeout: dialTimeout,
		})

		return c, err
	}

	//CloseRPC ...
	CloseRPC = func(v interface{}) error { return v.(*rpc.Client).Close() }
)
