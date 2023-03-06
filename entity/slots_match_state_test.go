package entity

import (
	"fmt"
	"testing"

	pb "github.com/ciaolink-game-platform/cgp-common/proto"
	"github.com/stretchr/testify/assert"
)

func TestSlotMatrix_ForEeachLine(t *testing.T) {
	type fields struct {
		Sm SlotMatrix
	}

	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
		{
			name: "Test Matrix_ForEeachLine",
			fields: fields{
				Sm: NewSiXiangMatrixNormal(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := tt.fields.Sm
			lastLine := 0
			sm.ForEeachLine(func(line int, symbols []pb.SiXiangSymbol) {
				assert.Equal(t, sm.Cols, len(symbols), fmt.Sprintf("in line %d", line))
				lastLine = line
			})
			assert.Equal(t, sm.Rows-1, lastLine, fmt.Sprintf("last line "))
		})
	}
}

func TestNewSiXiangMatrixLuckyDraw(t *testing.T) {
	name := "TestNewSiXiangMatrixLuckyDraw"
	t.Run(name, func(t *testing.T) {
		mapCountSymbol := make(map[pb.SiXiangSymbol]int)
		matrix := NewSiXiangMatrixLuckyDraw()
		assert.Equal(t, len(matrix.List), RowsMatrix*ColsMatrix)
		for _, symbol := range matrix.List {
			num := mapCountSymbol[symbol]
			num++
			mapCountSymbol[symbol] = num
		}
		for k, v := range mapCountSymbol {
			if k < pb.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_GOLD_1 {
				assert.Equal(t, 3, v)
			} else {
				assert.Equal(t, 1, v)
			}
		}
	})
}

func TestSlotMatrix_RowCol(t *testing.T) {
	type fields struct {
		List []pb.SiXiangSymbol
		Cols int
		Rows int
	}
	type args struct {
		id int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  int
	}{
		// TODO: Add test cases.
		{
			name: "Matrix_RowCol",
			fields: fields{
				Cols: 5,
				Rows: 3,
			},
			args: args{
				id: 10,
			},
			want:  2,
			want1: 0,
		},
		{
			name: "Matrix_RowCol",
			fields: fields{
				Cols: 5,
				Rows: 3,
			},
			args: args{
				id: 7,
			},
			want:  1,
			want1: 2,
		},
		{
			name: "Matrix_RowCol",
			fields: fields{
				Cols: 5,
				Rows: 3,
			},
			args: args{
				id: 4,
			},
			want:  0,
			want1: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := &SlotMatrix{
				List: tt.fields.List,
				Cols: tt.fields.Cols,
				Rows: tt.fields.Rows,
			}
			got, got1 := sm.RowCol(tt.args.id)
			if got != tt.want {
				t.Errorf("SlotMatrix.RowCol() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SlotMatrix.RowCol() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSlotMatrix_ToPbSlotMatrix(t *testing.T) {
	name := "TestSlotMatrix_ToPbSlotMatrix"
	t.Run(name, func(t *testing.T) {
		sm := NewSiXiangMatrixDragonPearl()
		result := sm.ToPbSlotMatrix()
		result.Lists[0] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED
		assert.NotEqual(t, sm.List[0], result.Lists[0])
		assert.Equal(t, len(result.Lists), len(sm.List))
		assert.Equal(t, result.Rows, sm.Rows)
		assert.Equal(t, result.Cols, sm.Cols)
	})

}
