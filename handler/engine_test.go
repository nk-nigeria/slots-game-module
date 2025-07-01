package handler

import (
	"fmt"
	"testing"

	pb "github.com/nk-nigeria/cgp-common/proto"
	"github.com/nk-nigeria/slots-game-module/entity"
	"github.com/stretchr/testify/assert"
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
				matchState: entity.NewSlotsMathState(&pb.Match{}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &slotsEngine{}
			matrix := e.SpinMatrix(tt.args.matchState.GetMatrix())
			assert.NotEmpty(t, matrix, "Matrix should not empty ")
			e.PrintMatrix(matrix)
		})
	}
}

func Test_slotsEngine_Random(t *testing.T) {
	type args struct {
		min int
		max int
	}
	tests := []struct {
		name string
		e    *slotsEngine
		args args
		want int
	}{
		// TODO: Add test cases.
		{
			name: "Random number",
			args: args{
				min: 0,
				max: 1000,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &slotsEngine{}
			got := e.Random(tt.args.min, tt.args.max)
			t.Logf("random number %d", got)
		})
	}
}

func Test_slotsEngine_SpinMatrix(t *testing.T) {
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
			name: "Spin matrix",
			args: args{
				matchState: entity.NewSlotsMathState(nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &slotsEngine{}
			matrix := e.SpinMatrix(tt.args.matchState.GetMatrix())
			e.PrintMatrix(matrix)
		})
	}
}

func Test_slotsEngine_SpreadWildInMatrix(t *testing.T) {
	type args struct {
		matchState *entity.SlotsMatchState
	}
	tests := []struct {
		name string
		e    *slotsEngine
		args args
		want entity.SlotMatrix
	}{
		// TODO: Add test cases.
		{
			name: "Spread matrix",
			args: args{
				matchState: entity.NewSlotsMathState(nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &slotsEngine{}
			matrix := e.SpinMatrix(tt.args.matchState.GetMatrix())
			e.PrintMatrix(matrix)
			t.Log("Spread Matrix")
			spreadMatrix := e.SpreadWildInMatrix(matrix)
			e.PrintMatrix(spreadMatrix)
		})
	}
}

func Test_slotsEngine_PaylineMatrix(t *testing.T) {
	name := "Test Playline Matrix"
	t.Run(name, func(t *testing.T) {
		engine := &slotsEngine{}
		matchState := entity.NewSlotsMathState(nil)
		matrix := engine.SpinMatrix(matchState.GetMatrix())
		spreadMatrix := engine.SpreadWildInMatrix(matrix)
		engine.PrintMatrix(spreadMatrix)
		paylines := engine.PaylineMatrix(spreadMatrix)
		assert.Equal(t, matrix.Rows, len(paylines), "payline size not same row matrix")
		for idx, payline := range paylines {
			fmt.Printf("line %d symbol %d occur %d \r\n", idx, payline.Symbol.Number(), payline.NumOccur)
		}
	})
}

func Test_slotsEngine_FilterSymbol(t *testing.T) {
	name := "Test Filter Payline"
	t.Run(name, func(t *testing.T) {
		engine := &slotsEngine{}
		matchState := entity.NewSlotsMathState(nil)
		matrix := engine.SpinMatrix(matchState.GetMatrix())
		spreadMatrix := engine.SpreadWildInMatrix(matrix)
		engine.PrintMatrix(spreadMatrix)
		paylines := engine.PaylineMatrix(spreadMatrix)
		for _, payline := range engine.FilterPayline(paylines, func(numOccur int) bool {
			return numOccur >= 3
		}) {
			fmt.Printf("payline id %d symbol %d occur %d \r\n", payline.Id, payline.Symbol.Number(), payline.NumOccur)
		}
	})
}

func Test_slotsEngine_Process(t *testing.T) {
	name := "Test Enginre Process"
	t.Run(name, func(t *testing.T) {
		engine := &slotsEngine{}
		matchState := entity.NewSlotsMathState(nil)
		matrix := entity.SlotMatrix{
			List: []pb.SiXiangSymbol{
				2048, 16, 16, 128, 64,
				4095, 1, 128, 8, 64,
				16, 1, 256, 4095, 1,
			},
			Cols: 5,
			Rows: 3,
		}
		matchState.SetMatrix(matrix)
		spreadMatrix := engine.SpreadWildInMatrix(matrix)
		engine.PrintMatrix(spreadMatrix)
		paylines := engine.PaylineMatrix(spreadMatrix)
		paylinesFilter := engine.FilterPayline(paylines, func(numOccur int) bool {
			return numOccur >= 3
		})
		matchState.SetPaylines(paylinesFilter)
		chipsMcb := int64(22222)
		for _, payline := range matchState.GetPaylines() {
			payline.Rate = engine.RatioPayline(payline)
			payline.Chips = int64(payline.Rate * float64(chipsMcb))
		}
		for _, payline := range matchState.GetPaylines() {
			t.Logf("payline id %d, symbol %d occur %d ratio %v chips %d",
				payline.Id, payline.Symbol, payline.NumOccur, payline.Rate, payline.Chips)
		}
	})
}
