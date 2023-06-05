package sixiang

import (
	"testing"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	api "github.com/ciaolink-game-platform/cgp-common/proto"
	"github.com/stretchr/testify/assert"
)

func Test_sixiangBonusEngine_NewGame(t *testing.T) {
	name := "Test_sixiangBonusEngine_NewGame"
	t.Run(name, func(t *testing.T) {
		e := &sixiangBonusEngine{}
		matchState := entity.NewSlotsMathState(nil)
		got, err := e.NewGame(matchState)
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, len(entity.ListSymbolSiXiangBonusGame), len(matchState.MatrixSpecial.List))
		assert.Nil(t, matchState.SpinSymbols)
		trackSymbol := make(map[api.SiXiangSymbol]int)
		matchState.MatrixSpecial.ForEeach(func(_, _, _ int, symbol api.SiXiangSymbol) {
			num := trackSymbol[symbol]
			num++
			trackSymbol[symbol] = num
		})
		for _, v := range trackSymbol {
			assert.Equal(t, 1, v)
		}
	})
}

func Test_sixiangBonusEngine_Process(t *testing.T) {
	name := "Test_sixiangBonusEngine_Process"
	t.Run(name, func(t *testing.T) {
		e := &sixiangBonusEngine{}
		matchState := entity.NewSlotsMathState(nil)
		got, err := e.NewGame(matchState)
		assert.NoError(t, err)
		assert.NotNil(t, got)
		got, err = e.Process(matchState)
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, 1, len(matchState.SpinSymbols))
		assert.Equal(t, 1, len(matchState.MatrixSpecial.TrackFlip))
		for k := range matchState.MatrixSpecial.TrackFlip {
			assert.Equal(t, matchState.SpinSymbols[0].Symbol, matchState.MatrixSpecial.List[k])
		}
	})
}

func Test_sixiangBonusEngine_Finish(t *testing.T) {
	type args struct {
	}
	type test struct {
		name string
		arg  *entity.SlotsMatchState
		want *api.SlotDesk
	}
	engine := &sixiangBonusEngine{}
	tests := make([]test, 0)
	arrSym := entity.ListSymbolSiXiangBonusGame
	arrNextGame := []api.SiXiangGame{
		api.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_DRAGON_PEARL,
		api.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_LUCKDRAW,
		api.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_GOLDPICK,
		api.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_RAPIDPAY,
	}
	for idx, sym := range arrSym {
		matchState := entity.NewSlotsMathState(nil)
		matchState.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS
		engine.NewGame(matchState)
		idRandom := 0
		matchState.MatrixSpecial.ForEeach(func(idx, _, _ int, symbol api.SiXiangSymbol) {
			if symbol == sym {
				idRandom = idx
			}
		})
		matchState.MatrixSpecial.Flip(idRandom)
		matchState.SpinSymbols = []*api.SpinSymbol{
			{Symbol: sym},
		}
		slotDesk := &api.SlotDesk{
			CurrentSixiangGame: api.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS,
			NextSixiangGame:    arrNextGame[idx],
			IsFinishGame:       true,
		}
		test := test{
			name: "Test_sixiangBonusEngine_Finish_" + sym.String(),
			arg:  matchState,
			want: slotDesk,
		}
		tests = append(tests, test)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.Finish(tt.arg)
			assert.NoError(t, err)
			assert.NotNil(t, result)
			slotDesk := result.(*api.SlotDesk)
			assert.Equal(t, tt.want.CurrentSixiangGame, slotDesk.CurrentSixiangGame)
			assert.Equal(t, tt.want.NextSixiangGame, slotDesk.NextSixiangGame)
			assert.Equal(t, tt.want.IsFinishGame, slotDesk.IsFinishGame)
		})
	}
}
