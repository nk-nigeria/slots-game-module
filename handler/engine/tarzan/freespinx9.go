package tarzan

import (
	"math"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

const maxGemSpinFreeSpinX9 = 9

var _ lib.Engine = &freespinx9{}

type freespinx9 struct {
	normal
}

func NewFreeSpinX9(randomIntFn func(int, int) int) lib.Engine {
	e := NewNormal(randomIntFn)
	engine := e.(*normal)
	engine.maxDropLetterSymbol = 0
	engine.maxDropFreeSpin = math.MaxInt
	engine.maxDropTarzanSymbol = 1
	engine.allowDropFreeSpinx9 = false
	freespinx9Engine := &freespinx9{
		normal: *engine,
	}
	return freespinx9Engine
}

// NewGame implements lib.Engine
func (e *freespinx9) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	// s.ChipWinByGame[s.CurrentSiXiangGame] = 0
	s.ChipStat.ResetChipWin(s.CurrentSiXiangGame)
	// s.LineWinByGame[s.CurrentSiXiangGame] = 0
	s.ChipStat.ResetLineWin(s.CurrentSiXiangGame)
	s.CountLineCrossFreeSpinSymbol = 0
	s.NumSpinLeft = maxGemSpinFreeSpinX9
	return matchState, nil
}

func (e *freespinx9) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	if s.NumSpinLeft <= 0 {
		return nil, entity.ErrorSpinReachMax
	}
	e.normal.Process(matchState)
	s.NumSpinLeft--
	return matchState, nil
}

// Finish implements lib.Engine
func (e *freespinx9) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	// prevChipWin := s.ChipWinByGame[s.CurrentSiXiangGame]
	prevChipWin := s.ChipStat.ChipWin(s.CurrentSiXiangGame)
	// prevLineWin := s.LineWinByGame[s.CurrentSiXiangGame]
	prevLineWin := s.ChipStat.LineWin(s.CurrentSiXiangGame)
	result, err := e.normal.Finish(matchState)
	if err != nil {
		return result, err
	}
	// check if payline pass freespin symbol
	for _, payline := range s.Paylines() {
		num := 0
		for _, val := range payline.Indexs {
			if s.TrackIndexFreeSpinSymbol[int(val)] {
				num++
			}
		}
		if num >= 3 {
			s.CountLineCrossFreeSpinSymbol++
		}
	}

	slotDesk := result.(*pb.SlotDesk)
	slotDesk.IsFinishGame = s.NumSpinLeft <= 0
	if slotDesk.IsFinishGame {
		// clean
		s.TrackIndexFreeSpinSymbol = make(map[int]bool)
		s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_NORMAL
	} else {
		s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_TARZAN_FREESPINX9
	}

	slotDesk.CurrentSixiangGame = s.CurrentSiXiangGame
	slotDesk.NextSixiangGame = s.NextSiXiangGame
	// s.ChipWinByGame[s.CurrentSiXiangGame] = s.ChipWinByGame[s.CurrentSiXiangGame] + prevChipWin
	s.ChipStat.AddChipWin(s.CurrentSiXiangGame, prevChipWin)
	// s.LineWinByGame[s.CurrentSiXiangGame] = s.LineWinByGame[s.CurrentSiXiangGame] + prevLineWin
	s.ChipStat.AddLineWin(s.CurrentSiXiangGame, prevLineWin)
	// Finish when gem spin = 0
	// tiền thưởng = (tổng số tiền thắng trong 9 Freespin) x (hệ số nhân bonus ở trên)
	// slotDesk.TotalChipsWinByGame = s.ChipWinByGame[s.CurrentSiXiangGame]
	slotDesk.TotalChipsWinByGame = s.ChipStat.ChipWin(s.CurrentSiXiangGame)
	if s.CountLineCrossFreeSpinSymbol > 0 {
		slotDesk.TotalChipsWinByGame *= int64(s.CountLineCrossFreeSpinSymbol)
	}
	slotDesk.ChipsWin = slotDesk.TotalChipsWinByGame
	slotDesk.NumSpinLeft = int64(s.NumSpinLeft)
	return slotDesk, err
}
