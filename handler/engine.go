package handler

import (
	"fmt"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
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

func (e *slotsEngine) Process(matchState interface{}) (interface{}, error) {
	return nil, nil
}
func (e *slotsEngine) Finish(matchState interface{}) (interface{}, error) {
	return nil, nil
}

func (e *slotsEngine) SpinMatrix(matchState *entity.SlotsMatchState) {
	matrix, cols, rows := matchState.GetMatrix()
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			matrix[row][col] = e.Random(0, len(entity.ListSymbol))
		}
	}
	matchState.SetMatrix(matrix)
}

func (e *slotsEngine) PrintMatrix(matchState *entity.SlotsMatchState) {
	matrix, cols, rows := matchState.GetMatrix()
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			fmt.Printf("%5d", matrix[row][col])
		}
		fmt.Println("")
	}
}
