package handler

import (
	"fmt"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
	"github.com/ciaolink-game-platform/cgp-common/utilities"
)

var _ lib.Engine = &slotsEngine{}

var RellsAllowScatter = map[int]bool{0: true, 2: true, 4: true}

func AllowScatter(col int) bool {
	return RellsAllowScatter[col]
}

type slotsEngine struct {
}

func NewSlotsEngine() lib.Engine {
	engine := slotsEngine{}
	return &engine
}

func (e *slotsEngine) NewGame(matchState interface{}) (interface{}, error) {
	return nil, nil
}

func (e *slotsEngine) Random(min, max int) int {
	return utilities.RandomNumber(min, max)
}

func (e *slotsEngine) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	matrix := e.SpinMatrix(s.GetMatrix())
	s.SetMatrix(matrix)
	spreadMatrix := e.SpreadWildInMatrix(matrix)
	s.SetSpreadMMatrix(spreadMatrix)
	paylines := e.PaylineMatrix(spreadMatrix)
	paylines = e.FilterPayline(paylines, func(numOccur int) bool {
		return numOccur >= 3
	})
	chipsMcb := s.GetBetInfo().Chips
	for _, payline := range paylines {
		payline.Rate = e.RatioPayline(payline)
		payline.Chips = int64(payline.Rate) * chipsMcb
	}
	s.SetPaylines(paylines)
	return s, nil
}

func (e *slotsEngine) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	slotDesk := &pb.SlotDesk{}
	// set matrix spin
	{
		sm := s.GetMatrix()
		slotDesk.Matrix = sm.ToPbSlotMatrix()
	}
	slotDesk.ChipsMcb = s.GetBetInfo().GetChips()
	// set matrix spread matrix wild symbol
	{
		sm := s.GetSpreadMatrix()
		slotDesk.SpreadMatrix = sm.ToPbSlotMatrix()
	}
	// add payline result
	{
		slotDesk.Paylines = s.GetPaylines()
		for _, payline := range slotDesk.Paylines {
			slotDesk.TotalChipsWin += payline.GetChips()
		}
	}
	// check if win bonus game
	{
		nextSiXiangGame := e.GetNextSiXiangGame(s)
		// s.AddTrackingPlayBonusGame(nextSiXiangGame)
		s.NextSiXiangGame = nextSiXiangGame
		slotDesk.SixiangGame = nextSiXiangGame
	}
	return slotDesk, nil
}

func (e *slotsEngine) SpinMatrix(matrix entity.SlotMatrix) entity.SlotMatrix {
	// matrix := matchState.GetMatrix()
	mapColExistScatter := make(map[int]bool)
	matrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		for {
			numRandom := e.Random(0, len(entity.ListSymbol))
			symbol := entity.ListSymbol[numRandom]
			if symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER {
				// Scatter only allow appear in list RellsAllowScatter
				if !AllowScatter(col) {
					continue
				}
				// check if symbol scatter already exist in this row
				if mapColExistScatter[col] {
					continue
				}
				mapColExistScatter[col] = true
			}
			matrix.List[idx] = symbol

			break
		}
		return
	})
	return matrix
}

func (e *slotsEngine) SpreadWildInMatrix(matrix entity.SlotMatrix) entity.SlotMatrix {
	// matrix := matchState.GetMatrix()
	spreadMatrix := entity.SlotMatrix{
		List: make([]pb.SiXiangSymbol, len(matrix.List), len(matrix.List)),
		Cols: matrix.Cols,
		Rows: matrix.Rows,
	}

	matrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		if symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_WILD {
			cols := spreadMatrix.Cols
			spreadMatrix.List[idx] = symbol
			for row := 0; row < matrix.Rows; row++ {
				id := row*cols + col
				spreadMatrix.List[id] = symbol
			}

			return
		}
		if spreadMatrix.List[idx] != pb.SiXiangSymbol_SI_XIANG_SYMBOL_WILD {
			spreadMatrix.List[idx] = symbol
		}
		return
	})
	return spreadMatrix
}

