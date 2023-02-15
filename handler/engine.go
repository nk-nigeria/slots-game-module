package handler

import (
	"github.com/ciaolink-game-platform/cgp-common/lib"
	"github.com/ciaolink-game-platform/cgp-common/utilities"
)

var _ lib.Engine = &slotsEngine{}

type slotsEngine struct {
}

func NewSlotsEngine() lib.Engine {
	engine := slotsEngine{}
	return &engine
}

func (e *slotsEngine) NewGame(matchState interface{}) (interface{}, error) {
	return nil, nil
}

func (e *slotsEngine) Random(min, max int) int {
	return utilities.RandomNumber(min, max)
}

func (e *slotsEngine) Finish(matchState interface{}) interface{} {
	return nil
}
