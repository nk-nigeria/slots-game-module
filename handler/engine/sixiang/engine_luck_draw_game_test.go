package sixiang

import (
	"reflect"
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
		matrix1 := entity.NewMatrixLuckyDraw()
		entity.ShuffleMatrix(matrix1)
		t.Log("matrix 1")
		engine.PrintMatrix(matrix1)
		matrix2 := entity.NewMatrixLuckyDraw()
		entity.ShuffleMatrix(matrix2)
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
			trackSpinSymbol := make(map[api.SiXiangSymbol]int, 0)
			var trackSpinSymbolExpect = map[api.SiXiangSymbol]int{
				api.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_MINOR:  3,
				api.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_MAJOR:  3,
				api.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_MEGA:   3,
				api.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_GRAND:  3,
				api.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_GOLD_1: 1,
				api.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_GOLD_2: 1,
				api.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_GOLD_3: 1,
			}

			matchState.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_LUCKDRAW
			e.NewGame(matchState)

			for i := int32(0); i < int32(matchState.MatrixSpecial.Size); i++ {
				e.Process(matchState)
				for _, spin := range matchState.SpinSymbols {
					num := trackSpinSymbol[spin.Symbol]
					num++
					trackSpinSymbol[spin.Symbol] = num
				}

			}
			assert.Equal(t, len(trackSpinSymbolExpect), len(trackSpinSymbol))
			for k, v := range trackSpinSymbol {
				assert.Equal(t, trackSpinSymbolExpect[k], v)
			}
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
		matchState.MatrixSpecial = entity.NewMatrixLuckyDraw()
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
				GameReward: &pb.GameReward{
					ChipsWin: 0,
				},
				ChipsMcb: bet.Chips,
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
			matchState.MatrixSpecial = entity.NewMatrixLuckyDraw()
			engine := NewLuckyDrawEngine(
				func(min, max int) int { return min },
				func(min, max float64) float64 { return min },
			)
			engine.NewGame(matchState)
			idFlip := entity.RandomInt(0, matchState.MatrixSpecial.Size)
			matchState.MatrixSpecial.List[idFlip] = symbol
			matchState.MatrixSpecial.TrackFlip[idFlip] = true
			matchState.SetBetInfo(bet)

			matchState.SpinSymbols[0].Symbol = symbol
			test := test{
				name: "Test_luckyDrawEngine_Finish lucky symbok" + symbol.String(),

				want: pb.SlotDesk{
					CurrentSixiangGame: api.SiXiangGame_SI_XIANG_GAME_LUCKDRAW,
					NextSixiangGame:    api.SiXiangGame_SI_XIANG_GAME_LUCKDRAW,
					BigWin:             api.BigWin_BIG_WIN_UNSPECIFIED,
					WinJp:              api.WinJackpot_WIN_JACKPOT_UNSPECIFIED,
					ChipsMcb:           bet.Chips,
					GameReward: &pb.GameReward{
						ChipsWin: bet.Chips * int64(entity.ListSymbolLuckyDraw[symbol].Value.Min),
					},
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
			matchState.MatrixSpecial = entity.NewMatrixLuckyDraw()
			matchState.SetBetInfo(bet)
			engine := NewLuckyDrawEngine(
				func(min, max int) int { return min },
				func(min, max float64) float64 { return min },
			)
			engine.NewGame(matchState)
			matchState.MatrixSpecial.List[1] = api.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_GOLD_1
			matchState.MatrixSpecial.TrackFlip[1] = true
			matchState.SpinSymbols[0].Symbol = api.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_GOLD_1
			engine.Finish(matchState)

			matchState.MatrixSpecial.List[2] = symbol
			matchState.MatrixSpecial.List[5] = symbol
			matchState.MatrixSpecial.List[8] = symbol

			matchState.MatrixSpecial.TrackFlip[2] = true
			matchState.MatrixSpecial.TrackFlip[5] = true
			matchState.MatrixSpecial.TrackFlip[8] = true

			matchState.SpinSymbols[0].Symbol = symbol
			test := test{
				name: "Test_luckyDrawEngine_Finish lucky symbok" + symbol.String(),

				want: pb.SlotDesk{
					CurrentSixiangGame: api.SiXiangGame_SI_XIANG_GAME_LUCKDRAW,
					NextSixiangGame:    api.SiXiangGame_SI_XIANG_GAME_NORMAL,
					GameReward: &pb.GameReward{
						ChipsWin: bet.Chips*int64(entity.ListSymbolLuckyDraw[api.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_GOLD_1].Value.Min) +
							bet.Chips*int64(entity.ListSymbolLuckyDraw[symbol].Value.Min),
					},
					ChipsMcb: bet.Chips,
				},
			}
			test.want.BigWin, test.want.WinJp = entity.LuckySymbolToReward(symbol)
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
			assert.Equal(t, tt.want.GameReward.ChipsWin, slotDesk.GameReward.ChipsWin)
			assert.Equal(t, tt.want.ChipsMcb, slotDesk.ChipsMcb)
		})
	}
}

