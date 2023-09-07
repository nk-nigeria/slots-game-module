package entity

import pb "github.com/ciaolink-game-platform/cgp-common/proto"

func PbSlotMatrixForEeach(sm *pb.SlotMatrix, fn func(idx int, row, col int32, symbol pb.SiXiangSymbol)) {
	row := int32(0)
	col := int32(0)
	cols := sm.Cols
	for idx, symbol := range sm.Lists {
		if idx != 0 {
			if col%cols == 0 {
				row++
				col = 0
			}
		}
		fn(idx, row, col, symbol)
		col++
	}
}
