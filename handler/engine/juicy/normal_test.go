package juicy

import (
	"testing"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	api "github.com/ciaolink-game-platform/cgp-common/proto"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
	"github.com/stretchr/testify/assert"
)

func Test_normal_SpinMatrix(t *testing.T) {
	name := "Test_normal_SpinMatrix"
	e := NewNormal(nil)
	engine := e.(*normal)
	numWildRatio1_0 := 0
	numWildRatio1_2 := 0
	numWildRatio1_5 := 0
	numWildRatio2_0 := 0
	m := entity.NewSlotMatrix(5, 3)
	countWildSymbol := func(matrix entity.SlotMatrix) int {
		num := 0
		matrix.ForEeach(func(idx, row, col int, symbol api.SiXiangSymbol) {
			if symbol == api.SiXiangSymbol_SI_XIANG_SYMBOL_WILD {
				num++
			}
		})
		return num
	}
	checkWildPositition := func(matrix entity.SlotMatrix) {
		matrix.ForEachCol(func(col int, symbols []pb.SiXiangSymbol) {
			numWild := 0
			for _, symbol := range symbols {
				if symbol == api.SiXiangSymbol_SI_XIANG_SYMBOL_WILD {
					numWild++
				}
				assert.Equal(t, true, numWild == 0 || col != entity.Col_1)
			}
		})
	}
	t.Run(name, func(t *testing.T) {
		for i := 0; i < 10000; i++ {
			matrix := engine.SpinMatrix(m, ratioWild1_0)
			numWildRatio1_0 += countWildSymbol(matrix)
			checkWildPositition(matrix)

			matrix1_2 := engine.SpinMatrix(m, ratioWild1_2)
			numWildRatio1_2 += countWildSymbol(matrix1_2)
			checkWildPositition(matrix1_2)

			matrix1_5 := engine.SpinMatrix(m, ratioWild1_5)
			numWildRatio1_5 += countWildSymbol(matrix1_5)
			checkWildPositition(matrix1_5)

			matrix2_0 := engine.SpinMatrix(m, ratioWild2_0)
			numWildRatio2_0 += countWildSymbol(matrix2_0)
			checkWildPositition(matrix2_0)
		}

		assert.Less(t, ratioWild1_0, ratioWild1_2)
		assert.Less(t, ratioWild1_2, ratioWild1_5)
		assert.Less(t, ratioWild1_5, ratioWild2_0)
		t.Logf("ratio 1.0 %d 1.2 %d 1.5 %d 2.0 %d \r\n", numWildRatio1_0, numWildRatio1_2, numWildRatio1_5, numWildRatio2_0)
	})
}

func TestNewNormal(t *testing.T) {
	name := "TestNewNormal"
	t.Run(name, func(t *testing.T) {
		e := NewNormal(nil)
		assert.NotNil(t, e)
	})
}

func Test_normal_NewGame(t *testing.T) {
	name := "Test_normal_NewGame"
	t.Run(name, func(t *testing.T) {
		e := NewNormal(nil)
		s := entity.NewSlotsMathState(nil)
		_, err := e.NewGame(s)
		t.Logf("%v", s.Matrix)
		assert.NoError(t, err)
		mapAllowSymbol := make(map[pb.SiXiangSymbol]bool)
		for _, symbol := range entity.JuiceAllSymbols {
			mapAllowSymbol[symbol] = true
		}
		t.Logf("matrix %v", s.Matrix.List)
		assert.Equal(t, entity.RowsJuicynMatrix*entity.ColsJuicyMatrix, len(s.Matrix.List))
		s.Matrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
			assert.Equal(t, true, mapAllowSymbol[symbol])
		})

	})
}

