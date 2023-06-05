package juicy

import (
	"testing"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
	"github.com/stretchr/testify/assert"
)

func Test_fruitRain_Process(t *testing.T) {
	name := "Test_fruitRain_Process"
	t.Run(name, func(t *testing.T) {
		for i := 0; i < 1000; i++ {
			e := NewFruitRain(nil)
			s := entity.NewSlotsMathState(nil)
			e.NewGame(s)
			mapFruitbasketIndex := make(map[int]pb.SiXiangSymbol)
			for {
				_, err := e.Process(s)
				// check index of fruitbase not change after process
				for idx, symbol := range mapFruitbasketIndex {
					assert.Equal(t, symbol, s.MatrixSpecial.List[idx])
				}
				s.MatrixSpecial.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
					if entity.IsFruitBasketSymbol(symbol) {
						mapFruitbasketIndex[idx] = symbol
					}
				})
				if err != nil {
					break
				}
			}
		}
	})
}
