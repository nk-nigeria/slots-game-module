package juicy

import (
	"fmt"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

// todo implement chip win with % total chip win
var _ lib.Engine = &fruitRain{}

type fruitRain struct {
	randomIntFn func(min, max int) int
	// autoRefillGemSpin bool
	// matrixFruitRainBasket []pb.SiXiangSymbol
}

func NewFruitRain(randomIntFn func(min, max int) int) lib.Engine {
	e := &fruitRain{}
	if randomIntFn == nil {
		e.randomIntFn = entity.RandomInt
	} else {
		e.randomIntFn = randomIntFn
	}
	return e
}

// NewGame implements lib.Engine
func (e *fruitRain) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	s.WinJp = pb.WinJackpot_WIN_JACKPOT_UNSPECIFIED
	if s.NumSpinLeft <= 0 {
		matrixSpecial := entity.NewJuicyMatrix()
		s.MatrixSpecial = &matrixSpecial
		s.NumSpinLeft = 3
		s.GameConfig.AddGiftSpin = true
		m := entity.NewJuicyFruitRainMaxtrix()
		s.WildMatrix = m
		s.SpinList = make([]*pb.SpinSymbol, 0)
		s.MatrixSpecial.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
			s.SpinList = append(s.SpinList, &pb.SpinSymbol{
				Symbol:    pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED,
				Row:       int32(row),
				Col:       int32(col),
				Index:     int32(idx),
				Ratio:     0,
				WinJp:     pb.WinJackpot_WIN_JACKPOT_UNSPECIFIED,
				WinAmount: 0,
			})
			for {
				randomSymbol := entity.JuicySpinSymbol(e.randomIntFn, entity.JuiceAllSymbols)
				if entity.IsFruitBasketSymbol(randomSymbol) {
					continue
				}
				s.MatrixSpecial.List[idx] = randomSymbol
				s.MatrixSpecial.Flip(idx)
				break
			}
		})
		s.ChipStat.Reset(s.CurrentSiXiangGame)
	}
	return s, nil
}

// Process implements lib.Engine
func (e *fruitRain) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	if s.NumSpinLeft <= 0 {
		return s, entity.ErrorSpinReachMax
	}
	s.IsSpinChange = true
	// matrix := s.MatrixSpecial
	matrix := e.SpinMatrix(*s.MatrixSpecial)
	s.SpinSymbols = make([]*pb.SpinSymbol, 0)
	s.NumSpinLeft--

	s.MatrixSpecial.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		// keep symbol if fruitbasket
		if entity.IsFruitBasketSymbol(symbol) {
			return
		}
		newSymbol := matrix.List[idx]
		if entity.IsFruitBasketSymbol(newSymbol) {
			fmt.Println("")
		}
		if newSymbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_SPIN {
			newSymbol = s.WildMatrix.List[idx]
		}
		s.MatrixSpecial.List[idx] = newSymbol
		s.MatrixSpecial.Flip(idx)

		if entity.IsFruitBasketSymbol(newSymbol) {
			s.NumSpinLeft = 3
			s.GameConfig.AddGiftSpin = false
			spinSymbol := s.SpinList[idx]
			spinSymbol.Symbol = newSymbol
			val := entity.JuicyBasketSymbol[spinSymbol.Symbol]
			spinSymbol.Ratio = float32(e.randomIntFn(int(val.Value.Min), int(val.Value.Max)))
			spinSymbol.WinJp = entity.JuicySpinSymbolToJp(spinSymbol.Symbol)
			// s.SpinList = append(s.SpinList, spinSymbol)
			s.SpinSymbols = append(s.SpinSymbols, spinSymbol)
		}

	})
	// Nếu trong 3 lượt đầu tiên mà không có Giỏ trái cây nào được thêm mới vào màn hình,
	//  user được tặng 3 lượt freespin nữa. việc tặng chỉ xuất hiện 1 lần duy nhất.
	if s.NumSpinLeft == 0 && s.GameConfig.AddGiftSpin {
		s.NumSpinLeft = 3
		s.GameConfig.AddGiftSpin = false
	}
	return s, nil
}

// Random implements lib.Engine
func (e *fruitRain) Random(min int, max int) int {
	return e.randomIntFn(min, max)
}

