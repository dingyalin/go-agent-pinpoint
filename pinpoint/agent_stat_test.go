package pinpoint

import (
	"reflect"
	"testing"

	pinpoint "github.com/dingyalin/pinpoint-go-agent/thrift/dto/pinpoint"
)

func Test_getTAgentStat(t *testing.T) {
	type args struct {
		agentID         string
		startTime       int64
		collectInterval int64
	}
	tests := []struct {
		name string
		args args
		want *pinpoint.TAgentStat
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTAgentStat(tt.args.agentID, tt.args.startTime, tt.args.collectInterval); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getAgentStat() = %v, want %v", got, tt.want)
			}
		})
	}
}
