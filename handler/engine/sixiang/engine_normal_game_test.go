package sixiangengine

import (
	"fmt"
	"testing"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
	"github.com/stretchr/testify/assert"
)

func Test_normalEngine_InitMatrix(t *testing.T) {
	type args struct {
		matchState *entity.SixiangMatchState
	}
	tests := []struct {
		name string
		e    *normalEngine
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
			e := &normalEngine{}
			matrix := e.SpinMatrix(tt.args.matchState.GetMatrix())
			assert.NotEmpty(t, matrix, "Matrix should not empty ")
			e.PrintMatrix(matrix)
		})
	}
}

func Test_normalEngine_Random(t *testing.T) {
	type args struct {
		min int
		max int
	}
	tests := []struct {
		name string
		e    *normalEngine
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
			e := &normalEngine{}
			got := e.Random(tt.args.min, tt.args.max)
			t.Logf("random number %d", got)
		})
	}
}

func Test_normalEngine_SpinMatrix(t *testing.T) {
	type args struct {
		matchState *entity.SixiangMatchState
	}
	tests := []struct {
		name string
		e    *normalEngine
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
			e := &normalEngine{}
			matrix := e.SpinMatrix(tt.args.matchState.GetMatrix())
			e.PrintMatrix(matrix)
		})
	}
}

func Test_normalEngine_SpreadWildInMatrix(t *testing.T) {
	type args struct {
		matchState *entity.SixiangMatchState
	}
	tests := []struct {
		name string
		e    *normalEngine
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
			e := &normalEngine{}
			matrix := e.SpinMatrix(tt.args.matchState.GetMatrix())
			e.PrintMatrix(matrix)
			t.Log("Spread Matrix")
			spreadMatrix := e.SpreadWildInMatrix(matrix)
			e.PrintMatrix(spreadMatrix)
		})
	}
}

func Test_normalEngine_PaylineMatrix(t *testing.T) {
	name := "Test Playline Matrix"
	t.Run(name, func(t *testing.T) {
		engine := &normalEngine{}
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

func Test_normalEngine_FilterSymbol(t *testing.T) {
	name := "Test Filter Payline"
	t.Run(name, func(t *testing.T) {
		engine := &normalEngine{}
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

func Test_normalEngine_Process(t *testing.T) {
	name := "Test Enginre Process"
	t.Run(name, func(t *testing.T) {
		engine := &normalEngine{}
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

func Test_normalEngine_Process_2(t *testing.T) {
	name := "Test Enginre Process"
	t.Run(name, func(t *testing.T) {
		engine := &normalEngine{}
		matchState := entity.NewSlotsMathState(nil)
		matchState.SetBetInfo(&pb.InfoBet{
			Chips:       11,
			ReqSpecGame: int32(pb.SiXiangGame_SI_XIANG_GAME_LUCKDRAW),
		})
		engine.NewGame(matchState)
		engine.Process(matchState)
		result, _ := engine.Finish(matchState)
		engine.PrintMatrix(matchState.GetMatrix())
		slotdesk := result.(*pb.SlotDesk)
		t.Logf("%v", slotdesk)
	})
}