// Finish implements lib.Engine
func (e *fruitRain) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	if !s.IsSpinChange {
		return s.LastResult, nil
	}
	s.IsSpinChange = false

	isFinish := s.NumSpinLeft == 0
	numFruitBasket := 0
	s.MatrixSpecial.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		if entity.IsFruitBasketSymbol(symbol) {
			numFruitBasket++
		}
	})
	lineWin := int64(0)
	for _, spin := range s.SpinSymbols {
		lineWin += int64(spin.Ratio)
		if entity.IsFruitJPSymbol(spin.Symbol) {
			spin.WinAmount = 0
			jpHistory := s.WinJPHistoryJuice()
			switch spin.WinJp {
			case pb.WinJackpot_WIN_JACKPOT_MINI:
				spin.WinAmount = jpHistory.Minor.Chips
			case pb.WinJackpot_WIN_JACKPOT_MINOR:
				spin.WinAmount = jpHistory.Minor.Chips
			case pb.WinJackpot_WIN_JACKPOT_MAJOR:
				spin.WinAmount = jpHistory.Major.Chips
				s.GetAndResetChipAccumt(pb.WinJackpot_WIN_JACKPOT_MAJOR)
			}
		} else {
			spin.WinAmount = int64(spin.Ratio) * s.Bet().Chips / 20
		}
		s.SpinList[spin.Index].WinAmount = spin.WinAmount
	}
	chipWin := lineWin * s.Bet().Chips / 20
	totalChipWin := int64(0)
	totalLineWin := int64(0)
	for _, spin := range s.SpinList {
		totalChipWin += spin.WinAmount
		totalLineWin += int64(spin.Ratio)
	}
	if numFruitBasket == len(s.MatrixSpecial.List) {
		isFinish = true
		s.WinJp = pb.WinJackpot_WIN_JACKPOT_GRAND
		totalChipWin += s.WinJPHistory().Grand.Chips
		s.GetAndResetChipAccumt(pb.WinJackpot_WIN_JACKPOT_GRAND)
	}
	if s.WinJp != pb.WinJackpot_WIN_JACKPOT_UNSPECIFIED {
		isFinish = true
	}
	s.NextSiXiangGame = s.CurrentSiXiangGame
	if isFinish {
		s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_NORMAL
		s.WildMatrix = entity.SlotMatrix{}
		s.NumSpinLeft = 0
	}
	slotDesk := &pb.SlotDesk{
		ChipsMcb: s.Bet().Chips,
		Matrix:   s.MatrixSpecial.ToPbSlotMatrix(),
		GameReward: &pb.GameReward{
			ChipsWin:            chipWin,
			LineWin:             lineWin,
			RatioWin:            float32(lineWin) / 100.0,
			TotalLineWin:        totalLineWin,
			TotalChipsWinByGame: totalChipWin,
			TotalRatioWin:       float32(totalLineWin) / 100.0,
		},
		CurrentSixiangGame: s.CurrentSiXiangGame,
		NextSixiangGame:    s.NextSiXiangGame,
		NumSpinLeft:        int64(s.NumSpinLeft),
		IsFinishGame:       isFinish,
	}
	slotDesk.Matrix.SpinLists = s.SpinList
	s.ChipStat.AddChipWin(s.CurrentSiXiangGame, slotDesk.GameReward.TotalChipsWinByGame)
	s.LastResult = slotDesk
	return slotDesk, nil
}

func (e *fruitRain) Loop(s interface{}) (interface{}, error) {
	return s, nil
}

func (e *fruitRain) Info(s interface{}) (interface{}, error) {
	return s, nil
}

func (e *fruitRain) SpinMatrix(matrix entity.SlotMatrix) entity.SlotMatrix {
	spinMatrix := entity.NewSlotMatrix(matrix.Rows, matrix.Cols)
	spinMatrix.List = make([]pb.SiXiangSymbol, spinMatrix.Size)
	spinMatrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		for {
			randomSymbol := entity.JuicySpinSymbol(e.randomIntFn, entity.JuiceAllSymbols)
			if entity.IsFruitJPSymbol(randomSymbol) {
				continue
			}
			spinMatrix.Flip(idx)
			spinMatrix.List[idx] = randomSymbol
			break
		}
	})
	return spinMatrix
}
