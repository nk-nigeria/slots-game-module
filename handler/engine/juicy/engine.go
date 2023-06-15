package juicy

import (
	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"

	pb "github.com/ciaolink-game-platform/cgp-common/proto"
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
	return e.engines[s.CurrentSiXiangGame].Finish(matchState)
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
