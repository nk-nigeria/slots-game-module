package engine

import (
	"testing"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	api "github.com/ciaolink-game-platform/cgp-common/proto"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
	"github.com/stretchr/testify/assert"
)

func Test_nluckyDrawEngine_ShuffleMatrix(t *testing.T) {
	name := "Engine_ShuffleMatrix"
	t.Run(name, func(t *testing.T) {
		engine := &luckyDrawEngine{}
		matrix1 := entity.NewSiXiangMatrixLuckyDraw()
		ShuffleMatrix(matrix1)
		t.Log("matrix 1")
		engine.PrintMatrix(matrix1)
		matrix2 := entity.NewSiXiangMatrixLuckyDraw()
		ShuffleMatrix(matrix2)
		t.Log("matrix 2")
		engine.PrintMatrix(matrix2)
		countSameSymbol := 0
		matrix1.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
			if symbol == matrix2.List[idx] {
				countSameSymbol++
			}
		})
		assert.NotEqual(t, len(matrix1.List), countSameSymbol)
	})
}

func Test_luckyDrawEngine_Process(t *testing.T) {
	type args struct {
		matchState interface{}
	}
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{
			name: "Test_luckyDrawEngine_Process",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewLuckyDrawEngine(
				func(min, max int) int { return min },
				func(min, max float64) float64 { return min },
			)
			matchState := entity.NewSlotsMathState(nil)
			mapTrackDraw := make(map[api.SiXiangSymbol]bool, 0)
			matchState.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_LUCKDRAW
			e.NewGame(matchState)
			for i := int32(0); i < int32(matchState.MatrixSpecial.Size); i++ {
				matchState.MatrixSpecial.List[i] = api.SiXiangSymbol(i)
			}
			for i := int32(0); i < int32(matchState.MatrixSpecial.Size); i++ {
				e.Process(matchState)
				if mapTrackDraw[matchState.SpinSymbol.Symbol] {
					t.Fatalf("draw duplicate symbol")
					break
				}
				mapTrackDraw[api.SiXiangSymbol(matchState.SpinSymbol.Symbol)] = true
			}
			t.Log("Done")
		})
	}
}

