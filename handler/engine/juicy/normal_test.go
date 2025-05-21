package juicy

import (
	"testing"

	"github.com/nakamaFramework/cgb-slots-game-module/entity"
	api "github.com/nakamaFramework/cgp-common/proto"
	pb "github.com/nakamaFramework/cgp-common/proto"
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
	numUnspecSymbol := func(matrix entity.SlotMatrix) int {
		num := 0
		matrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
			if symbol == api.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED {
				num++
			}
		})
		return num
	}
	checkValidSymbol := func(matrix entity.SlotMatrix) bool {
		valid := true
		matrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
			if symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_GRAND {
				valid = false
				return
			}
		})
		return valid
	}
	t.Run(name, func(t *testing.T) {
		for i := 0; i < 100000; i++ {
			matrix := engine.SpinMatrix(m, ratioWild1_0)
			numUnspecSymbol := numUnspecSymbol(matrix)
			assert.Equal(t, int(0), int(numUnspecSymbol))
			numWildRatio1_0 += countWildSymbol(matrix)
			checkWildPositition(matrix)
			assert.Equal(t, true, checkValidSymbol(matrix))

			matrix1_2 := engine.SpinMatrix(m, ratioWild1_2)
			numWildRatio1_2 += countWildSymbol(matrix1_2)
			checkWildPositition(matrix1_2)
			assert.Equal(t, true, checkValidSymbol(matrix1_2))

			matrix1_5 := engine.SpinMatrix(m, ratioWild1_5)
			numWildRatio1_5 += countWildSymbol(matrix1_5)
			checkWildPositition(matrix1_5)
			assert.Equal(t, true, checkValidSymbol(matrix1_5))

			matrix2_0 := engine.SpinMatrix(m, ratioWild2_0)
			numWildRatio2_0 += countWildSymbol(matrix2_0)
			checkWildPositition(matrix2_0)
			assert.Equal(t, true, checkValidSymbol(matrix2_0))

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
				assert.LessOrEqual(t, 3, len(paylines))
				// payline := paylines[0]
				// assert.Equal(t, int32(4), payline.NumOccur)
				// assert.Equal(t, api.SiXiangSymbol_SI_XIANG_SYMBOL_J, payline.GetSymbol())
				// assert.Equal(t, int(pair.Key), int(payline.GetId()))
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
				assert.LessOrEqual(t, 1, len(paylines))
				// payline := paylines[0]
				// assert.Equal(t, int32(5), payline.NumOccur)
				// assert.Equal(t, api.SiXiangSymbol_SI_XIANG_SYMBOL_J, payline.GetSymbol())
				// assert.Equal(t, int(pair.Key), int(payline.GetId()))
			}
			// 5 symbol xxxxx
			{
				matrix := newEmptyMatrix()
				matrix.List[ids[0]] = api.SiXiangSymbol_SI_XIANG_SYMBOL_J
				matrix.List[ids[1]] = api.SiXiangSymbol_SI_XIANG_SYMBOL_WILD
				matrix.List[ids[2]] = api.SiXiangSymbol_SI_XIANG_SYMBOL_WILD
				matrix.List[ids[3]] = api.SiXiangSymbol_SI_XIANG_SYMBOL_A
				matrix.List[ids[4]] = api.SiXiangSymbol_SI_XIANG_SYMBOL_Q
				paylines := engine.Paylines(matrix)
				assert.NotNil(t, paylines)
				assert.LessOrEqual(t, 1, len(paylines))
				// payline := paylines[0]
				// assert.Equal(t, int32(5), payline.NumOccur)
				// assert.Equal(t, api.SiXiangSymbol_SI_XIANG_SYMBOL_A, payline.GetSymbol())
				// assert.Equal(t, int(pair.Key), int(payline.GetId()))
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

func Test_normal_Paylines_2(t *testing.T) {
	name := "Test_normal_Paylines_2"
	e := NewNormal(nil)
	engine := e.(*normal)
	s := entity.NewSlotsMathState(nil)
	matrixSpecial := entity.NewJuicyMatrix()
	s.MatrixSpecial = &matrixSpecial
	s.MatrixSpecial.List[0] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_MANGOSTEEN
	s.MatrixSpecial.List[1] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STRAWBERRY
	s.MatrixSpecial.List[2] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_WATERMELON
	s.MatrixSpecial.List[3] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STONE_DIAMOND
	s.MatrixSpecial.List[4] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_WATERMELON

	s.MatrixSpecial.List[5] = api.SiXiangSymbol_SI_XIANG_SYMBOL_K
	s.MatrixSpecial.List[6] = api.SiXiangSymbol_SI_XIANG_SYMBOL_Q
	s.MatrixSpecial.List[7] = api.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER
	s.MatrixSpecial.List[8] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STRAWBERRY
	s.MatrixSpecial.List[9] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STONE_GREEN

	s.MatrixSpecial.List[10] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_MAJOR
	s.MatrixSpecial.List[11] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STONE_GREEN
	s.MatrixSpecial.List[12] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_WATERMELON
	s.MatrixSpecial.List[13] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STONE_GREEN
	s.MatrixSpecial.List[14] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_WATERMELON

	t.Run(name, func(t *testing.T) {
		paylines := engine.Paylines(*s.MatrixSpecial)
		assert.NotNil(t, paylines)
		assert.Equal(t, 3, len(paylines))
		assert.Equal(t, paylines[0].Symbol, api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STONE_GREEN)
		assert.Equal(t, paylines[0].Indices, []int32{13, 9})

	})
}

func Test_normal_Paylines_điều_kiện_ăn_1_line_bị_sai(t *testing.T) {
	name := "Test_normal_Paylines_điều_kiện_ăn_1_line_bị_sai"
	e := NewNormal(nil)
	engine := e.(*normal)
	s := entity.NewSlotsMathState(nil)
	matrixSpecial := entity.NewJuicyMatrix()
	s.MatrixSpecial = &matrixSpecial
	s.MatrixSpecial.List[0] = api.SiXiangSymbol_SI_XIANG_SYMBOL_K
	s.MatrixSpecial.List[1] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_WATERMELON
	s.MatrixSpecial.List[2] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_WATERMELON
	s.MatrixSpecial.List[3] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STONE_VIOLET
	s.MatrixSpecial.List[4] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STONE_DIAMOND

	s.MatrixSpecial.List[5] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STRAWBERRY
	s.MatrixSpecial.List[6] = api.SiXiangSymbol_SI_XIANG_SYMBOL_A
	s.MatrixSpecial.List[7] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STONE_DIAMOND
	s.MatrixSpecial.List[8] = api.SiXiangSymbol_SI_XIANG_SYMBOL_K
	s.MatrixSpecial.List[9] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_MANGOSTEEN

	s.MatrixSpecial.List[10] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_WATERMELON
	s.MatrixSpecial.List[11] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_WATERMELON
	s.MatrixSpecial.List[12] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STONE_GREEN
	s.MatrixSpecial.List[13] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_WATERMELON
	s.MatrixSpecial.List[14] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STONE_DIAMOND

	t.Run(name, func(t *testing.T) {
		paylines := engine.Paylines(*s.MatrixSpecial)
		assert.NotNil(t, paylines)
		assert.Equal(t, 0, len(paylines))

	})
}

func Test_normal_Paylines_tính_sai_tiền_ăn_line(t *testing.T) {
	name := "Test_normal_Paylines_tính_sai_tiền_ăn_line"
	e := NewNormal(nil)
	engine := e.(*normal)
	s := entity.NewSlotsMathState(nil)
	s.Bet().Chips = 100
	matrixSpecial := entity.NewJuicyMatrix()
	s.MatrixSpecial = &matrixSpecial
	e.Process(s)
	s.MatrixSpecial.List[0] = api.SiXiangSymbol_SI_XIANG_SYMBOL_K
	s.MatrixSpecial.List[1] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_WATERMELON
	s.MatrixSpecial.List[2] = api.SiXiangSymbol_SI_XIANG_SYMBOL_J
	s.MatrixSpecial.List[3] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STONE_VIOLET
	s.MatrixSpecial.List[4] = api.SiXiangSymbol_SI_XIANG_SYMBOL_K

	s.MatrixSpecial.List[5] = api.SiXiangSymbol_SI_XIANG_SYMBOL_DIAMOND
	s.MatrixSpecial.List[6] = api.SiXiangSymbol_SI_XIANG_SYMBOL_Q
	s.MatrixSpecial.List[7] = api.SiXiangSymbol_SI_XIANG_SYMBOL_K
	s.MatrixSpecial.List[8] = api.SiXiangSymbol_SI_XIANG_SYMBOL_J
	s.MatrixSpecial.List[9] = api.SiXiangSymbol_SI_XIANG_SYMBOL_Q

	s.MatrixSpecial.List[10] = api.SiXiangSymbol_SI_XIANG_SYMBOL_J
	s.MatrixSpecial.List[11] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_MANGOSTEEN
	s.MatrixSpecial.List[12] = api.SiXiangSymbol_SI_XIANG_SYMBOL_Q
	s.MatrixSpecial.List[13] = api.SiXiangSymbol_SI_XIANG_SYMBOL_Q
	s.MatrixSpecial.List[14] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_MANGOSTEEN

	s.SetMatrix(*s.MatrixSpecial)
	s.SetWildMatrix(engine.WildMatrix(s.Matrix))
	s.SetPaylines(engine.Paylines(s.WildMatrix))

	t.Run(name, func(t *testing.T) {
		result, err := e.Finish(s)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		slotDesk := result.(*api.SlotDesk)
		assert.Equal(t, int64(125), slotDesk.GameReward.TotalChipsWinByGame)

	})
}

func Test_normal_3_scatter_seq(t *testing.T) {
	name := "Test_normal_3_scatter_seq"
	e := NewNormal(nil)
	// engine := e.(*normal)
	s := entity.NewSlotsMathState(nil)
	s.Bet().Chips = 100
	e.NewGame(s)
	e.Process(s)
	t.Run(name, func(t *testing.T) {
		s.Matrix.List[0] = api.SiXiangSymbol_SI_XIANG_SYMBOL_K
		s.Matrix.List[1] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_WATERMELON
		s.Matrix.List[2] = api.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER
		s.Matrix.List[3] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STONE_VIOLET
		s.Matrix.List[4] = api.SiXiangSymbol_SI_XIANG_SYMBOL_K

		s.Matrix.List[5] = api.SiXiangSymbol_SI_XIANG_SYMBOL_K
		s.Matrix.List[6] = api.SiXiangSymbol_SI_XIANG_SYMBOL_Q
		s.Matrix.List[7] = api.SiXiangSymbol_SI_XIANG_SYMBOL_K
		s.Matrix.List[8] = api.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER
		s.Matrix.List[9] = api.SiXiangSymbol_SI_XIANG_SYMBOL_Q

		s.Matrix.List[10] = api.SiXiangSymbol_SI_XIANG_SYMBOL_J
		s.Matrix.List[11] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_MANGOSTEEN
		s.Matrix.List[12] = api.SiXiangSymbol_SI_XIANG_SYMBOL_Q
		s.Matrix.List[13] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_MANGOSTEEN
		s.Matrix.List[14] = api.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER
		result, err := e.Finish(s)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		slotDesk := result.(*api.SlotDesk)
		assert.Equal(t, api.SiXiangGame_SI_XIANG_GAME_JUICE_FRUIT_BASKET, slotDesk.NextSixiangGame)
	})
}

func Test_normal_3_scatter_seq_2(t *testing.T) {
	name := "Test_normal_3_scatter_seq"
	e := NewNormal(nil)
	// engine := e.(*normal)
	s := entity.NewSlotsMathState(nil)
	s.Bet().Chips = 100
	e.NewGame(s)
	e.Process(s)
	t.Run(name, func(t *testing.T) {
		s.Matrix.List[0] = api.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER
		s.Matrix.List[1] = api.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER
		s.Matrix.List[2] = api.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER
		s.Matrix.List[3] = api.SiXiangSymbol_SI_XIANG_SYMBOL_K
		s.Matrix.List[4] = api.SiXiangSymbol_SI_XIANG_SYMBOL_Q

		s.Matrix.List[5] = api.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER
		s.Matrix.List[6] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STONE_GREEN
		s.Matrix.List[7] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_MANGOSTEEN
		s.Matrix.List[8] = api.SiXiangSymbol_SI_XIANG_SYMBOL_K
		s.Matrix.List[9] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_MANGOSTEEN

		s.Matrix.List[10] = api.SiXiangSymbol_SI_XIANG_SYMBOL_Q
		s.Matrix.List[11] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STRAWBERRY
		s.Matrix.List[12] = api.SiXiangSymbol_SI_XIANG_SYMBOL_Q
		s.Matrix.List[13] = api.SiXiangSymbol_SI_XIANG_SYMBOL_J
		s.Matrix.List[14] = api.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER
		result, err := e.Finish(s)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		slotDesk := result.(*api.SlotDesk)
		assert.Equal(t, api.SiXiangGame_SI_XIANG_GAME_JUICE_FRUIT_BASKET, slotDesk.NextSixiangGame)
	})
}

func Test_normal_3_scatter_not_seq(t *testing.T) {
	name := "Test_normal_3_scatter_not_seq"
	e := NewNormal(nil)
	// engine := e.(*normal)
	s := entity.NewSlotsMathState(nil)
	s.Bet().Chips = 100
	e.NewGame(s)
	e.Process(s)
	t.Run(name, func(t *testing.T) {
		s.Matrix.List[0] = api.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER
		s.Matrix.List[1] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_WATERMELON
		s.Matrix.List[2] = api.SiXiangSymbol_SI_XIANG_SYMBOL_K
		s.Matrix.List[3] = api.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER
		s.Matrix.List[4] = api.SiXiangSymbol_SI_XIANG_SYMBOL_K

		s.Matrix.List[5] = api.SiXiangSymbol_SI_XIANG_SYMBOL_J
		s.Matrix.List[6] = api.SiXiangSymbol_SI_XIANG_SYMBOL_Q
		s.Matrix.List[7] = api.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER
		s.Matrix.List[8] = api.SiXiangSymbol_SI_XIANG_SYMBOL_J
		s.Matrix.List[9] = api.SiXiangSymbol_SI_XIANG_SYMBOL_Q

		s.Matrix.List[10] = api.SiXiangSymbol_SI_XIANG_SYMBOL_J
		s.Matrix.List[11] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_MANGOSTEEN
		s.Matrix.List[12] = api.SiXiangSymbol_SI_XIANG_SYMBOL_Q
		s.Matrix.List[13] = api.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER
		s.Matrix.List[14] = api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_MANGOSTEEN
		result, err := e.Finish(s)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		slotDesk := result.(*api.SlotDesk)
		assert.Equal(t, api.SiXiangGame_SI_XIANG_GAME_NORMAL, slotDesk.NextSixiangGame)
	})
}

func Test_normal_Paylines_Panic(t *testing.T) {
	name := "Test_normal_Paylines_Panic"
	e := NewNormal(nil)
	engine := e.(*normal)
	s := entity.NewSlotsMathState(nil)
	s.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_NORMAL
	engine.NewGame(s)
	t.Run(name, func(t *testing.T) {
		for i := 0; i < 10000; i++ {
			s.NumSpinLeft = 2
			engine.Process(s)
			engine.Finish(s)
		}
	})
}

func Test_normal_Paylines_Panic_2(t *testing.T) {
	name := "Test_normal_Paylines_Panicc"
	e := NewNormal(nil)
	engine := e.(*normal)
	s := entity.NewSlotsMathState(nil)
	s.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_NORMAL
	engine.NewGame(s)

	arr := []int{4356, 4354, 4353, 2, 4357, 4353, 2, 65535, 4357, 4358, 4357, 4359, 2, 4357, 4356}
	for idx, val := range arr {
		s.Matrix.List[idx] = api.SiXiangSymbol(val)
	}
	t.Run(name, func(t *testing.T) {
		e.Process(s)
		e.Finish(s)
	})
}

func Test_normal_GetNextSiXiangGame(t *testing.T) {
	name := "Test_normal_GetNextSiXiangGame"
	e := NewNormal(nil)
	engine := e.(*normal)
	s := entity.NewSlotsMathState(nil)
	t.Run(name, func(t *testing.T) {
		s.GameConfig.NumScatterSeq = 3
		s.NumFruitBasket = 0
		nextGame := engine.GetNextSiXiangGame(s)
		assert.Equal(t, pb.SiXiangGame_SI_XIANG_GAME_JUICE_FRUIT_BASKET, nextGame)

		s.GameConfig.NumScatterSeq = 0
		s.Matrix = entity.NewJuicyMatrix()
		nextGame = engine.GetNextSiXiangGame(s)
		assert.Equal(t, pb.SiXiangGame_SI_XIANG_GAME_NORMAL, nextGame)

		s.NumFruitBasket = 6
		nextGame = engine.GetNextSiXiangGame(s)
		assert.Equal(t, pb.SiXiangGame_SI_XIANG_GAME_JUICE_FRUIT_RAIN, nextGame)

		s.GameConfig.NumScatterSeq = 3
		nextGame = engine.GetNextSiXiangGame(s)
		assert.Equal(t, pb.SiXiangGame_SI_XIANG_GAME_JUICE_FRUIT_BASKET, nextGame)
	})
}

func Test_normal_Only_Payline_Finish(t *testing.T) {
	name := "Test_normal_only_payline_Finish"
	e := NewNormal(nil)
	// s := entity.NewSlotsMathState(nil)
	// e.NewGame(s)
	// e.Process(s)
	t.Run(name, func(t *testing.T) {
		// test payline line win with num scatter sequence
		listNumScatterSeq := []int{1, 2, 3, 4, 5}
		for _, numScatterSeq := range listNumScatterSeq {
			s := entity.NewSlotsMathState(nil)
			chipMcb := 100
			s.SetBetInfo(&pb.InfoBet{
				Chips: int64(chipMcb),
			})
			s.Matrix = entity.NewJuicyMatrix()
			s.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_NORMAL
			s.Matrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
				s.Matrix.List[idx] = api.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED
			})
			for i := 0; i < numScatterSeq; i++ {
				s.Matrix.List[i] = api.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER
			}
			// paylineSymbols := s.Matrix.ListFromIndexs(ids)
			lineWin := 100
			payline := &pb.Payline{
				Symbol:   api.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED,
				NumOccur: 0,
				Id:       1,
				Rate:     float64(lineWin),
			}
			s.SetPaylines([]*pb.Payline{payline})
			result, err := e.Finish(s)
			assert.NoError(t, err)
			assert.NotNil(t, result)
			slotDesk := result.(*api.SlotDesk)
			assert.NotNil(t, slotDesk)
			assert.Equal(t, int(100), int(slotDesk.ChipsMcb))
			if numScatterSeq < 3 {
				assert.Equal(t, float32(1), float32(s.GameConfig.RatioBasket))
				// assert.Equal(t, int(lineWin), int(s.LineWinByGame[s.CurrentSiXiangGame]))
				assert.Equal(t, int(lineWin), int(s.ChipStat.LineWin(s.CurrentSiXiangGame)))
				assert.Equal(t, api.SiXiangGame_SI_XIANG_GAME_NORMAL, slotDesk.NextSixiangGame)
				assert.Equal(t, api.SiXiangGame_SI_XIANG_GAME_NORMAL, s.NextSiXiangGame)
				assert.Equal(t, int(lineWin*int(slotDesk.ChipsMcb)/100), int(slotDesk.GameReward.ChipsWin))

			} else {
				assert.Equal(t, float32(numScatterSeq), float32(s.GameConfig.RatioBasket))
				// assert.Equal(t, int(lineWin*numScatterSeq), int(s.LineWinByGame[s.CurrentSiXiangGame]))
				assert.Equal(t, int(lineWin*numScatterSeq), int(s.ChipStat.LineWin(s.CurrentSiXiangGame)))
				assert.Equal(t, int(lineWin*numScatterSeq), int(slotDesk.GameReward.ChipsWin))
				assert.Equal(t, api.SiXiangGame_SI_XIANG_GAME_JUICE_FRUIT_BASKET, slotDesk.NextSixiangGame)
				assert.Equal(t, api.SiXiangGame_SI_XIANG_GAME_JUICE_FRUIT_BASKET, s.NextSiXiangGame)
				assert.Equal(t, int(lineWin*int(float64(slotDesk.ChipsMcb)/100.0*float64(s.GameConfig.RatioBasket))), int(slotDesk.GameReward.ChipsWin))
			}
			assert.Equal(t, 0, s.NumFruitBasket)
		}

	})
}

func Test_normal_Only_fruitbasket_Finish(t *testing.T) {
	name := "Test_normal_only_payline_Finish"
	e := NewNormal(nil)
	engine := e.(*normal)
	// s := entity.NewSlotsMathState(nil)
	// e.NewGame(s)
	// e.Process(s)
	t.Run(name, func(t *testing.T) {
		// test payline line win with num scatter sequence
		listNumScatterSeq := []int{1, 2, 3, 4, 5}
		for _, numScatterSeq := range listNumScatterSeq {
			listNumFruitbasket := []int{0, 1, 2, 3, 4, 5, 6, 7}
			// ids := pair.Value
			for _, numFruitbasket := range listNumFruitbasket {
				chipMcb := 100
				s := entity.NewSlotsMathState(nil)
				s.SetBetInfo(&pb.InfoBet{
					Chips: int64(chipMcb),
				})
				s.Matrix = entity.NewJuicyMatrix()
				s.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_NORMAL
				s.Matrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
					s.Matrix.List[idx] = api.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED
				})
				for i := 0; i < numScatterSeq; i++ {
					s.Matrix.List[i+2*entity.ColsJuicyMatrix] = api.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER
				}
				// t.Logf("num numScatterSeq %d num fruitBasket %d ratioFruitBasket %d",
				// 	s.GameConfig.NumScatterSeq, float32(s.NumFruitBasket), float32(s.GameConfig.RatioBasket))
				listFruitbasket := []pb.SiXiangSymbol{
					api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_MAJOR,
					api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_MINOR,
					api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_MINI,
				}
				rationFruitBasket := int(engine.transformNumScaterSeqToRationFruitBasket(numScatterSeq))
				s.GameConfig.NumScatterSeq = int64(numScatterSeq)
				lineWin := 0
				for i := 0; i < numFruitbasket; i++ {
					symbol := entity.ShuffleSlice(listFruitbasket)[0]
					s.Matrix.List[i] = symbol
					val := entity.JuicyBasketSymbol[symbol]
					lineWin += int(val.Value.Min) * rationFruitBasket
				}

				s.SetPaylines([]*pb.Payline{})
				result, err := e.Finish(s)
				assert.Equal(t, float32(rationFruitBasket), float32(s.GameConfig.RatioBasket))
				assert.NoError(t, err)
				assert.NotNil(t, result)
				slotDesk := result.(*api.SlotDesk)
				assert.NotNil(t, slotDesk)
				assert.Equal(t, int(chipMcb), int(slotDesk.ChipsMcb))
				nextGame := api.SiXiangGame_SI_XIANG_GAME_NORMAL
				chipWin := int(lineWin * int(slotDesk.ChipsMcb) / 100)
				if numScatterSeq < 3 {
					// assert.Equal(t, int(lineWin), int(s.LineWinByGame[s.CurrentSiXiangGame]))
					assert.Equal(t, int(lineWin), int(s.ChipStat.LineWin(s.CurrentSiXiangGame)))
					assert.Equal(t, chipWin, int(slotDesk.GameReward.ChipsWin))
					assert.Equal(t, chipWin, int(s.ChipStat.ChipWin(s.CurrentSiXiangGame)))
					if numFruitbasket >= 6 {
						nextGame = api.SiXiangGame_SI_XIANG_GAME_JUICE_FRUIT_RAIN
					}
				} else {
					assert.Equal(t, int(lineWin), int(s.ChipStat.LineWin(s.CurrentSiXiangGame)))
					assert.Equal(t, int(chipWin), int(slotDesk.GameReward.ChipsWin))
					assert.Equal(t, chipWin, int(slotDesk.GameReward.ChipsWin))
					assert.Equal(t, chipWin, int(s.ChipStat.ChipWin(s.CurrentSiXiangGame)))
					nextGame = api.SiXiangGame_SI_XIANG_GAME_JUICE_FRUIT_BASKET
				}
				// t.Logf("num numScatterSeq %d num fruitBasket %d ratioFruitBasket %d",
				// 	s.GameConfig.NumScatterSeq, float32(s.NumFruitBasket), float32(s.GameConfig.RatioBasket))
				assert.Equal(t, nextGame, slotDesk.NextSixiangGame)
				assert.Equal(t, nextGame, s.NextSiXiangGame)
				assert.Equal(t, numFruitbasket, s.NumFruitBasket)
			}
		}
	})
}
