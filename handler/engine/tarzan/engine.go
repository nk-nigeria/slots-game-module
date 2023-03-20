package tarzan

import (
	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

var _ lib.Engine = &tarzanEngine{}

type tarzanEngine struct {
	engines     map[pb.SiXiangGame]lib.Engine
}

func NewEngine() lib.Engine {
	e := &tarzanEngine{
		engines: make(map[pb.SiXiangGame]lib.Engine),
	}
	e.engines[pb.SiXiangGame_SI_XIANG_GAME_TARZAN_NORMAL] = NewNormal(nil)
	e.engines[pb.SiXiangGame_SI_XIANG_GAME_TARZAN_JUNGLE_TREASURE] = NewJungleTrease()
	return e
}

// NewGame implements lib.Engine
func (e *tarzanEngine) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.TarzanMatchState)
	engine := e.engines[s.CurrentSiXiangGame]
	engine.NewGame(matchState)
	return matchState, nil
}

// Process implements lib.Engine
func (e *tarzanEngine) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.TarzanMatchState)
	engine := e.engines[s.CurrentSiXiangGame]
	return engine.Process(matchState)
}

// Random implements lib.Engine
func (e *tarzanEngine) Random(min int, max int) int {
	panic("unimplemented")
}

// Finish implements lib.Engine
func (e *tarzanEngine) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.TarzanMatchState)
	engine := e.engines[s.CurrentSiXiangGame]
	return engine.Finish(matchState)
}


