package tarzan

import (
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
	engine.maxDropFreeSpin = 2
	engine.maxDropTarzanSymbol = 1
	// engine.allowDropFreeSpinx9 = false
	freespinx9Engine := &freespinx9{
		normal: *engine,
	}
	return freespinx9Engine
}

// NewGame implements lib.Engine
func (e *freespinx9) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	s.ChipStat.Reset(s.CurrentSiXiangGame)
	s.CountLineCrossFreeSpinSymbol = 0
	s.NumSpinLeft = maxGemSpinFreeSpinX9
	s.SpinList = make([]*pb.SpinSymbol, 0)
	return matchState, nil
}

func (e *freespinx9) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	if s.NumSpinLeft <= 0 {
		return nil, entity.ErrorSpinReachMax
	}
	if len(s.LetterSymbol) > 0 {
		s.LetterSymbol = make(map[pb.SiXiangSymbol]bool)
		s.SaveGameJson()
	}
	e.normal.Process(matchState)
	s.LetterSymbol = make(map[pb.SiXiangSymbol]bool)
	s.NumSpinLeft--
	return matchState, nil
}

// Finish implements lib.Engine
func (e *freespinx9) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	// prevChipWin := s.ChipStat.ChipWin(s.CurrentSiXiangGame)
	// prevLineWin := s.ChipStat.LineWin(s.CurrentSiXiangGame)
	result, err := e.normal.Finish(matchState)
	if err != nil {
		return result, err
	}
	// check if payline pass freespin symbol
	for _, payline := range s.Paylines() {
		num := 0
		for _, val := range payline.GetIndices() {
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
	s.ChipStat.AddChipWin(s.CurrentSiXiangGame, slotDesk.GameReward.ChipsWin)
	s.ChipStat.AddLineWin(s.CurrentSiXiangGame, slotDesk.GameReward.LineWin)
	// Finish when gem spin = 0
	// tiền thưởng = (tổng số tiền thắng trong 9 Freespin) x (hệ số nhân bonus ở trên)
	slotDesk.GameReward.ChipsWin = s.ChipStat.ChipWin(s.CurrentSiXiangGame)
	slotDesk.GameReward.LineWin = s.ChipStat.LineWin(s.CurrentSiXiangGame)
	slotDesk.GameReward.TotalChipsWinByGame = s.ChipStat.TotalChipWin(s.CurrentSiXiangGame)
	slotDesk.GameReward.TotalLineWin = s.ChipStat.TotalLineWin(s.CurrentSiXiangGame)
	if s.CountLineCrossFreeSpinSymbol < 1 {
		s.CountLineCrossFreeSpinSymbol = 1
	}
	slotDesk.GameReward.RatioBonus = float32(s.CountLineCrossFreeSpinSymbol)
	if slotDesk.GameReward.RatioBonus > 1 {
		slotDesk.GameReward.TotalChipsWinByGame *= int64(slotDesk.GameReward.RatioBonus)
	}
	slotDesk.LetterSymbols = make([]pb.SiXiangSymbol, 0)
	slotDesk.NumSpinLeft = int64(s.NumSpinLeft)
	return slotDesk, err
}

func (e *freespinx9) Loop(s interface{}) (interface{}, error) {
	return s, nil
}

func (e *freespinx9) Info(matchState interface{}) (interface{}, error) {
	return nil, nil
}
