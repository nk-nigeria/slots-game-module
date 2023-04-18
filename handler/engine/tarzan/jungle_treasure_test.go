package tarzan

import (
	"fmt"
	"testing"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	api "github.com/ciaolink-game-platform/cgp-common/proto"
	"github.com/stretchr/testify/assert"
)

func Test_jungleTreasure_Process_Num_Spin(t *testing.T) {
	name := "Test_jungleTreasure_Process_num_spin"
	engine := NewJungTreasure(nil)
	t.Run(name, func(t *testing.T) {
		for i := 0; i < 100; i++ {
			numSpin := 0
			matchState := entity.NewSlotsMathState(nil)
			matchState.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_TARZAN_JUNGLE_TREASURE
			engine.NewGame(matchState)

			for {
				_, err := engine.Process(matchState)
				if err == entity.ErrorSpinReachMax {
					break
				}
				engine.Finish(matchState)
				numSpin++
				assert.NoError(t, err)
			}
			assert.Equal(t, false, numSpin < 8, fmt.Sprintf("num spin 8 <= num_spin <= 15, current num_spin = %d", numSpin))
			assert.Equal(t, false, numSpin > 15, fmt.Sprintf("num spin 7 <= num_spin <= 15, current num_spin = %d", numSpin))
		}
	})
}

func Test_jungleTreasure_Process(t *testing.T) {
	name := "Test_jungleTreasure_Process"
	engine := NewJungTreasure(nil)
	t.Run(name, func(t *testing.T) {
		for i := 0; i < 100; i++ {
			matchState := entity.NewSlotsMathState(nil)
			matchState.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_TARZAN_JUNGLE_TREASURE
			engine.NewGame(matchState)
			result, err := engine.Process(matchState)
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, 4, matchState.NumSpinLeft)
			assert.NotZero(t, len(matchState.SpinSymbols))
			assert.NotZero(t, len(matchState.MatrixSpecial.TrackFlip))
		}
		{
			matchState := entity.NewSlotsMathState(nil)
			matchState.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_TARZAN_JUNGLE_TREASURE
			engine.NewGame(matchState)
			for {
				prevGemSpin := matchState.NumSpinLeft
				result, err := engine.Process(matchState)
				if err == entity.ErrorSpinReachMax {
					break
				}

				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, prevGemSpin-1, matchState.NumSpinLeft)
				assert.NotZero(t, len(matchState.SpinSymbols))
				assert.NotZero(t, len(matchState.MatrixSpecial.TrackFlip))
			}
		}
	})
}

func Test_jungleTreasure_Finish(t *testing.T) {
	name := "Test_jungleTreasure_Finish"
	engine := NewJungTreasure(nil)
	matchStates := make([]entity.SlotsMatchState, 0)
	for sym := range entity.TarzanJungleTreasureSymbol {
		spin := &api.SpinSymbol{
			Symbol: sym,
		}
		matchState := entity.NewSlotsMathState(nil)
		matchState.SetBetInfo(&api.InfoBet{
			Chips: 100,
		})
		matchState.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_TARZAN_JUNGLE_TREASURE
		engine.NewGame(matchState)
		matchState.SpinSymbols = append(matchState.SpinSymbols, spin)
		matchStates = append(matchStates, *matchState)
	}
	t.Run(name, func(t *testing.T) {
		for _, s := range matchStates {
			prevGemspin := s.NumSpinLeft
			result, err := engine.Finish(&s)
			assert.NoError(t, err)
			assert.NotNil(t, result)
			slotDesk := result.(*api.SlotDesk)
			assert.NotZero(t, len(s.SpinSymbols))
			switch s.SpinSymbols[0].Symbol {
			case api.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN_MORE_TURNX2, api.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN_MORE_TURNX3:
				v := entity.TarzanJungleTreasureSymbol[s.SpinSymbols[0].Symbol]
				prevGemspin += int(v.Value.Min)
			default:
				symInfo := entity.TarzanJungleTreasureSymbol[s.SpinSymbols[0].Symbol]
				assert.Equal(t, true, s.LineWinByGame[s.CurrentSiXiangGame] >= int(symInfo.Value.Min), fmt.Sprintf("sym %s val %d, expext > %v", s.SpinSymbols[0].Symbol.String(), s.LineWinByGame[s.CurrentSiXiangGame], symInfo.Value.Min))
				assert.Equal(t, true, s.LineWinByGame[s.CurrentSiXiangGame] <= int(symInfo.Value.Max), fmt.Sprintf("sym  %s val %d, expext < %v", s.SpinSymbols[0].Symbol.String(), s.LineWinByGame[s.CurrentSiXiangGame], symInfo.Value.Max))
			}
			assert.Equal(t, prevGemspin, s.NumSpinLeft)
			assert.Equal(t, slotDesk.CurrentSixiangGame, s.CurrentSiXiangGame)
			assert.Equal(t, slotDesk.NextSixiangGame, s.NextSiXiangGame)
		}
	})
}
