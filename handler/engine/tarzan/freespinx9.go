package tarzan

import (
	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

var _ lib.Engine = &freespinx9{}

type freespinx9 struct {
	normal
}

func NewFreeSpinX9(randomIntFn func(int, int) int) lib.Engine {
	e := NewNormal(randomIntFn)
	engine := e.(*normal)
	engine.allowSpinLetterSymbol = false
	engine.allowSpinFreeSpinSymbol = false
	engine.allowSpinTarzanSymbol = false
	freespinx9Engine := &freespinx9{
		normal: *engine,
	}
	return freespinx9Engine
}

// NewGame implements lib.Engine
func (*freespinx9) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.TarzanMatchState)
	s.ChipWinByGame[s.CurrentSiXiangGame] = 0
	s.CountLineCrossFreeSpinSymbol = 0
	s.GemSpin = 9
	return matchState, nil
}

func (e *freespinx9) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.TarzanMatchState)
	if s.GemSpin <= 0 {
		return nil, ErrorSpinReachMaximum
	}
	e.normal.Process(matchState)
	s.GemSpin--
	return matchState, nil
}

// Finish implements lib.Engine
func (e *freespinx9) Finish(matchState interface{}) (interface{}, error) {
	result, err := e.normal.Finish(matchState)
	if err != nil {
		return result, err
	}
	s := matchState.(*entity.TarzanMatchState)
	// check if payline pass freespin symbol
	for _, payline := range s.Paylines {
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
	s.ChipWinByGame[s.CurrentSiXiangGame] += slotDesk.ChipsWin
	// Finish when gem spin = 0
	slotDesk.IsFinishGame = s.GemSpin <= 0
	// tiền thưởng = (tổng số tiền thắng trong 9 Freespin) x (hệ số nhân bonus ở trên)
	slotDesk.TotalChipsWinByGame = s.ChipWinByGame[s.CurrentSiXiangGame] * int64(s.CountLineCrossFreeSpinSymbol)
	slotDesk.ChipsWin = slotDesk.TotalChipsWinByGame
	return slotDesk, err
}
