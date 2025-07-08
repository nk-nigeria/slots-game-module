package tarzan

import (
	"github.com/nk-nigeria/slots-game-module/entity"
	"github.com/nk-nigeria/cgp-common/lib"
	pb "github.com/nk-nigeria/cgp-common/proto"
)

var _ lib.Engine = &jungleTreasure{}

type jungleTreasure struct {
	randomIntFn func(int, int) int
	// đảm bảo 100% sẽ có 1 lần ra turn x3
}

func NewJungTreasure(randomIntFn func(int, int) int) lib.Engine {
	e := &jungleTreasure{}
	if randomIntFn == nil {
		e.randomIntFn = entity.RandomInt
	} else {
		e.randomIntFn = randomIntFn
	}
	return e
}

// NewGame implements lib.Engine
func (e *jungleTreasure) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	if s.MatrixSpecial == nil {
		matrixSpecial := entity.NewTarzanJungleTreasureMatrix()
		matrixSpecial = entity.ShuffleMatrix(matrixSpecial)
		s.MatrixSpecial = &matrixSpecial
		s.SpinSymbols = nil
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
		})
		s.ChipStat.Reset(s.CurrentSiXiangGame)
	}
	if s.NumSpinLeft <= 0 {
		s.NumSpinLeft = 5
		s.TurnSureSpinSpecial = e.randomIntFn(1, s.NumSpinLeft+1)
	}

	return s, nil
}

// Process implements lib.Engine
func (e *jungleTreasure) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	if s.NumSpinLeft == 0 {
		return nil, entity.ErrorSpinReachMax
	}
	spinIndex := s.Bet().Id
	if spinIndex < 0 {
		return nil, entity.ErrorSpinIndexRequired
	}
	if s.MatrixSpecial.IsFlip(int(spinIndex)) {
		return nil, entity.ErrorSpinIndexAleadyTaken
	}
	if s.Bet().Id < 0 {
		return s, entity.ErrorInfoBetInvalid
	}
	if s.MatrixSpecial.IsFlip(int(s.Bet().Id)) {
		return s, entity.ErrorInfoBetInvalid
	}
	s.IsSpinChange = true
	var randIdx int
	var randSymbol pb.SiXiangSymbol
	if s.NumSpinLeft != s.TurnSureSpinSpecial {
		randIdx, randSymbol = int(spinIndex), s.MatrixSpecial.Flip(int(spinIndex))
	} else {
		// 100% spin symbol SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN_MORE_TURNX3
		s.MatrixSpecial.ForEeachNotFlip(func(idx, row, col int, symbol pb.SiXiangSymbol) {
			if symbol != pb.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN_MORE_TURNX3 {
				return
			}
			// swap
			s.MatrixSpecial.List[spinIndex], s.MatrixSpecial.List[idx] = s.MatrixSpecial.List[idx], s.MatrixSpecial.List[spinIndex]
			randIdx, randSymbol = int(spinIndex), s.MatrixSpecial.Flip(int(spinIndex))
		})
	}
	if randSymbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN_MORE_TURNX3 {
		// nếu đã quay được turnx3,
		// ko cần force ra turnx3 nữa
		s.TurnSureSpinSpecial = -1
	}
	s.MatrixSpecial.Flip(randIdx)
	row, col := s.MatrixSpecial.RowCol(randIdx)
	spin := &pb.SpinSymbol{
		Symbol: randSymbol,
		Row:    int32(row),
		Col:    int32(col),
		Index:  int32(randIdx),
	}
	symInfo := entity.TarzanJungleTreasureSymbol[spin.Symbol]
	spin.Ratio = float32(e.randomIntFn(int(symInfo.Value.Min), int(symInfo.Value.Max))) / float32(100)
	s.SpinSymbols = []*pb.SpinSymbol{spin}
	s.SpinList[randIdx] = spin
	s.NumSpinLeft--
	return matchState, nil
}

// Finish implements lib.Engine
func (e *jungleTreasure) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	if !s.IsSpinChange {
		return nil, entity.ErrorSpinNotChange
	}
	lineWin := 0
	chipsWin := 0
	totalLineWin := 0
	totalChipsWinByGame := 0
	for _, spin := range s.SpinSymbols {
		switch spin.Symbol {
		case pb.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN_MORE_TURNX2:
			s.NumSpinLeft += 2
		case pb.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN_MORE_TURNX3:
			s.NumSpinLeft += 3
		}
		spin.WinAmount = int64(float64(s.Bet().Chips) * float64(spin.Ratio))
		s.SpinList[spin.Index].WinAmount = spin.WinAmount
		lineWin += int(spin.Ratio * 100)
		chipsWin += int(spin.WinAmount)
	}
	for _, spin := range s.SpinList {
		totalLineWin += int(spin.GetRatio() * 100)
		totalChipsWinByGame += int(spin.WinAmount)
	}

	if s.NumSpinLeft <= 0 {
		s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_NORMAL
	}
	slotDesk := &pb.SlotDesk{
		Matrix:             s.MatrixSpecial.ToPbSlotMatrix(),
		CurrentSixiangGame: s.CurrentSiXiangGame,
		NextSixiangGame:    s.NextSiXiangGame,
		ChipsMcb:           s.Bet().GetChips(),
		IsFinishGame:       s.NumSpinLeft == 0,
		NumSpinLeft:        int64(s.NumSpinLeft),
		GameReward: &pb.GameReward{
			ChipsWin:            int64(chipsWin),
			LineWin:             int64(lineWin),
			RatioWin:            float32(lineWin / 100),
			TotalLineWin:        int64(totalLineWin),
			TotalRatioWin:       float32(totalLineWin / 100),
			TotalChipsWinByGame: int64(totalChipsWinByGame),
		},
		SpinSymbols: s.SpinSymbols,
	}
	s.MatrixSpecial.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		if !s.MatrixSpecial.TrackFlip[idx] {
			slotDesk.Matrix.Lists[idx] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED
		}
	})
	// reset letter
	if slotDesk.IsFinishGame {
		s.LetterSymbol = make(map[pb.SiXiangSymbol]bool)
		s.MatrixSpecial = nil
		s.SpinSymbols = nil
	}
	slotDesk.Matrix.SpinLists = s.SpinList
	s.ChipStat.Reset(s.CurrentSiXiangGame)
	s.ChipStat.AddChipWin(s.CurrentSiXiangGame, slotDesk.GameReward.TotalChipsWinByGame)
	s.ChipStat.AddLineWin(s.CurrentSiXiangGame, slotDesk.GameReward.TotalLineWin)
	s.LastResult = slotDesk
	return slotDesk, nil
}

// Random implements lib.Engine
func (e *jungleTreasure) Random(min int, max int) int {
	return e.randomIntFn(min, max)
}

func (e *jungleTreasure) Loop(s interface{}) (interface{}, error) {
	return s, nil
}

func (e *jungleTreasure) Info(matchState interface{}) (interface{}, error) {
	return nil, nil
}
