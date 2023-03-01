package engine

import "github.com/ciaolink-game-platform/cgp-common/lib"

var _ lib.Engine = &goldPickEngine{}

type goldPickEngine struct {
}

func NewGoldPickEngine() lib.Engine {
	engine := goldPickEngine{}
	return &engine
}

func (e *goldPickEngine) NewGame(matchState interface{}) (interface{}, error) {
	return nil, nil
}

func (e *goldPickEngine) Random(min, max int) int {
	return 0
}

func (e *goldPickEngine) Process(matchState interface{}) (interface{}, error) {
	return nil, nil
}

func (e *goldPickEngine) Finish(matchState interface{}) (interface{}, error) {
	return nil, nil
}
