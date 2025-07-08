package juicy

import (
	"time"

	"github.com/nk-nigeria/slots-game-module/entity"
	"github.com/nk-nigeria/cgp-common/lib"

	pb "github.com/nk-nigeria/cgp-common/proto"
)

var _ lib.Engine = &engine{}

type engine struct {
	engines map[pb.SiXiangGame]lib.Engine
}

func NewEngine() lib.Engine {
	e := &engine{}
	e.engines = make(map[pb.SiXiangGame]lib.Engine)
	e.engines[pb.SiXiangGame_SI_XIANG_GAME_NORMAL] = NewNormal(nil)
	e.engines[pb.SiXiangGame_SI_XIANG_GAME_JUICE_FRUIT_BASKET] = NewFruitBaseket()
	e.engines[pb.SiXiangGame_SI_XIANG_GAME_JUICE_FREE_GAME] = NewFreeGame(nil)
	e.engines[pb.SiXiangGame_SI_XIANG_GAME_JUICE_FRUIT_RAIN] = NewFruitRain(nil)
	return e
}

// Finish implements lib.Engine
func (e *engine) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	engine := e.engines[s.CurrentSiXiangGame]
	result, err := engine.Finish(matchState)
	if err != nil {
		return result, err
	}
	slotDesk := result.(*pb.SlotDesk)
	if slotDesk == nil {
		return result, err
	}
	slotDesk.WinJpHistory = s.WinJPHistoryJuice()
	slotDesk.GameConfig = nil
	if slotDesk.IsFinishGame {
		s.GameConfig = entity.GameConfigFreeGame(0)
		if s.CurrentSiXiangGame == pb.SiXiangGame_SI_XIANG_GAME_JUICE_FRUIT_BASKET && s.GameConfig != nil {
			slotDesk.GameConfig = s.GameConfig.GameConfig
		}
	}
	s.LastResult = slotDesk
	return slotDesk, nil
}

func (e *engine) Info(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	var matrix *pb.SlotMatrix
	var spreadMatrix *pb.SlotMatrix
	switch s.CurrentSiXiangGame {
	case pb.SiXiangGame_SI_XIANG_GAME_NORMAL:
		spreadMatrix = s.WildMatrix.ToPbSlotMatrix()
		matrix = spreadMatrix
	case pb.SiXiangGame_SI_XIANG_GAME_JUICE_FRUIT_BASKET:
		matrix = s.MatrixSpecial.ToPbSlotMatrix()
		for idx, symbol := range s.MatrixSpecial.List {
			if s.MatrixSpecial.IsFlip(idx) {
				matrix.Lists[idx] = symbol
			} else {
				matrix.Lists[idx] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED
			}
		}
		spreadMatrix = matrix
	default:
		matrix = s.MatrixSpecial.ToPbSlotMatrix()
		spreadMatrix = matrix

	}
	matrix.SpinLists = s.SpinList
	spreadMatrix.SpinLists = s.SpinList
	slotdesk := &pb.SlotDesk{
		Matrix:             matrix,
		SpreadMatrix:       spreadMatrix,
		ChipsMcb:           s.Bet().Chips,
		CurrentSixiangGame: s.CurrentSiXiangGame,
		NextSixiangGame:    s.NextSiXiangGame,
		TsUnix:             time.Now().Unix(),
		SpinSymbols:        s.SpinSymbols,
		NumSpinLeft:        int64(s.NumSpinLeft),
		InfoBet:            s.Bet(),
		WinJpHistory:       s.WinJPHistoryJuice(),
		BetLevels:          entity.BetLevels[:],
		GameReward: &pb.GameReward{
			UpdateWallet:        false,
			TotalChipsWinByGame: s.ChipStat.TotalChipWin(s.CurrentSiXiangGame),
			TotalLineWin:        s.ChipStat.TotalLineWin(s.CurrentSiXiangGame),
		},
		WinJp:      s.WinJp,
		GameConfig: nil,
	}
	return slotdesk, nil
}

// NewGame implements lib.Engine
func (e *engine) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	return e.engines[s.CurrentSiXiangGame].NewGame(matchState)
}

// Process implements lib.Engine
func (e *engine) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	return e.engines[s.CurrentSiXiangGame].Process(s)
}
func (e *engine) Loop(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	engine := e.engines[s.CurrentSiXiangGame]
	return engine.Loop(s)
}

// Random implements lib.Engine
func (e *engine) Random(min int, max int) int {
	return e.Random(min, max)
}