func Test_luckyDrawEngine_GetNextSiXiangGame(t *testing.T) {
	type fields struct {
		randomIntFn   func(min, max int) int
		randomFloat64 func(min, max float64) float64
	}
	type args struct {
		s *entity.SlotsMatchState
	}
	type test struct {
		name   string
		fields fields
		args   args
		want   pb.SiXiangGame
	}
	e := NewLuckyDrawEngine(func(min, max int) int { return min }, func(min, max float64) float64 { return min })
	engine := e.(*luckyDrawEngine)
	tests := make([]test, 0)
	{
		arr := []pb.SiXiangSymbol{
			api.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_GOLD_1,
			api.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_GOLD_2,
			api.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_GOLD_3,
		}
		for _, sym := range arr {
			test := test{
				name: "luckyDrawEngine_GetNextSiXiangGame not finish",
				want: api.SiXiangGame_SI_XIANG_GAME_LUCKDRAW,
			}
			matchState := entity.NewSlotsMathState(nil)
			matchState.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_LUCKDRAW
			engine.NewGame(matchState)
			test.args = args{
				s: matchState,
			}
			for id, symbol := range matchState.MatrixSpecial.List {
				if symbol == sym {
					matchState.MatrixSpecial.TrackFlip[id] = true
				}
			}
			trackJP := make(map[api.SiXiangSymbol]int)
			for id, symbol := range matchState.MatrixSpecial.List {
				if int(symbol) < int(pb.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_GOLD_1) {
					num := trackJP[symbol]
					if num >= 2 {
						continue
					}
					matchState.MatrixSpecial.TrackFlip[id] = true
					num++
					trackJP[symbol] = num
				}
			}
			tests = append(tests, test)
		}
	}
	{
		arr := []pb.SiXiangSymbol{
			api.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_MINOR,
			api.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_MINOR,
			api.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_MINOR,
		}
		for _, sym := range arr {
			test := test{
				name: "luckyDrawEngine_GetNextSiXiangGame finish",
				want: api.SiXiangGame_SI_XIANG_GAME_NORMAL,
			}
			matchState := entity.NewSlotsMathState(nil)
			matchState.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_LUCKDRAW
			engine.NewGame(matchState)
			test.args = args{
				s: matchState,
			}
			for id, symbol := range matchState.MatrixSpecial.List {
				if symbol == sym {
					matchState.MatrixSpecial.TrackFlip[id] = true
				}
			}
			trackJP := make(map[api.SiXiangSymbol]int)
			for id, symbol := range matchState.MatrixSpecial.List {
				if int(symbol) >= int(pb.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_GOLD_1) {
					num := trackJP[symbol]
					if num >= 2 {
						continue
					}
					matchState.MatrixSpecial.TrackFlip[id] = true
					num++
					trackJP[symbol] = num
				}
			}
			tests = append(tests, test)
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := engine.GetNextSiXiangGame(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("luckyDrawEngine.GetNextSiXiangGame() = %v, want %v", got, tt.want)
			}
		})
	}
}
