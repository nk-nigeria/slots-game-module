package tarzan

import (
	"math"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

var _ lib.Engine = &normal{}

type normal struct {
	maxDropTarzanSymbol int
	maxDropFreeSpin     int
	allowDropFreeSpinx9 bool
	maxDropLetterSymbol int
	randomIntFn         func(int, int) int
}

// todo save JUNGLE when exit
func NewNormal(randomIntFn func(int, int) int) lib.Engine {
	e := &normal{
		maxDropTarzanSymbol: 1,
		maxDropLetterSymbol: 1,
		maxDropFreeSpin:     math.MaxInt,
		allowDropFreeSpinx9: true,
	}
	if randomIntFn == nil {
		e.randomIntFn = RandomInt
	} else {
		e.randomIntFn = randomIntFn
	}
	return e
}

// NewGame implements lib.Engine
func (e *normal) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	matrix := entity.NewTarzanMatrix()
	s.SetMatrix(e.SpinMatrix(matrix))
	s.TrackIndexFreeSpinSymbol = make(map[int]bool)
	return s, nil
}

// Process implements lib.Engine
func (e *normal) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	matrix := s.Matrix
	matrix = e.SpinMatrix(matrix)
	s.SetMatrix(matrix)
	s.SetWildMatrix(e.TarzanSwing(matrix))
	// s.TrackIndexFreeSpinSymbol = make(map[int]bool)

	matrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		if entity.TarzanLetterSymbol[symbol] {
			s.AddCollectionSymbol(0, symbol)
		}
	})
	return matchState, nil
}

// Finish implements lib.Engine
func (e *normal) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	s.ChipWinByGame[s.CurrentSiXiangGame] = 0
	s.LineWinByGame[s.CurrentSiXiangGame] = 0
	slotDesk := &pb.SlotDesk{}
	slotDesk.Paylines = e.Paylines(s.WildMatrix)
	slotDesk.ChipsMcb = s.Bet().Chips
	lineWin := len(slotDesk.Paylines)
	matrix := s.Matrix
	matrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		if symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN {
			lineWin += 500
		}
	})
	s.ChipWinByGame[s.CurrentSiXiangGame] = int64(lineWin/100) * slotDesk.ChipsMcb
	s.LineWinByGame[s.CurrentSiXiangGame] = lineWin
	s.NextSiXiangGame = e.GetNextSiXiangGame(s)
	// if next game is freex9, save index freespin symbol
	if s.NextSiXiangGame == pb.SiXiangGame_SI_XIANG_GAME_TARZAN_FREESPINX9 {
		matrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
			if symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_FREE_SPIN {
				s.TrackIndexFreeSpinSymbol[idx] = true
			}
		})
	}
	slotDesk.Matrix = s.Matrix.ToPbSlotMatrix()
	slotDesk.SpreadMatrix = s.WildMatrix.ToPbSlotMatrix()
	slotDesk.ChipsWin = s.ChipWinByGame[s.CurrentSiXiangGame]
	slotDesk.TotalChipsWinByGame = slotDesk.ChipsWin
	slotDesk.CurrentSixiangGame = s.CurrentSiXiangGame
	slotDesk.NextSixiangGame = s.NextSiXiangGame
	slotDesk.IsFinishGame = true
	return slotDesk, nil
}

// Random implements lib.Engine
func (e *normal) Random(min int, max int) int {
	return e.randomIntFn(min, max)
}

