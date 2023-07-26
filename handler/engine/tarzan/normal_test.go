package tarzan

import (
	"testing"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	api "github.com/ciaolink-game-platform/cgp-common/proto"
	"github.com/stretchr/testify/assert"
)

func Test_normal_NewGame(t *testing.T) {
	name := "Test_normal_NewGame"
	t.Run(name, func(t *testing.T) {
		e := NewNormal(func(i1, i2 int) int {
			return i1
		})
		s := entity.NewSlotsMathState(nil)
		got, err := e.NewGame(s)
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, entity.ColsTarzanMatrix*entity.RowsTarzanMatrix, len(s.Matrix.List))
		matrix := s.Matrix
		matrix.ForEeach(func(idx, row, col int, symbol api.SiXiangSymbol) {
			assert.NotEqual(t, api.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED, symbol)
		})
	})
}

func Test_normal_SpinMatrix(t *testing.T) {
	name := "Test_normal_SpinMatrix"
	s := entity.NewSlotsMathState(nil)
	e := NewNormal(func(i1, i2 int) int { return i1 })
	e.NewGame(s)
	engine := e.(*normal)
	for i := 0; i < 5000; i++ {
		t.Run(name, func(t *testing.T) {
			matrix := engine.SpinMatrix(s.Matrix)
			assert.Equal(t, entity.ColsTarzanMatrix*entity.RowsTarzanMatrix, len(matrix.List))
			numTarzanSymbol := 0
			numLetterSymbol := 0
			numFreeSpin := 0
			matrix.ForEeach(func(idx, row, col int, symbol api.SiXiangSymbol) {
				if symbol == api.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN {
					numTarzanSymbol++
				}
				if entity.TarzanLetterSymbol[symbol] {
					numLetterSymbol++
				}
				if symbol == api.SiXiangSymbol_SI_XIANG_SYMBOL_FREE_SPIN {
					numFreeSpin++
				}
			})
			assert.Equal(t, true, numTarzanSymbol <= engine.maxDropTarzanSymbol)
			assert.Equal(t, true, numLetterSymbol <= engine.maxDropLetterSymbol)
			assert.Equal(t, true, numFreeSpin <= engine.maxDropFreeSpin)
		})
	}
}

func Test_normal_SpinMatrix_Freex9(t *testing.T) {
	name := "Test_normal_SpinMatrix_FreeX9"
	indexFreeSpin := 0
	for idx, sym := range entity.TarzanSymbols {
		if sym == api.SiXiangSymbol_SI_XIANG_SYMBOL_FREE_SPIN {
			indexFreeSpin = idx
			break
		}
	}
	s := entity.NewSlotsMathState(nil)
	e := NewNormal(func(i1, i2 int) int { return indexFreeSpin })
	e.NewGame(s)
	engine := e.(*normal)
	assert.Equal(t, true, engine.maxDropFreeSpin >= 2)
	// assert.Equal(t, true, engine.allowDropFreeSpinx9)
	for i := 0; i < 5000; i++ {
		t.Run(name, func(t *testing.T) {
			matrix := engine.SpinMatrix(s.Matrix)
			matrix = engine.SpinMatrix(s.Matrix)
			matrix = engine.SpinMatrix(s.Matrix)
			numLetterSymbol := 0
			matrix.ForEeach(func(idx, row, col int, symbol api.SiXiangSymbol) {
				if symbol == api.SiXiangSymbol_SI_XIANG_SYMBOL_FREE_SPIN {
					numLetterSymbol++
				}
			})
			assert.Equal(t, true, numLetterSymbol >= 3)
		})
	}
}

func Test_normal_TarzanSwing(t *testing.T) {
	name := "Test_normal_TarzanSwing"
	s := entity.NewSlotsMathState(nil)
	e := NewNormal(func(i1, i2 int) int { return i1 })
	e.NewGame(s)
	engine := e.(*normal)
	for i := 0; i < 100; i++ {

		t.Run(name, func(t *testing.T) {
			matrix := engine.SpinMatrix(s.Matrix)
			hasTarzanSymbol := false
			matrix.ForEeach(func(idx, row, col int, symbol api.SiXiangSymbol) {
				if symbol == api.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN {
					hasTarzanSymbol = true
				}
			})
			if !hasTarzanSymbol {
				matrix.List[len(matrix.List)-1] = api.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN
			}
			swingMatrix := engine.TarzanSwing(matrix)
			assert.Equal(t, matrix.Size, swingMatrix.Size)
			assert.NotEqual(t, 0, swingMatrix.Size)
			swingMatrix.ForEeach(func(idx, row, col int, symbol api.SiXiangSymbol) {
				assert.NotEqual(t, api.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED, symbol)
			})
			matrix.ForEeach(func(idx, row, col int, symbol api.SiXiangSymbol) {
				if entity.TarzanMidSymbol[symbol] {
					assert.Equal(t, api.SiXiangSymbol_SI_XIANG_SYMBOL_WILD, swingMatrix.List[idx])
				} else {
					assert.NotEqual(t, api.SiXiangSymbol_SI_XIANG_SYMBOL_WILD, swingMatrix.List[idx])
					assert.NotEqual(t, api.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED, swingMatrix.List[idx])
				}
			})
		})
	}
}

func Test_normal_Paylines(t *testing.T) {
	name := "Test_normal_TarzanSwing"
	t.Run(name, func(t *testing.T) {
		matrix := entity.NewTarzanMatrix()
		indexs, _ := entity.PaylineTarzanMapping.Get(12)
		for _, val := range indexs {
			matrix.List[val] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JANE
		}
		e := NewNormal(func(i1, i2 int) int { return i1 })
		engine := e.(*normal)
		paylines := engine.Paylines(matrix)
		assert.NotNil(t, paylines)
		assert.Equal(t, 1, len(paylines))
		assert.Equal(t, int32(12), paylines[0].Id)
		assert.Equal(t, api.SiXiangSymbol_SI_XIANG_SYMBOL_JANE, paylines[0].Symbol)
	})
}

func Test_normal_Finish(t *testing.T) {
	name := "Tarzan Normal Finish"
	matchState := entity.NewSlotsMathState(nil)
	engine := NewNormal(func(i1, i2 int) int {
		return i1
	})
	matchState.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_NORMAL
	matchState.SetBetInfo(&api.InfoBet{
		Chips: 100,
	})
	engine.NewGame(matchState)
	engine.Process(matchState)
	t.Run(name, func(t *testing.T) {
		result, err := engine.Finish(matchState)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		slotdesk := result.(*api.SlotDesk)
		assert.Equal(t, matchState.Bet().Chips, slotdesk.ChipsMcb)
		assert.Equal(t, matchState.CurrentSiXiangGame, slotdesk.CurrentSixiangGame)
		assert.Equal(t, matchState.NextSiXiangGame, slotdesk.NextSixiangGame)
		assert.NotZero(t, len(slotdesk.Matrix.Lists))
		assert.NotZero(t, len(slotdesk.Paylines))
		assert.NotZero(t, len(slotdesk.SpreadMatrix.Lists))
	})
}

func Test_normal_Process_Stress(t *testing.T) {
	name := "Test_normal_Process_Stress"
	t.Run(name, func(t *testing.T) {
		e := NewNormal(nil)
		matchState := entity.NewSlotsMathState(nil)
		for i := 0; i < 10000; i++ {
			got, err := e.Process(matchState)
			assert.NoError(t, err)
			assert.NotNil(t, got)
		}
	})
}
