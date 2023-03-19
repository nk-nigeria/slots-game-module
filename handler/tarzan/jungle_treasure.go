package tarzan

import "github.com/ciaolink-game-platform/cgp-common/lib"

var _ lib.Engine = &jungleTreasure{}

type jungleTreasure struct{}

func NewJungleTrease() lib.Engine {
	return &jungleTreasure{}
}

// Finish implements lib.Engine
func (*jungleTreasure) Finish(matchState interface{}) (interface{}, error) {
	panic("unimplemented")
}

// NewGame implements lib.Engine
func (*jungleTreasure) NewGame(matchState interface{}) (interface{}, error) {
	panic("unimplemented")
}

// Process implements lib.Engine
func (*jungleTreasure) Process(matchState interface{}) (interface{}, error) {
	panic("unimplemented")
}

// Random implements lib.Engine
func (*jungleTreasure) Random(min int, max int) int {
	panic("unimplemented")
}
