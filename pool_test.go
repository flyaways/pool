package pool

import (
	"reflect"
	"testing"
)

func TestNewChannelPool(t *testing.T) {
	type args struct {
		poolConfig *PoolConfig
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
			got, err := NewChannelPool(tt.args.poolConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewChannelPool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewChannelPool() = %v, want %v", got, tt.want)
			}
		})
	}
}
