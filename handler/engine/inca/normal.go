package inca

import (
	"time"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

type normal struct {
	randomFn func(min, max int) int
}

func NewEngine() lib.Engine {
	e := NewNormal(nil)
	return e
}

func NewNormal(randomIntFn func(int, int) int) lib.Engine {
	e := &normal{}
	if randomIntFn != nil {
		e.randomFn = randomIntFn
	} else {
		e.randomFn = entity.RandomInt
	}
	return e
}

// Finish implements lib.Engine.
func (e *normal) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	paylines := s.Paylines()
	totalWin := int64(0)
	for _, payline := range paylines {
		payline.Chips = s.Bet().Chips * int64(payline.Rate) / 20
		totalWin += payline.Chips
	}
	{
		numScatter := e.countScatterByCol(s.Matrix)
		ratio := entity.IncaRatioPaylineMap[pb.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER][int32(numScatter)]
		totalWin += int64(ratio) * s.Bet().Chips / 20
	}
	slotDesk := &pb.SlotDesk{
		GameReward: &pb.GameReward{
			ChipsWin:            totalWin,
			TotalChipsWinByGame: totalWin,
			UpdateWallet:        true,
		},
		ChipsMcb:           s.Bet().Chips,
		CurrentSixiangGame: s.CurrentSiXiangGame,
		NextSixiangGame:    s.NextSiXiangGame,
		Matrix:             s.Matrix.ToPbSlotMatrix(),
		SpreadMatrix:       s.WildMatrix.ToPbSlotMatrix(),
		Paylines:           paylines,
		IsFinishGame:       true,
		NumSpinLeft:        int64(s.NumSpinLeft),
		BetLevels:          entity.BetLevels[:],
	}
	s.LastResult = slotDesk
	return slotDesk, nil
}

// Info implements lib.Engine.
func (*normal) Info(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	slotdesk := &pb.SlotDesk{
		Matrix:             s.Matrix.ToPbSlotMatrix(),
		SpreadMatrix:       s.WildMatrix.ToPbSlotMatrix(),
		ChipsMcb:           s.Bet().Chips,
		CurrentSixiangGame: s.CurrentSiXiangGame,
		NextSixiangGame:    s.NextSiXiangGame,
		TsUnix:             time.Now().Unix(),
		NumSpinLeft:        int64(s.NumSpinLeft),
		InfoBet:            s.Bet(),
	}
	return slotdesk, nil
}

// Loop implements lib.Engine.
func (*normal) Loop(matchState interface{}) (interface{}, error) {
	return nil, nil
}

// NewGame implements lib.Engine.
func (*normal) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	matrix := entity.NewSlotMatrix(entity.RowsIncaMatrix, entity.ColsIncaMatrix)
	s.SetMatrix(matrix)
	s.SetWildMatrix(matrix)
	s.NumSpinLeft = -1
	return matchState, nil
}

// Process implements lib.Engine.
func (e *normal) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	s.SetMatrix(e.SpinMatrix(s.Matrix))
	s.SetWildMatrix(e.SpreadWildInMatrix(s.Matrix))
	s.SetPaylines(e.Paylines(*s.MatrixSpecial))
	return s, nil
}

// Random implements lib.Engine.
func (e *normal) Random(min int, max int) int {
	return e.randomFn(min, max)
}

func (e *normal) SpinMatrix(matrix entity.SlotMatrix) entity.SlotMatrix {
	spinMatrix := entity.NewSlotMatrix(matrix.Rows, matrix.Cols)
	var randSymbol pb.SiXiangSymbol
	symbols := entity.ShuffleSlice(entity.IncalAllSymbol)
	matrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		for {
			randSymbol = symbols[e.Random(0, len(symbols))]
			if randSymbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_WILD && col < entity.Col_2 {
				continue
			}
			spinMatrix.List[idx] = randSymbol
			break
		}
	})
	return matrix
}

func (e *normal) Paylines(matrix entity.SlotMatrix) []*pb.Payline {
	paylines := make([]*pb.Payline, 0)
	for pair := entity.MapJuicyPaylineIdx.Oldest(); pair != nil; pair = pair.Next() {
		paylineIndexs, isPayline := matrix.IsIncaPayline(matrix, pair.Value)
		if !isPayline {
			continue
		}
		payline := &pb.Payline{
			Id:       int32(pair.Key),
			Symbol:   pb.SiXiangSymbol_SI_XIANG_SYMBOL_WILD,
			NumOccur: int32(len(paylineIndexs)),
		}
		for _, symIdx := range paylineIndexs {
			symbol := matrix.List[symIdx]
			if symbol != pb.SiXiangSymbol_SI_XIANG_SYMBOL_WILD {
				payline.Symbol = symbol
				break
			}
		}
		for _, val := range paylineIndexs {
			payline.Indices = append(payline.GetIndices(), int32(val))
		}
		payline.Rate = entity.IncaRatioPaylineMap[payline.Symbol][payline.NumOccur]
		if payline.Symbol != pb.SiXiangSymbol_SI_XIANG_SYMBOL_WILD {
			containSymbolWild := false
			for _, idx := range payline.Indices {
				sym := matrix.List[idx]
				if sym == pb.SiXiangSymbol_SI_XIANG_SYMBOL_WILD {
					containSymbolWild = true
					break
				}
			}
			if containSymbolWild {
				payline.Rate *= 2
			}
		}
		if payline.Rate > 0 {
			payline.Id = int32(pair.Key)
			paylines = append(paylines, payline)
		}
	}
	return paylines
}

func (e *normal) SpreadWildInMatrix(matrix entity.SlotMatrix) entity.SlotMatrix {
	return matrix
}

func (e *normal) countScatterByCol(matrix entity.SlotMatrix) int {
	count := 0
	matrix.ForEachCol(func(col int, symbols []pb.SiXiangSymbol) {
		for _, symbol := range symbols {
			if symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER {
				count++
				break
			}
		}
	})
	return count
}
