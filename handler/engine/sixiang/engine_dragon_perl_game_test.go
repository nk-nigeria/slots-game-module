package sixiang

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/nakamaFramework/cgb-slots-game-module/entity"
	api "github.com/nakamaFramework/cgp-common/proto"
	pb "github.com/nakamaFramework/cgp-common/proto"
	"github.com/stretchr/testify/assert"
)

func Test_dragonPearlEngine_NewGame(t *testing.T) {

	type args struct {
		matchState *entity.SlotsMatchState
	}
	matchState := entity.NewSlotsMathState(nil)
	matchState.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_DRAGON_PEARL
	matchStateExpect := entity.NewSlotsMathState(nil)
	matchStateExpect.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_DRAGON_PEARL
	matchStateExpect.NumSpinLeft = defaultDragonPearlGemSpin
	{
		list := make([]api.SiXiangSymbol, 0)
		for k := range entity.ListEyeSiXiang {
			list = append(list, k)
		}
		// matchStateExpect.EyeSymbolRemains = entity.ShuffleSlice(list)
	}
	matrixSpecial := entity.NewMatrixDragonPearl()
	matchStateExpect.MatrixSpecial = &matrixSpecial
	tests := []struct {
		name string
		args args
		want *entity.SlotsMatchState
	}{
		// TODO: Add test cases.
		{
			name: "Test init new dragon pearl game",
			args: args{
				matchState: matchState,
			},
			want: matchStateExpect,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewDragonPearlEngine(4, nil, nil)
			got, err := e.NewGame(tt.args.matchState)
			assert.NotNil(t, got)
			assert.NoError(t, err)
			assert.Equal(t, tt.want.CurrentSiXiangGame, tt.args.matchState.CurrentSiXiangGame)
			assert.Equal(t, tt.want.NumSpinLeft, tt.args.matchState.NumSpinLeft)
			// assert.Equal(t, tt.want.EyeSymbolRemains, tt.args.matchState.EyeSymbolRemains)
			// assert.Equal(t,
			// 	tt.want.CollectionSymbol[tt.args.matchState.CurrentSiXiangGame][int(tt.args.matchState.Bet().Chips)],
			// 	tt.args.matchState.CollectionSymbol[tt.args.matchState.CurrentSiXiangGame][int(tt.args.matchState.Bet().GetChips())])
			trackSym := make(map[api.SiXiangSymbol]int)
			trackSymExpect := make(map[api.SiXiangSymbol]int)
			for _, sym := range tt.want.MatrixSpecial.List {
				num := trackSymExpect[sym]
				num++
				trackSymExpect[sym] = num
			}
			for _, sym := range tt.args.matchState.MatrixSpecial.List {
				num := trackSym[sym]
				num++
				trackSym[sym] = num
			}
			for k, v := range trackSymExpect {
				assert.Equal(t, v, trackSym[k])
			}
		})
	}
}

func Test_dragonPearlEngine_Process(t *testing.T) {
	name := "Test_dragonPearlEngine_Process check spin num symbol not over maximum appear"
	trackSymbolSpin := make(map[api.SiXiangSymbol]int)
	t.Run(name, func(t *testing.T) {

		e := NewDragonPearlEngine(4, nil, nil)
		matchState := entity.NewSlotsMathState(nil)
		matchState.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_DRAGON_PEARL
		trackSymbolSpinExpect := make(map[api.SiXiangSymbol]int)
		trackSymbolSpinExpect[api.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_BIRD] = 1
		trackSymbolSpinExpect[api.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_DRAGON] = 1
		trackSymbolSpinExpect[api.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_TIGER] = 1
		trackSymbolSpinExpect[api.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_WARRIOR] = 1
		trackSymbolSpinExpect[api.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_GEM_RANDOM1] = 4
		trackSymbolSpinExpect[api.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_GEM_RANDOM2] = 3
		trackSymbolSpinExpect[api.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_GEM_RANDOM3] = 2
		trackSymbolSpinExpect[api.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_GEM_RANDOM4] = 2
		trackSymbolSpinExpect[api.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_GEM_RANDOM5] = 1
		e.NewGame(matchState)
		matchState.Bet().Chips = 1000
		for {
			gemSpinRemain := 3
			matchState.NumSpinLeft = gemSpinRemain
			got, err := e.Process(matchState)
			assert.NotNil(t, got)
			if err != nil {
				t.Log(err)
				break
			}
			for _, spin := range matchState.SpinSymbols {
				num := trackSymbolSpin[spin.Symbol]
				num++
				trackSymbolSpin[spin.Symbol] = num
			}
			spin := matchState.SpinSymbols[0]
			msg := fmt.Sprintf("symbol %s", spin.Symbol)
			switch spin.Symbol {
			case api.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_BIRD:
				assert.Equal(t, gemSpinRemain-1+3, matchState.NumSpinLeft, msg)
			case api.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_DRAGON: //todo
				assert.Equal(t, gemSpinRemain, matchState.NumSpinLeft, msg)
			case api.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_TIGER:
				assert.Equal(t, gemSpinRemain, matchState.NumSpinLeft, msg)
			case api.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_WARRIOR:
				assert.Equal(t, 4, len(matchState.SpinSymbols), msg)
				assert.Equal(t, gemSpinRemain, matchState.NumSpinLeft, msg)
			default:
				assert.Equal(t, 1, len(matchState.SpinSymbols), msg)
				assert.Equal(t, gemSpinRemain-1, matchState.NumSpinLeft, msg)
			}
			t.Logf("len trackflip %d", len(matchState.MatrixSpecial.TrackFlip))
			if len(matchState.MatrixSpecial.TrackFlip) >= 15 {
				break
			}
			// }
		}
		for k, v := range trackSymbolSpin {
			t.Logf("symbol %s, spin %d time", k.String(), v)
			assert.Equal(t, v, trackSymbolSpinExpect[k])
		}
	})
}

