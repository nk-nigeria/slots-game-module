package entity

import (
	"fmt"
	"testing"

	pb "github.com/nk-nigeria/cgp-common/proto"
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
				Sm: NewSlotMatrix(),
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