func (e *slotsEngine) PaylineMatrix(matrix entity.SlotMatrix) []*pb.Payline {
	paylines := make([]*pb.Payline, 0)
	payline := &pb.Payline{}
	// matrix.ForEeachLine(func(line int, symbols []pb.SiXiangSymbol) {
	for _, indexs := range entity.MapPaylineIdx {
		symbols := matrix.ListFromIndexs(indexs)
		for id, val := range entity.ListSymbol {
			numOccur := 0
			if val == pb.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER {
				continue
			}
			if val == pb.SiXiangSymbol_SI_XIANG_SYMBOL_WILD {
				continue
			}
			for _, symbol := range symbols {
				if (symbol & val) > 0 {
					numOccur++
					continue
				}
				if numOccur >= 3 {
					break
				}
				numOccur = 0
			}
			if numOccur > int(payline.NumOccur) {
				payline.Id = int32(id)
				payline.NumOccur = int32(numOccur)
				payline.Symbol = val
			}
			if numOccur >= 3 {
				break
			}
		}
		paylines = append(paylines, payline)
		payline = &pb.Payline{}
	}
	return paylines
}

func (e *slotsEngine) FilterPayline(paylines []*pb.Payline, fn func(numOccur int) bool) []*pb.Payline {
	list := make([]*pb.Payline, 0)
	for _, payline := range paylines {
		if fn(int(payline.GetNumOccur())) {
			list = append(list, payline)
		}
	}
	return list
}
func (e *slotsEngine) GetNextSiXiangGame(matchState *entity.SlotsMatchState) pb.SiXiangGame {
	matrix := matchState.GetMatrix()
	numScatter := 0
	matrix.ForEeachLine(func(line int, symbols []pb.SiXiangSymbol) {
		if !RellsAllowScatter[line] {
			return
		}
		for _, symbol := range symbols {
			if symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER {
				numScatter++
				break
			}
		}
	})
	if numScatter == 3 {
		return pb.SiXiangGame_SI_XIANG_GAME_BONUS
	}
	return pb.SiXiangGame_SI_XIANG_GAME_NOMAL
}

func (e *slotsEngine) RatioPayline(payline *pb.Payline) float64 {
	switch payline.Symbol {
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_10,
		pb.SiXiangSymbol_SI_XIANG_SYMBOL_J,
		pb.SiXiangSymbol_SI_XIANG_SYMBOL_Q, pb.SiXiangSymbol_SI_XIANG_SYMBOL_K:
		if payline.NumOccur == 3 {
			return 0.5
		}
		if payline.NumOccur == 4 {
			return 2.5
		}
		if payline.NumOccur == 5 {
			return 5
		}
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_BLUE_DRAGON:
		if payline.NumOccur == 3 {
			return 2
		}
		if payline.NumOccur == 4 {
			return 10
		}
		if payline.NumOccur == 5 {
			return 20
		}
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_WHITE_TIGER:
		if payline.NumOccur == 3 {
			return 1.5
		}
		if payline.NumOccur == 4 {
			return 7.5
		}
		if payline.NumOccur == 5 {
			return 15
		}
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_WARRIOR:
		if payline.NumOccur == 3 {
			return 1.2
		}
		if payline.NumOccur == 4 {
			return 6
		}
		if payline.NumOccur == 5 {
			return 12
		}
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_VERMILION_BIRD:
		if payline.NumOccur == 3 {
			return 1
		}
		if payline.NumOccur == 4 {
			return 5
		}
		if payline.NumOccur == 5 {
			return 10
		}
	}
	return 0
}

func (e *slotsEngine) PrintMatrix(matrix entity.SlotMatrix) {
	// matrix := matchState.GetMatrix()
	matrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		if idx != 0 && col == 0 {
			fmt.Println("")
		}
		fmt.Printf("%8d", symbol.Number())
		return
	})
	fmt.Println("")
}
