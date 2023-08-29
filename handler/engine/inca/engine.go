package inca

import (
	"time"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

func NewEngine() lib.Engine {
	e := &engine{}
	e.engines = make(map[pb.SiXiangGame]lib.Engine)
	e.engines[pb.SiXiangGame_SI_XIANG_GAME_NORMAL] = NewNormal(nil)
	e.engines[pb.SiXiangGame_SI_XIANG_GAME_INCA_FREE_GAME] = NewFreeGame(nil)
	return e
}

// NewGame implements lib.Engine
func (e *engine) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	return e.engines[s.CurrentSiXiangGame].NewGame(matchState)
}

// Process implements lib.Engine
func (e *engine) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	s.IsSpinChange = true
	return e.engines[s.CurrentSiXiangGame].Process(s)
}
func (e *engine) Loop(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	engine := e.engines[s.CurrentSiXiangGame]
	return engine.Loop(s)
}

// Random implements lib.Engine
func (e *engine) Random(min int, max int) int {
	// return e.Random(min, max)
	return 0
}

type engine struct {
	engines map[pb.SiXiangGame]lib.Engine
}

// Finish implements lib.Engine.
func (e *engine) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	if !s.IsSpinChange {
		return s.LastResult, nil
	}
	s.IsSpinChange = false
	engine := e.engines[s.CurrentSiXiangGame]
	return engine.Finish(s)
}

// Info implements lib.Engine.
func (e *engine) Info(matchState interface{}) (interface{}, error) {
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
	if s.CurrentSiXiangGame == pb.SiXiangGame_SI_XIANG_GAME_INCA_FREE_GAME {
		slotdesk.GameConfig = s.GameConfig.GameConfig
	}
	return slotdesk, nil
}
