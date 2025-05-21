package entity

import (
	"testing"

	api "github.com/nakamaFramework/cgp-common/proto"
	pb "github.com/nakamaFramework/cgp-common/proto"
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
		list, isPayline := matrix.IsTarzanPayline(matrix, indexs)
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
			list, isPayline = matrix.IsTarzanPayline(matrix, indexs)
			assert.Equal(t, true, isPayline)
			assert.Equal(t, indexs, list)
		}
		// test not payline
		matrix.List[lastIdx] = api.SiXiangSymbol_SI_XIANG_SYMBOL_LETTER_J
		_, isPayline = matrix.IsTarzanPayline(matrix, indexs)
		assert.Equal(t, false, isPayline)
		// test wild at begin
		matrix.List[5] = api.SiXiangSymbol_SI_XIANG_SYMBOL_WILD
		for _, sym := range arr {
			matrix.List[lastIdx] = sym
			list, isPayline = matrix.IsTarzanPayline(matrix, indexs)
			assert.Equal(t, true, isPayline)
			assert.Equal(t, indexs, list)
		}
	})
}

func TestNewSlotMatrix(t *testing.T) {
	name := "TestNewSlotMatrix"
	t.Run(name, func(t *testing.T) {
		rows := 5
		cols := 3
		matrix := NewSlotMatrix(rows, cols)

		assert.Equal(t, rows*cols, matrix.Size)
		assert.NotNil(t, matrix.List)
		assert.NotNil(t, matrix.TrackFlip)
	})
}

func TestSlotMatrix_ForEeach(t *testing.T) {
	slostMatrix := NewSiXiangMatrixNormal()
	slostMatrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		t.Logf("id: %d row %d col %d", idx, row, col)
	})
	assert.Equal(t, true, false)
}

func TestSlotMatrix_ForEeachNotFlip(t *testing.T) {
	type fields struct {
		List      []pb.SiXiangSymbol
		Cols      int
		Rows      int
		Size      int
		TrackFlip map[int]bool
	}
	type args struct {
		fn func(idx, row, col int, symbol pb.SiXiangSymbol)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := &SlotMatrix{
				List:      tt.fields.List,
				Cols:      tt.fields.Cols,
				Rows:      tt.fields.Rows,
				Size:      tt.fields.Size,
				TrackFlip: tt.fields.TrackFlip,
			}
			sm.ForEeachNotFlip(tt.args.fn)
		})
	}
}
