package pool

import (
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	type args struct {
		network     string
		address     string
		initialCap  int
		maxCap      int
		dialTimeout time.Duration
		idleTimeout time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    Pool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.network, tt.args.address, tt.args.initialCap, tt.args.maxCap, tt.args.dialTimeout, tt.args.idleTimeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
