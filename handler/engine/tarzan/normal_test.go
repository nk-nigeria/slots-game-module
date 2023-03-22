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
		s := entity.NewTarzanMatchState(nil)
		got, err := e.NewGame(s)
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, entity.ColsTarzanMatrix*entity.RowsTarzanMatrix, len(s.Matrix.List))
		s.Matrix.ForEeach(func(idx, row, col int, symbol api.SiXiangSymbol) {
			assert.NotEqual(t, api.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED, symbol)
		})
	})
}

func Test_normal_SpinMatrix(t *testing.T) {
	name := "Test_normal_SpinMatrix"
	s := entity.NewTarzanMatchState(nil)
	e := NewNormal(func(i1, i2 int) int { return i1 })
	e.NewGame(s)
	engine := e.(*normal)
	for i := 0; i < 5000; i++ {
		t.Run(name, func(t *testing.T) {
			matrix := engine.SpinMatrix(s.Matrix)
			assert.Equal(t, entity.ColsTarzanMatrix*entity.RowsTarzanMatrix, len(matrix.List))
			numTarzanSymbol := 0
			numLetterSymbol := 0
			matrix.ForEeach(func(idx, row, col int, symbol api.SiXiangSymbol) {
				if symbol == api.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN {
					numTarzanSymbol++
				}
				if entity.TarzanLetterSymbol[symbol] {
					numLetterSymbol++
				}
			})
			assert.Equal(t, true, numTarzanSymbol <= 1)
			assert.Equal(t, true, numLetterSymbol <= 1)
		})
	}
}

func Test_normal_TarzanSwing(t *testing.T) {
	name := "Test_normal_TarzanSwing"
	s := entity.NewTarzanMatchState(nil)
	e := NewNormal(func(i1, i2 int) int { return i1 })
	e.NewGame(s)
	engine := e.(*normal)
	for i := 0; i < 5000; i++ {

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
			matrix.ForEeach(func(idx, row, col int, symbol api.SiXiangSymbol) {
				if entity.TarzanMidSymbol[symbol] {
					assert.Equal(t, api.SiXiangSymbol_SI_XIANG_SYMBOL_WILD, swingMatrix.List[idx])
				} else {
					assert.NotEqual(t, api.SiXiangSymbol_SI_XIANG_SYMBOL_WILD, swingMatrix.List[idx])
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
