package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
)

func Test_slotsEngine_InitMatrix(t *testing.T) {
	type args struct {
		matchState *entity.SlotsMatchState
	}
	tests := []struct {
		name string
		e    *slotsEngine
		args args
	}{
		// TODO: Add test cases.
		{
			name: "Test Init Matrix",
			args: args{
				matchState: entity.NewSlotsMathState(&lib.MatchLabel{}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &slotsEngine{}
			e.SpinMatrix(tt.args.matchState)
			matrix, _, _ := tt.args.matchState.GetMatrix()
			assert.NotEmpty(t, matrix, "Matrix should not empty ")
			e.PrintMatrix(tt.args.matchState)
		})
	}
}
