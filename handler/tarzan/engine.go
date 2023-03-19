package tarzan

import "github.com/ciaolink-game-platform/cgp-common/lib"

var _ lib.Engine = &tarzanEngine{}

type tarzanEngine struct{}

// Finish implements lib.Engine
func (*tarzanEngine) Finish(matchState interface{}) (interface{}, error) {
	panic("unimplemented")
}

// NewGame implements lib.Engine
func (*tarzanEngine) NewGame(matchState interface{}) (interface{}, error) {
	panic("unimplemented")
}

// Process implements lib.Engine
func (*tarzanEngine) Process(matchState interface{}) (interface{}, error) {
	panic("unimplemented")
}

// Random implements lib.Engine
func (*tarzanEngine) Random(min int, max int) int {
	panic("unimplemented")
}
