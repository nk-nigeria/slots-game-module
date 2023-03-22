package entity

import (
	"testing"

	api "github.com/ciaolink-game-platform/cgp-common/proto"
	"github.com/stretchr/testify/assert"
)

func TestSlotMatrix_IsPayline(t *testing.T) {
	name := "TestSlotMatrix_IsPayline"
	t.Run(name, func(t *testing.T) {
		matrix := NewTarzanMatrix()
		indexs := []int{5, 6, 7, 8, 9}
		for _, val := range indexs {
			matrix.List[val] = api.SiXiangSymbol_SI_XIANG_SYMBOL_GORILLE + api.SiXiangSymbol(val)
		}
		list, isPayline := matrix.IsPayline(indexs)
		assert.Equal(t, false, isPayline)
		assert.Equal(t, 0, len(list))

		for _, val := range indexs {
			matrix.List[val] = api.SiXiangSymbol_SI_XIANG_SYMBOL_LETTER_E
		}
		arr := []api.SiXiangSymbol{
			api.SiXiangSymbol_SI_XIANG_SYMBOL_LETTER_E,
			api.SiXiangSymbol_SI_XIANG_SYMBOL_WILD,
			api.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN,
		}

		lastIdx := indexs[len(indexs)-1]
		for _, sym := range arr {
			matrix.List[lastIdx] = sym
			list, isPayline = matrix.IsPayline(indexs)
			assert.Equal(t, true, isPayline)
			assert.Equal(t, indexs, list)
		}
		matrix.List[lastIdx] = api.SiXiangSymbol_SI_XIANG_SYMBOL_LETTER_J
		_, isPayline = matrix.IsPayline(indexs)
		assert.Equal(t, false, isPayline)
	})
}