func (e *normal) SpinMatrix(matrix entity.SlotMatrix) entity.SlotMatrix {
	numTarzanSymbolSpin := 0
	numLetterSymbolSpin := 0
	numFreeSpinSymbolSpin := 0
	matrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
	loop:
		for {
			numRandom := e.Random(0, len(entity.TarzanSymbols)-1)
			randSymbol := entity.TarzanSymbols[numRandom]
			switch randSymbol {
			// Tarzan symbol chỉ xuất hiện ở col 5 và chỉ xuất hiện 1 lần
			case pb.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN:
				if col != entity.Col_5 || numTarzanSymbolSpin >= e.maxDropTarzanSymbol {
					continue loop
				}
				numTarzanSymbolSpin++
			// chỉ xuất hiện free spin ở col 3, 4, 5
			case pb.SiXiangSymbol_SI_XIANG_SYMBOL_FREE_SPIN:
				if row < entity.Col_3 || numFreeSpinSymbolSpin >= e.maxDropLetterSymbol {
					continue loop
				}
				// kiểm tra điều kiện cho phép ra freespin symbol
				// nhưng không cho phép ra freespinx9 game
				if !e.allowDropFreeSpinx9 && e.countFreeSpinSymbolByCol(matrix) >= 2 {
					continue loop
				}
				numFreeSpinSymbolSpin++
			}
			// Letter symbol only one per spin
			if entity.TarzanLetterSymbol[randSymbol] {
				if numLetterSymbolSpin >= e.maxDropLetterSymbol {
					continue loop
				}
				numLetterSymbolSpin++
			}
			matrix.List[idx] = randSymbol
			break
		}
	})
	return matrix
}

func (e *normal) TarzanSwing(matrix entity.SlotMatrix) entity.SlotMatrix {
	swingMatrix := entity.SlotMatrix{
		List: make([]pb.SiXiangSymbol, matrix.Size),
		Cols: matrix.Cols,
		Rows: matrix.Rows,
		Size: matrix.Size,
	}
	copy(swingMatrix.List, matrix.List)
	hasTarzanSymbol := false
	matrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		if symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN {
			hasTarzanSymbol = true
		}
	})
	if !hasTarzanSymbol {
		return swingMatrix
	}
	swingMatrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		isMidSymbol := entity.TarzanMidSymbol[symbol]
		if isMidSymbol {
			swingMatrix.List[idx] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_WILD
		}
	})
	return swingMatrix
}

func (e *normal) GetNextSiXiangGame(s *entity.SlotsMatchState) pb.SiXiangGame {
	if s.SizeCollectionSymbol(0) == len(entity.TarzanLetterSymbol) {
		return pb.SiXiangGame_SI_XIANG_GAME_TARZAN_JUNGLE_TREASURE
	}
	matrix := s.Matrix
	nummFreeSpinSymbolPerCol := 0
	matrix.ForEachCol(func(col int, symbols []pb.SiXiangSymbol) {
		if col < entity.Col_3 {
			return
		}
		for _, sym := range symbols {
			if sym == pb.SiXiangSymbol_SI_XIANG_SYMBOL_FREE_SPIN {
				nummFreeSpinSymbolPerCol++
				break
			}
		}
	})
	if nummFreeSpinSymbolPerCol >= 3 {
		return pb.SiXiangGame_SI_XIANG_GAME_TARZAN_FREESPINX9
	}
	return pb.SiXiangGame_SI_XIANG_GAME_NORMAL
}

func (e *normal) Paylines(matrix entity.SlotMatrix) []*pb.Payline {
	paylines := make([]*pb.Payline, 0)
	for pair := entity.PaylineTarzanMapping.Oldest(); pair != nil; pair = pair.Next() {
		paylineIndexs, isPayline := matrix.IsPayline(pair.Value)
		if !isPayline {
			continue
		}
		payline := &pb.Payline{
			// Id: int32(idx),
			Indexs: make([]int32, 0),
		}
		payline.Id = int32(pair.Key)
		payline.Symbol = matrix.List[paylineIndexs[0]]
		payline.NumOccur = int32(len(paylineIndexs))
		for _, val := range paylineIndexs {
			payline.Indexs = append(payline.Indexs, int32(val))
		}
		paylines = append(paylines, payline)
	}
	return paylines
}

func (e *normal) countFreeSpinSymbolByCol(matrix entity.SlotMatrix) int {
	count := 0
	matrix.ForEachCol(func(col int, symbols []pb.SiXiangSymbol) {
		for _, symbol := range symbols {
			if symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_FREE_SPIN {
				count++
				break
			}
		}
	})
	return count
}
