package sixiang

import (
	"testing"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	api "github.com/ciaolink-game-platform/cgp-common/proto"
	"github.com/stretchr/testify/assert"
)

func Test_bonusEngine_NewGame(t *testing.T) {
	name := "bonusEngine_NewGame"
	t.Run(name, func(t *testing.T) {
		e := NewBonusEngine(nil)
		matchState := entity.NewSlotsMathState(nil)
		got, err := e.NewGame(matchState)
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, len(entity.ListSymbolBonusGame), len(matchState.MatrixSpecial.List))
		t.Log(matchState.MatrixSpecial)
	})
}

func Test_bonusEngine_Process(t *testing.T) {
	name := "bonusEngine_Process"
	t.Run(name, func(t *testing.T) {
		e := NewBonusEngine(nil)
		matchState := entity.NewSlotsMathState(nil)
		_, _ = e.NewGame(matchState)
		matchState.MatrixSpecial.List[4] = api.SiXiangSymbol_SI_XIANG_SYMBOL_BONUS_DRAGONBALL
		matchState.SetBetInfo(&api.InfoBet{
			Id: 4,
		})
		e.Process(matchState)
		matchState.SpinSymbols[0].Symbol = api.SiXiangSymbol_SI_XIANG_SYMBOL_BONUS_DRAGONBALL
	})
}

func Test_bonusEngine_FullFlow(t *testing.T) {
	type args struct {
		idRandom       int
		symbolWantDraw api.SiXiangSymbol
		chip           int64
	}
	type wants struct {
		nextGame api.SiXiangGame
		chips    int64
	}
	tests := []struct {
		name string
		args args
		want wants
	}{
		{
			name: "bonusEngine_FullFlow_draw_dragon_ball",
			args: args{
				idRandom:       2,
				symbolWantDraw: api.SiXiangSymbol_SI_XIANG_SYMBOL_BONUS_DRAGONBALL,
				chip:           1,
			},
			want: wants{
				nextGame: api.SiXiangGame_SI_XIANG_GAME_DRAGON_PEARL,
				chips:    0,
			},
		},
		{
			name: "bonusEngine_FullFlow_draw_dragon_luckdraw",
			args: args{
				idRandom:       2,
				symbolWantDraw: api.SiXiangSymbol_SI_XIANG_SYMBOL_BONUS_LUCKYDRAW,
				chip:           1,
			},
			want: wants{
				nextGame: api.SiXiangGame_SI_XIANG_GAME_LUCKDRAW,
				chips:    0,
			},
		},
		{
			name: "bonusEngine_FullFlow_draw_dragon_gold_pick",
			args: args{
				idRandom:       2,
				symbolWantDraw: api.SiXiangSymbol_SI_XIANG_SYMBOL_BONUS_GOLDPICK,
				chip:           1,
			},
			want: wants{
				nextGame: api.SiXiangGame_SI_XIANG_GAME_GOLDPICK,
				chips:    0,
			},
		},
		{
			name: "bonusEngine_FullFlow_draw_dragon_rapid_pay",
			args: args{
				idRandom:       2,
				symbolWantDraw: api.SiXiangSymbol_SI_XIANG_SYMBOL_BONUS_RAPIDPAY,
				chip:           1,
			},
			want: wants{
				nextGame: api.SiXiangGame_SI_XIANG_GAME_RAPIDPAY,
				chips:    0,
			},
		},
		{
			name: "bonusEngine_FullFlow_draw_dragon_gold_10",
			args: args{
				idRandom:       2,
				symbolWantDraw: api.SiXiangSymbol_SI_XIANG_SYMBOL_BONUS_GOLDX10,
				chip:           1,
			},
			want: wants{
				nextGame: api.SiXiangGame_SI_XIANG_GAME_NORMAL,
				chips: 1 *
					int64(entity.ListSymbolBonusGame[api.SiXiangSymbol_SI_XIANG_SYMBOL_BONUS_GOLDX10].Value.Min),
			},
		},
		{
			name: "bonusEngine_FullFlow_draw_dragon_gold_20",
			args: args{
				idRandom:       2,
				symbolWantDraw: api.SiXiangSymbol_SI_XIANG_SYMBOL_BONUS_GOLDX20,
				chip:           1,
			},
			want: wants{
				nextGame: api.SiXiangGame_SI_XIANG_GAME_NORMAL,
				chips: 1 *
					int64(entity.ListSymbolBonusGame[api.SiXiangSymbol_SI_XIANG_SYMBOL_BONUS_GOLDX20].Value.Min),
			},
		},
		{
			name: "bonusEngine_FullFlow_draw_dragon_gold_30",
			args: args{
				idRandom:       2,
				symbolWantDraw: api.SiXiangSymbol_SI_XIANG_SYMBOL_BONUS_GOLDX30,
				chip:           1,
			},
			want: wants{
				nextGame: api.SiXiangGame_SI_XIANG_GAME_NORMAL,
				chips: 1 *
					int64(entity.ListSymbolBonusGame[api.SiXiangSymbol_SI_XIANG_SYMBOL_BONUS_GOLDX30].Value.Min),
			},
		}, {
			name: "bonusEngine_FullFlow_draw_dragon_gold_50",
			args: args{
				idRandom:       2,
				symbolWantDraw: api.SiXiangSymbol_SI_XIANG_SYMBOL_BONUS_GOLDX50,
				chip:           1,
			},
			want: wants{
				nextGame: api.SiXiangGame_SI_XIANG_GAME_NORMAL,
				chips: 1 *
					int64(entity.ListSymbolBonusGame[api.SiXiangSymbol_SI_XIANG_SYMBOL_BONUS_GOLDX50].Value.Min),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewBonusEngine(func(min, max int) int {
				return tt.args.idRandom
			})
			matchState := entity.NewSlotsMathState(nil)
			matchState.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_BONUS
			_, _ = e.NewGame(matchState)
			matchState.SetBetInfo(&api.InfoBet{
				Chips: int64(tt.args.chip),
			})
			matchState.MatrixSpecial.List[tt.args.idRandom] = tt.args.symbolWantDraw
			e.Process(matchState)
			result, err := e.Finish(matchState)
			assert.NoError(t, err)
			assert.NotNil(t, result)
			slotDesk := result.(*api.SlotDesk)
			assert.Equal(t, tt.want.nextGame, slotDesk.NextSixiangGame)
			assert.Equal(t, tt.want.chips, slotDesk.GameReward.ChipsWin)
			assert.Equal(t, tt.want.chips, slotDesk.GameReward.ChipsWin)
			assert.Equal(t, tt.args.chip, slotDesk.ChipsMcb)
		})
	}
}
