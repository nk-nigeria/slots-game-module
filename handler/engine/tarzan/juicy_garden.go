package tarzan

import "github.com/ciaolink-game-platform/cgp-common/lib"

var _ lib.Engine = &juiceGarden{}

type juiceGarden struct {
}

// NewGame implements lib.Engine
func (*juiceGarden) NewGame(matchState interface{}) (interface{}, error) {
	panic("unimplemented")
}

// Process implements lib.Engine
func (*juiceGarden) Process(matchState interface{}) (interface{}, error) {
	panic("unimplemented")
}

// Finish implements lib.Engine
func (*juiceGarden) Finish(matchState interface{}) (interface{}, error) {
	panic("unimplemented")
}

// Random implements lib.Engine
func (*juiceGarden) Random(min int, max int) int {
	panic("unimplemented")
}
