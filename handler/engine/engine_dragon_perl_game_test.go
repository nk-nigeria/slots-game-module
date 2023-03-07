package engine

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	api "github.com/ciaolink-game-platform/cgp-common/proto"
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
	matchStateExpect.GemSpin = defaultDragonPearlGemSpin
	matchStateExpect.EyeSiXiangRemain = entity.ListEyeSiXiang[:]
	matchStateExpect.EyeSiXiangSpined = make([]api.SiXiangSymbol, 0)
	matchStateExpect.MatrixSpecial = entity.NewSiXiangMatrixDragonPearl()
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
			e := NewDragonPearlEngine(nil, nil)
			got, err := e.NewGame(tt.args.matchState)
			assert.NotNil(t, got)
			assert.NoError(t, err)
			assert.Equal(t, tt.want.CurrentSiXiangGame, tt.args.matchState.CurrentSiXiangGame)
			assert.Equal(t, tt.want.GemSpin, tt.args.matchState.GemSpin)
			assert.Equal(t, tt.want.RatioBonus, tt.args.matchState.RatioBonus)
			assert.Equal(t, tt.want.EyeSiXiangRemain, tt.args.matchState.EyeSiXiangRemain)
			assert.Equal(t, tt.want.EyeSiXiangSpined, tt.args.matchState.EyeSiXiangSpined)
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

		e := NewDragonPearlEngine(nil, nil)
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
		ratioBonus := int64(1)
		for {
			gemSpinRemain := int64(3)
			matchState.GemSpin = gemSpinRemain
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
				assert.Equal(t, gemSpinRemain-1+3, matchState.GemSpin, msg)
			case api.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_DRAGON: //todo
				assert.Equal(t, gemSpinRemain, matchState.GemSpin, msg)
			case api.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_TIGER:
				ratioBonus = 2
				assert.Equal(t, gemSpinRemain, matchState.GemSpin, msg)
				assert.Equal(t, ratioBonus, matchState.RatioBonus, msg)
			case api.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_WARRIOR:
				assert.Equal(t, 4, len(matchState.SpinSymbols), msg)
				assert.Equal(t, gemSpinRemain, matchState.GemSpin, msg)
			default:
				assert.Equal(t, 1, len(matchState.SpinSymbols), msg)
				assert.Equal(t, gemSpinRemain-1, matchState.GemSpin, msg)
			}
			// }
			assert.Equal(t, ratioBonus, matchState.RatioBonus)
		}
		for k, v := range trackSymbolSpin {
			t.Logf("symbol %s, spin %d time", k.String(), v)
			assert.Equal(t, v, trackSymbolSpinExpect[k])
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
	engine := NewDragonPearlEngine(func(min, max int) int { return min }, func(min, max float64) float64 { return min })

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
			matchState.GemSpin = 0
			matchState.RatioBonus = 2
			matchState.SpinSymbols = []*api.SpinSymbol{{
				Symbol: gem,
			}}
			slotDesk := &api.SlotDesk{
				CurrentSixiangGame: api.SiXiangGame_SI_XIANG_GAME_DRAGON_PEARL,
				NextSixiangGame:    api.SiXiangGame_SI_XIANG_GAME_NORMAL,
				IsFinishGame:       true,
				ChipsMcb:           matchState.GetBetInfo().GetChips(),
			}
			slotDesk.ChipsWin = int64(float64(matchState.RatioBonus) * float64(slotDesk.ChipsMcb) * float64(entity.ListSymbolDragonPearl[gem].Value.Min))
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
			assert.Equal(t, tt.want.slotDesk.ChipsWin, slotDesk.ChipsWin)
			assert.Equal(t, 15, len(slotDesk.Matrix.Lists))
		})
	}
}