func Test_dragonPearlEngine_Process_CheckMinMaxEyeFlip(t *testing.T) {
	name := "Test_dragonPearlEngine_Process_CheckMinMaxEyeFlip"
	t.Run(name, func(t *testing.T) {
		for i := 0; i < 10000; i++ {
			e := NewDragonPearlEngine(4, nil, nil)
			matchState := entity.NewSlotsMathState(nil)
			matchState.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_DRAGON_PEARL
			matchState.Bet().Chips = 1000
			e.NewGame(matchState)
			for {
				e.Process(matchState)
				res, err := e.Finish(matchState)
				assert.NoError(t, err)
				assert.NotNil(t, res)
				result, ok := res.(*api.SlotDesk)
				assert.Equal(t, true, ok)
				if result.IsFinishGame {
					numEye := matchState.MatrixSpecial.CountSymbolCond(func(index int, symbol api.SiXiangSymbol) bool {
						return entity.IsSixiangEyeSymbol(symbol) && matchState.MatrixSpecial.IsFlip(index)
					})
					assert.LessOrEqual(t, 1, numEye)
					assert.GreaterOrEqual(t, 3, numEye)
					break
				}
			}
		}
	})
}
func Test_dragonPearlEngine_Finish(t *testing.T) {
	type want struct {
		slotDesk   *api.SlotDesk
		matchState *entity.SlotsMatchState
	}
	type test struct {
		name    string
		args    *entity.SlotsMatchState
		want    want
		wantErr bool
	}
	engine := NewDragonPearlEngine(4, func(min, max int) int { return min }, func(min, max float64) float64 { return min })

	tests := make([]test, 0)
	// test gem random
	{
		arrGem := []api.SiXiangSymbol{
			api.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_GEM_RANDOM1,
			api.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_GEM_RANDOM2,
			api.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_GEM_RANDOM3,
			api.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_GEM_RANDOM4,
			api.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_GEM_RANDOM5,
		}
		for idx, gem := range arrGem {
			test := test{
				name: "Test finish random gem " + strconv.Itoa(idx+1),
			}
			matchState := entity.NewSlotsMathState(nil)
			engine.NewGame(matchState)
			matchState.SetBetInfo(&api.InfoBet{
				Chips: 100,
			})
			matchState.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_DRAGON_PEARL
			matchState.NumSpinLeft = 0
			idFlip := 0
			matchState.MatrixSpecial.ForEeach(func(idx, row, col int, symbol api.SiXiangSymbol) {
				if symbol == gem {
					idFlip = idx
				}
			})
			matchState.MatrixSpecial.Flip(idFlip)
			matchState.SpinSymbols = []*api.SpinSymbol{{
				Symbol: gem,
			}}
			slotDesk := &api.SlotDesk{
				CurrentSixiangGame: api.SiXiangGame_SI_XIANG_GAME_DRAGON_PEARL,
				NextSixiangGame:    api.SiXiangGame_SI_XIANG_GAME_NORMAL,
				IsFinishGame:       true,
				ChipsMcb:           matchState.Bet().GetChips(),
			}
			ratioBonus := float64(1)
			// ml := matchState.CollectionSymbolToSlice(matchState.CurrentSiXiangGame, int(matchState.Bet().Chips))
			// for _, eyeSym := range ml {
			// 	r := entity.ListEyeSiXiang[eyeSym.Symbol].Value.Min
			// 	if float64(r) > ratioBonus {
			// 		ratioBonus = float64(r)
			// 	}
			// }
			slotDesk.GameReward.ChipsWin = int64(ratioBonus * float64(slotDesk.ChipsMcb) * float64(entity.ListSymbolDragonPearl[gem].Value.Min))
			test.args = matchState
			test.want.slotDesk = slotDesk
			test.want.matchState = entity.NewSlotsMathState(nil)
			tests = append(tests, test)
		}
	}
	// test gem special, eye symbol
	{

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.Finish(tt.args)
			assert.NoError(t, err)
			assert.NotNil(t, result)
			slotDesk := result.(*api.SlotDesk)
			assert.Equal(t, tt.want.slotDesk.CurrentSixiangGame, slotDesk.CurrentSixiangGame)
			assert.Equal(t, tt.want.slotDesk.NextSixiangGame, slotDesk.NextSixiangGame)
			assert.Equal(t, tt.want.slotDesk.ChipsMcb, slotDesk.ChipsMcb)
			assert.Equal(t, tt.want.slotDesk.IsFinishGame, slotDesk.IsFinishGame)
			assert.Equal(t, tt.want.slotDesk.GameReward.ChipsWin, slotDesk.GameReward.ChipsWin)
			assert.Equal(t, 15, len(slotDesk.Matrix.Lists))
		})
	}
}