func Test_luckyDrawEngine_Finish(t *testing.T) {
	type fields struct {
		engine lib.Engine
	}
	type args struct {
		matchState interface{}
	}
	type test struct {
		name   string
		fields fields
		args   args
		want   api.SlotDesk
	}
	tests := make([]test, 0)
	// not spin
	{
		bet := &api.InfoBet{
			Chips: 1,
		}
		matchState := entity.NewSlotsMathState(nil)
		matchState.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_LUCKDRAW
		matchState.MatrixSpecial = entity.NewSiXiangMatrixLuckyDraw()
		matchState.SetBetInfo(bet)
		engine := NewLuckyDrawEngine(
			func(min, max int) int { return min },
			func(min, max float64) float64 { return min },
		)
		engine.NewGame(matchState)
		test := test{
			name: "Test_luckyDrawEngine_Finish not done",

			want: pb.SlotDesk{
				CurrentSixiangGame: api.SiXiangGame_SI_XIANG_GAME_LUCKDRAW,
				NextSixiangGame:    api.SiXiangGame_SI_XIANG_GAME_LUCKDRAW,
				BigWin:             api.BigWin_BIG_WIN_UNSPECIFIED,
				WinJp:              api.WinJackpot_WIN_JACKPOT_UNSPECIFIED,
				ChipsWinInSpin:     0,
				ChipsMcb:           bet.Chips,
			},
		}
		test.args.matchState = matchState
		test.fields.engine = engine
		tests = append(tests, test)
	}
	// 	// spin money gold1, gold2, gold3
	{
		arr := []api.SiXiangSymbol{api.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_GOLD_1,
			api.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_GOLD_2,
			api.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_GOLD_3}
		for _, symbol := range arr {
			bet := &api.InfoBet{
				Chips: 1,
			}
			matchState := entity.NewSlotsMathState(nil)
			matchState.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_LUCKDRAW
			matchState.MatrixSpecial = entity.NewSiXiangMatrixLuckyDraw()
			engine := NewLuckyDrawEngine(
				func(min, max int) int { return min },
				func(min, max float64) float64 { return min },
			)
			engine.NewGame(matchState)
			idFlip := RandomInt(0, matchState.MatrixSpecial.Size)
			matchState.MatrixSpecial.List[idFlip] = symbol
			matchState.MatrixSpecial.TrackFlip[idFlip] = true
			matchState.SetBetInfo(bet)

			matchState.SpinSymbol.Symbol = symbol
			test := test{
				name: "Test_luckyDrawEngine_Finish lucky symbok" + symbol.String(),

				want: pb.SlotDesk{
					CurrentSixiangGame: api.SiXiangGame_SI_XIANG_GAME_LUCKDRAW,
					NextSixiangGame:    api.SiXiangGame_SI_XIANG_GAME_LUCKDRAW,
					BigWin:             api.BigWin_BIG_WIN_UNSPECIFIED,
					WinJp:              api.WinJackpot_WIN_JACKPOT_UNSPECIFIED,
					ChipsWinInSpin:     bet.Chips * int64(entity.ListSymbolLuckyDraw[symbol].Value.Min),
					ChipsMcb:           bet.Chips,
				},
			}
			test.args.matchState = matchState
			test.fields.engine = engine
			tests = append(tests, test)
		}
	}
	// 	spin jackpot
	{
		arr := []api.SiXiangSymbol{api.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_MINOR,
			api.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_MAJOR,
			api.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_MEGA,
			api.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_GRAND}
		for _, symbol := range arr {
			bet := &api.InfoBet{
				Chips: 1,
			}
			matchState := entity.NewSlotsMathState(nil)
			matchState.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_LUCKDRAW
			matchState.MatrixSpecial = entity.NewSiXiangMatrixLuckyDraw()
			matchState.SetBetInfo(bet)
			engine := NewLuckyDrawEngine(
				func(min, max int) int { return min },
				func(min, max float64) float64 { return min },
			)
			engine.NewGame(matchState)
			matchState.MatrixSpecial.List[1] = api.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_GOLD_1
			matchState.MatrixSpecial.TrackFlip[1] = true
			matchState.SpinSymbol.Symbol = api.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_GOLD_1
			engine.Finish(matchState)

			matchState.MatrixSpecial.List[2] = symbol
			matchState.MatrixSpecial.List[5] = symbol
			matchState.MatrixSpecial.List[8] = symbol

			matchState.MatrixSpecial.TrackFlip[2] = true
			matchState.MatrixSpecial.TrackFlip[5] = true
			matchState.MatrixSpecial.TrackFlip[8] = true

			matchState.SpinSymbol.Symbol = symbol
			test := test{
				name: "Test_luckyDrawEngine_Finish lucky symbok" + symbol.String(),

				want: pb.SlotDesk{
					CurrentSixiangGame: api.SiXiangGame_SI_XIANG_GAME_LUCKDRAW,
					NextSixiangGame:    api.SiXiangGame_SI_XIANG_GAME_NOMAL,
					ChipsWinInSpin: bet.Chips*int64(entity.ListSymbolLuckyDraw[api.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_GOLD_1].Value.Min) +
						bet.Chips*int64(entity.ListSymbolLuckyDraw[symbol].Value.Min),
					ChipsMcb: bet.Chips,
				},
			}
			test.want.BigWin, test.want.WinJp = LuckySymbolToReward(symbol)
			test.args.matchState = matchState
			test.fields.engine = engine
			tests = append(tests, test)
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := tt.fields.engine
			result, err := engine.Finish(tt.args.matchState)
			assert.NoError(t, err)
			assert.NotNil(t, result)
			slotDesk := result.(*api.SlotDesk)
			assert.Equal(t, tt.want.CurrentSixiangGame, slotDesk.CurrentSixiangGame)
			assert.Equal(t, tt.want.NextSixiangGame, slotDesk.NextSixiangGame)
			assert.Equal(t, tt.want.BigWin, slotDesk.BigWin)
			assert.Equal(t, tt.want.WinJp, slotDesk.WinJp)
			assert.Equal(t, tt.want.ChipsWinInSpin, slotDesk.ChipsWinInSpin)
			assert.Equal(t, tt.want.ChipsMcb, slotDesk.ChipsMcb)
		})
	}
}