func Test_normal_Process(t *testing.T) {
	name := "Test_normal_Process"
	engine := NewNormal(nil)
	s := entity.NewSlotsMathState(nil)
	_, _ = engine.NewGame(s)
	t.Run(name, func(t *testing.T) {
		_, err := engine.Process(s)
		assert.NoError(t, err)
		assert.Equal(t, entity.RowsJuicynMatrix*entity.ColsJuicyMatrix, len(s.Matrix.List))
		assert.Equal(t, entity.RowsJuicynMatrix*entity.ColsJuicyMatrix, len(s.WildMatrix.List))
		assert.NotNil(t, s.Paylines())
	})
}

func Test_normal_Paylines(t *testing.T) {
	name := "Test_normal_Paylines"
	e := NewNormal(nil)
	engine := e.(*normal)
	newEmptyMatrix := func() entity.SlotMatrix {
		matrix := entity.NewJuicyMatrix()
		matrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
			matrix.List[idx] = api.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED
		})
		return matrix

	}
	t.Run(name, func(t *testing.T) {
		for pair := entity.MapJuicyPaylineIdx.Oldest(); pair != nil; pair = pair.Next() {
			ids := pair.Value
			// 3 symbol xxxyy
			{
				matrix := newEmptyMatrix()
				matrix.List[ids[0]] = api.SiXiangSymbol_SI_XIANG_SYMBOL_J
				matrix.List[ids[1]] = api.SiXiangSymbol_SI_XIANG_SYMBOL_J
				matrix.List[ids[2]] = api.SiXiangSymbol_SI_XIANG_SYMBOL_J
				paylines := engine.Paylines(matrix)
				assert.NotNil(t, paylines)
				assert.Equal(t, 1, len(paylines))
				payline := paylines[0]
				assert.Equal(t, int32(3), payline.NumOccur)
				assert.Equal(t, api.SiXiangSymbol_SI_XIANG_SYMBOL_J, payline.GetSymbol())
				assert.Equal(t, int(pair.Key), int(payline.GetId()))
			}
			// 3 symbol xwildxyy
			{
				matrix := newEmptyMatrix()
				matrix.List[ids[0]] = api.SiXiangSymbol_SI_XIANG_SYMBOL_J
				matrix.List[ids[1]] = api.SiXiangSymbol_SI_XIANG_SYMBOL_WILD
				matrix.List[ids[2]] = api.SiXiangSymbol_SI_XIANG_SYMBOL_J
				paylines := engine.Paylines(matrix)
				assert.NotNil(t, paylines)
				assert.Equal(t, 1, len(paylines))
				payline := paylines[0]
				assert.Equal(t, int32(3), payline.NumOccur)
				assert.Equal(t, api.SiXiangSymbol_SI_XIANG_SYMBOL_J, payline.GetSymbol())
				assert.Equal(t, int(pair.Key), int(payline.GetId()))
			}
			// 3 symbol xxwildyy
			{
				matrix := newEmptyMatrix()
				matrix.List[ids[0]] = api.SiXiangSymbol_SI_XIANG_SYMBOL_J
				matrix.List[ids[1]] = api.SiXiangSymbol_SI_XIANG_SYMBOL_J
				matrix.List[ids[2]] = api.SiXiangSymbol_SI_XIANG_SYMBOL_WILD
				paylines := engine.Paylines(matrix)
				assert.NotNil(t, paylines)
				assert.Equal(t, 1, len(paylines))
				payline := paylines[0]
				assert.Equal(t, int32(3), payline.NumOccur)
				assert.Equal(t, api.SiXiangSymbol_SI_XIANG_SYMBOL_J, payline.GetSymbol())
				assert.Equal(t, int(pair.Key), int(payline.GetId()))
			}
			// 4 symbol xxxxy
			{
				matrix := newEmptyMatrix()
				matrix.List[ids[0]] = api.SiXiangSymbol_SI_XIANG_SYMBOL_J
				matrix.List[ids[1]] = api.SiXiangSymbol_SI_XIANG_SYMBOL_J
				matrix.List[ids[2]] = api.SiXiangSymbol_SI_XIANG_SYMBOL_J
				matrix.List[ids[3]] = api.SiXiangSymbol_SI_XIANG_SYMBOL_J
				paylines := engine.Paylines(matrix)
				assert.NotNil(t, paylines)
				assert.Equal(t, 1, len(paylines))
				payline := paylines[0]
				assert.Equal(t, int32(4), payline.NumOccur)
				assert.Equal(t, api.SiXiangSymbol_SI_XIANG_SYMBOL_J, payline.GetSymbol())
				assert.Equal(t, int(pair.Key), int(payline.GetId()))
			}
			// 5 symbol xxxxx
			{
				matrix := newEmptyMatrix()
				matrix.List[ids[0]] = api.SiXiangSymbol_SI_XIANG_SYMBOL_J
				matrix.List[ids[1]] = api.SiXiangSymbol_SI_XIANG_SYMBOL_J
				matrix.List[ids[2]] = api.SiXiangSymbol_SI_XIANG_SYMBOL_J
				matrix.List[ids[3]] = api.SiXiangSymbol_SI_XIANG_SYMBOL_J
				matrix.List[ids[4]] = api.SiXiangSymbol_SI_XIANG_SYMBOL_J
				paylines := engine.Paylines(matrix)
				assert.NotNil(t, paylines)
				assert.Equal(t, 1, len(paylines))
				payline := paylines[0]
				assert.Equal(t, int32(5), payline.NumOccur)
				assert.Equal(t, api.SiXiangSymbol_SI_XIANG_SYMBOL_J, payline.GetSymbol())
				assert.Equal(t, int(pair.Key), int(payline.GetId()))
			}
			// 2 symbol xxyyy
			arrAllow := []pb.SiXiangSymbol{api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STONE_DIAMOND,
				api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STONE_VIOLET,
				api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STONE_GREEN,
			}
			for _, symbol := range arrAllow {
				matrix := newEmptyMatrix()
				matrix.List[ids[0]] = symbol
				matrix.List[ids[1]] = symbol

				paylines := engine.Paylines(matrix)
				assert.NotNil(t, paylines)
				t.Logf("paylinens %v", paylines)
				payline := paylines[0]
				assert.Equal(t, symbol, payline.GetSymbol())
			}
			arrNotAllow := []pb.SiXiangSymbol{
				api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STRAWBERRY,
				api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STRAWBERRY,
				api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_PINAPPLE,
				api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_MANGOSTEEN,
				api.SiXiangSymbol_SI_XIANG_SYMBOL_A,
				api.SiXiangSymbol_SI_XIANG_SYMBOL_J,
				api.SiXiangSymbol_SI_XIANG_SYMBOL_Q,
				api.SiXiangSymbol_SI_XIANG_SYMBOL_K,
			}
			for _, symbol := range arrNotAllow {
				matrix := newEmptyMatrix()
				matrix.List[ids[0]] = symbol
				matrix.List[ids[1]] = symbol

				paylines := engine.Paylines(matrix)
				assert.NotNil(t, paylines)
				assert.Equal(t, int(0), len(paylines))

			}
		}
	})
}

func Test_normal_GetNextSiXiangGame(t *testing.T) {
	name := "Test_normal_GetNextSiXiangGame"
	e := NewNormal(nil)
	engine := e.(*normal)
	s := entity.NewSlotsMathState(nil)
	s.Matrix = entity.NewJuicyMatrix()
	for i := 0; i < 3; i++ {
		s.Matrix.List[i] = api.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER
	}
	t.Run(name, func(t *testing.T) {
		nextGame := engine.GetNextSiXiangGame(s.Matrix)
		assert.Equal(t, pb.SiXiangGame_SI_XIANG_GAME_JUICE_FRUIT_BASKET, nextGame)
		nextGame = engine.GetNextSiXiangGame(entity.NewJuicyMatrix())
		assert.Equal(t, pb.SiXiangGame_SI_XIANG_GAME_NORMAL, nextGame)

	})
}
