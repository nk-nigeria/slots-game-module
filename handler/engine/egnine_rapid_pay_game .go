package engine

import "github.com/ciaolink-game-platform/cgp-common/lib"

var _ lib.Engine = &rapidPayEngine{}

type rapidPayEngine struct {
}

func NewRapidPayEngine() lib.Engine {
	engine := rapidPayEngine{}
	return &engine
}

func (e *rapidPayEngine) NewGame(matchState interface{}) (interface{}, error) {
	return nil, nil
}

func (e *rapidPayEngine) Random(min, max int) int {
	return 0
}

func (e *rapidPayEngine) Process(matchState interface{}) (interface{}, error) {
	return nil, nil
}

func (e *rapidPayEngine) Finish(matchState interface{}) (interface{}, error) {
	return nil, nil
}
