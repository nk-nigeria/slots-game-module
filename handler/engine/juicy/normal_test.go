package juicy

import (
	"testing"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	api "github.com/ciaolink-game-platform/cgp-common/proto"
	"github.com/stretchr/testify/assert"
)

func Test_normal_SpinMatrix(t *testing.T) {
	name := "Test_normal_SpinMatrix"
	e := NewNormal(nil)
	engine := e.(*normal)
	numWildRatio1_0 := 0
	numWildRatio1_2 := 0
	numWildRatio1_5 := 0
	numWildRatio2_0 := 0
	m := entity.NewSlotMatrix(5, 3)
	countWildSymbol := func(matrix entity.SlotMatrix) int {
		num := 0
		matrix.ForEeach(func(idx, row, col int, symbol api.SiXiangSymbol) {
			if symbol == api.SiXiangSymbol_SI_XIANG_SYMBOL_WILD {
				num++
			}
		})
		return num
	}
	t.Run(name, func(t *testing.T) {
		for i := 0; i < 10000; i++ {
			matrix := engine.SpinMatrix(m, ratioWild1_0)
			numWildRatio1_0 += countWildSymbol(matrix)

			matrix1_2 := engine.SpinMatrix(m, ratioWild1_2)
			numWildRatio1_2 += countWildSymbol(matrix1_2)

			matrix1_5 := engine.SpinMatrix(m, ratioWild1_5)
			numWildRatio1_5 += countWildSymbol(matrix1_5)

			matrix2_0 := engine.SpinMatrix(m, ratioWild2_0)
			numWildRatio2_0 += countWildSymbol(matrix2_0)
		}

		assert.Less(t, ratioWild1_0, ratioWild1_2)
		assert.Less(t, ratioWild1_2, ratioWild1_5)
		assert.Less(t, ratioWild1_5, ratioWild2_0)
		t.Logf("ratio 1.0 %d 1.2 %d 1.5 %d 2.0 %d \r\n", numWildRatio1_0, numWildRatio1_2, numWildRatio1_5, numWildRatio2_0)
	})
}