func Test_dragonPearlEngine_FinishWithJpSymbol(t *testing.T) {
	name := "Test_dragonPearlEngine_FinishWithJpSymbol"
	t.Run(name, func(t *testing.T) {
		e := NewDragonPearlEngine(4, nil, nil)
		matchState := entity.NewSlotsMathState(nil)
		matchState.Bet().Chips = 1000
		matchState.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_DRAGON_PEARL
		e.NewGame(matchState)
		e.Process(matchState)
		matchState.SpinSymbols[0].Symbol = api.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_JP_MEGA
		matchState.SpinSymbols[0].WinJp = api.WinJackpot_WIN_JACKPOT_MEGA
		idx := matchState.SpinSymbols[0].Index
		matchState.SpinList[idx].Symbol = matchState.SpinSymbols[0].Symbol
		matchState.SpinList[idx].WinJp = matchState.SpinSymbols[0].WinJp
		res, err := e.Finish(matchState)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		result, ok := res.(*api.SlotDesk)
		assert.Equal(t, true, ok)
		assert.Less(t, int64(0), result.GameReward.ChipsWin)
		assert.Less(t, int64(0), result.GameReward.TotalChipsWinByGame)

		assert.LessOrEqual(t, matchState.Bet().Chips*int64(pb.WinJackpot_WIN_JACKPOT_MEGA.Number()),
			result.GameReward.ChipsWin)
		assert.LessOrEqual(t, matchState.Bet().Chips*int64(pb.WinJackpot_WIN_JACKPOT_MEGA.Number()),
			result.GameReward.TotalChipsWinByGame)

	})
}

func Test_dragonPearlEngine_EyeTiger(t *testing.T) {
	name := "Test_dragonPearlEngine_EyeTiger"
	t.Run(name, func(t *testing.T) {
		e := NewDragonPearlEngine(4, nil, nil)
		matchState := entity.NewSlotsMathState(nil)
		matchState.Bet().Chips = 1000
		matchState.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_DRAGON_PEARL
		e.NewGame(matchState)
		e.Process(matchState)
		matchState.SpinSymbols[0].Symbol = api.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_TIGER
		matchState.SpinSymbols[0].WinJp = api.WinJackpot_WIN_JACKPOT_UNSPECIFIED
		idx := matchState.SpinSymbols[0].Index
		matchState.SpinList[idx].Symbol = matchState.SpinSymbols[0].Symbol
		matchState.SpinList[idx].WinJp = matchState.SpinSymbols[0].WinJp
		assert.Less(t, int64(0), matchState.LastResult.GameReward.TotalChipsWinByGame)
		lastResult := *matchState.LastResult
		res, err := e.Finish(matchState)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		result, ok := res.(*api.SlotDesk)
		t.Logf("prev ratio %v, current ratio %v",
			lastResult.GameReward.TotalRatioWin,
			result.GameReward.TotalRatioWin,
		)
		assert.Equal(t, true, ok)
		assert.Less(t, int64(0), result.GameReward.TotalChipsWinByGame)
		assert.LessOrEqual(t, int64(float64(lastResult.GameReward.TotalChipsWinByGame)*1.99),
			result.GameReward.TotalChipsWinByGame)
	})
}
