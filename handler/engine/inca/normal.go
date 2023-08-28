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
	if s.CurrentSiXiangGame == pb.SiXiangGame_SI_XIANG_GAME_NORMAL {
		s.ChipStat.ResetChipWin(s.CurrentSiXiangGame)
		s.ChipStat.ResetChipWin(pb.SiXiangGame_SI_XIANG_GAME_INCA_FREE_GAME)
		s.GameConfig.NumScatterSeq = int64(e.countScatterByCol(s.Matrix))
	}
	s.ChipStat.AddChipWin(s.CurrentSiXiangGame, totalWin)
	s.NextSiXiangGame = e.GetNextSiXiangGame(s)
	slotDesk := &pb.SlotDesk{
		GameReward: &pb.GameReward{
			ChipsWin:            s.ChipStat.ChipWin(s.CurrentSiXiangGame),
			TotalChipsWinByGame: s.ChipStat.TotalChipWin(s.CurrentSiXiangGame),
		},
		ChipsMcb:           s.Bet().Chips,
		CurrentSixiangGame: s.CurrentSiXiangGame,
		NextSixiangGame:    s.NextSiXiangGame,
		Matrix:             s.Matrix.ToPbSlotMatrix(),
		SpreadMatrix:       s.WildMatrix.ToPbSlotMatrix(),
		Paylines:           paylines,
		IsFinishGame:       s.NumSpinLeft <= 0,
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
		BetLevels:          entity.BetLevels[:],
	}
	return slotdesk, nil
}

// Loop implements lib.Engine.
func (*normal) Loop(matchState interface{}) (interface{}, error) {
	return nil, nil
}

// NewGame implements lib.Engine.
func (e *normal) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	matrix := entity.NewSlotMatrix(entity.RowsIncaMatrix, entity.ColsIncaMatrix)
	matrix.List = make([]pb.SiXiangSymbol, matrix.Size)
	matrix = e.SpinMatrix(matrix)
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
	s.SetPaylines(e.Paylines(s.Matrix))
	return s, nil
}

// Random implements lib.Engine.
func (e *normal) Random(min int, max int) int {
	return e.randomFn(min, max)
}

func (e *normal) SpinMatrix(matrix entity.SlotMatrix) entity.SlotMatrix {
	spinMatrix := entity.NewSlotMatrix(matrix.Rows, matrix.Cols)
	spinMatrix.List = make([]pb.SiXiangSymbol, spinMatrix.Size)
	var randSymbol pb.SiXiangSymbol
	listSymbols := entity.ShuffleSlice(entity.IncalAllSymbol)
	matrix.ForEeach(func(idx, _, col int, _ pb.SiXiangSymbol) {
		for {
			randSymbol = listSymbols[e.Random(0, len(listSymbols))]
			if randSymbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER {
				if col == entity.Col_1 || col == entity.Col_5 {
					continue
				}
			}
			spinMatrix.List[idx] = randSymbol
			break
		}
	})
	return spinMatrix
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

func (e *normal) GetNextSiXiangGame(s *entity.SlotsMatchState) pb.SiXiangGame {
	if s.GameConfig.NumScatterSeq >= 3 {
		return pb.SiXiangGame_SI_XIANG_GAME_INCA_FREE_GAME
	}
	if s.NumSpinLeft <= 0 {
		return pb.SiXiangGame_SI_XIANG_GAME_NORMAL
	}
	return s.CurrentSiXiangGame
}
