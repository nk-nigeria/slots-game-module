package engine

import (
	"testing"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	api "github.com/ciaolink-game-platform/cgp-common/proto"
	"github.com/stretchr/testify/assert"
)

func TestRandomFloat64(t *testing.T) {
	type args struct {
		min float64
		max float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
		{
			name: "RandomFloat64 < 1.0",
			args: args{
				min: 0.1,
				max: 0.5,
			},
			want: 0.5,
		},
		{
			name: "RandomFloat64 > 1.0",
			args: args{
				min: 1,
				max: 5,
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RandomFloat64(tt.args.min, tt.args.max)
			if got > tt.want {
				t.Errorf("RandomFloat64() = %v, want %v", got, tt.want)
			}
			t.Logf("random %v", got)
		})
	}
}

func TestShuffleSlice(t *testing.T) {
	type args struct {
		slice []api.SiXiangSymbol
	}
	tests := []struct {
		name string
		args args
		want []api.SiXiangSymbol
	}{
		// TODO: Add test cases.
		{
			name: "TestShuffleSlice",
			args: args{
				slice: entity.ListEyeSiXiang[:],
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ShuffleSlice(tt.args.slice)
			assert.Equal(t, len(got), len(tt.args.slice))
			assert.Equal(t, got, tt.args.slice)
			t.Log(got)
		})
	}
}
