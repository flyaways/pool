package pool

import (
	"bufio"
	"encoding/gob"
	"io"
	"net/rpc"
	"testing"
	"time"
)

func TestGobCoDec_WriteRequest(t *testing.T) {
	type fields struct {
		Timeout time.Duration
		Closer  io.ReadWriteCloser
		Decoder *gob.Decoder
		Encoder *gob.Encoder
		EncBuf  *bufio.Writer
	}
	type args struct {
		r    *rpc.Request
		body interface{}
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
			c := &GobCoDec{
				Timeout: tt.fields.Timeout,
				Closer:  tt.fields.Closer,
				Decoder: tt.fields.Decoder,
				Encoder: tt.fields.Encoder,
				EncBuf:  tt.fields.EncBuf,
			}
			if err := c.WriteRequest(tt.args.r, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("GobCoDec.WriteRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGobCoDec_timeoutCoder(t *testing.T) {
	type fields struct {
		Timeout time.Duration
		Closer  io.ReadWriteCloser
		Decoder *gob.Decoder
		Encoder *gob.Encoder
		EncBuf  *bufio.Writer
	}
	type args struct {
		e   interface{}
		msg string
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
			c := &GobCoDec{
				Timeout: tt.fields.Timeout,
				Closer:  tt.fields.Closer,
				Decoder: tt.fields.Decoder,
				Encoder: tt.fields.Encoder,
				EncBuf:  tt.fields.EncBuf,
			}
			if err := c.timeoutCoder(tt.args.e, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("GobCoDec.timeoutCoder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
