package sixiang

import (
	"fmt"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

var _ lib.Engine = &normalEngine{}

type normalEngine struct {
	randomIntFn func(min, max int) int
}

func NewNormalEngine(randFn func(min, max int) int) lib.Engine {
	engine := normalEngine{
		randomIntFn: randFn,
	}
	if engine.randomIntFn == nil {
		engine.randomIntFn = entity.RandomInt
	}
	return &engine
}

func init() {

}
func AllowScatter(col int) bool {
	return entity.RowsAllowScatter[col]
}

func (e *normalEngine) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	s.SpinList = make([]*pb.SpinSymbol, 0)
	s.SpinSymbols = make([]*pb.SpinSymbol, 0)
	if len(s.Matrix.List) > 0 {
		spreadMatrix := e.SpreadWildInMatrix(s.Matrix)
		s.SetWildMatrix(spreadMatrix)
		return s, nil
	}
	matrix := entity.NewSiXiangMatrixNormal()
	matrix = e.SpinMatrix(matrix)
	s.SetMatrix(matrix)
	s.NumSpinLeft = -1
	s.ChipStat.Reset(s.CurrentSiXiangGame)
	s.SpinList = make([]*pb.SpinSymbol, 0)
	return s, nil
}

func (e *normalEngine) Random(min, max int) int {
	return entity.RandomInt(min, max)
}

func (e *normalEngine) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	s.IsSpinChange = true
	matrix := e.SpinMatrix(s.Matrix)
	if s.Bet().GetReqSpecGame() != 0 {
		matrix.List[0+5] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER
		matrix.List[2+5] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER
		matrix.List[4+5] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER
	}
	s.SetMatrix(matrix)
	spreadMatrix := e.SpreadWildInMatrix(matrix)
	s.SetWildMatrix(spreadMatrix)
	// logic
	s.SetPaylines(make([]*pb.Payline, 0))
	if e.CheckJpMatrix(spreadMatrix) {
		s.WinJp = pb.WinJackpot_WIN_JACKPOT_GRAND
	} else {
		paylines := e.PaylineMatrix(spreadMatrix)
		paylinesFilter := e.FilterPayline(paylines, func(numOccur int) bool {
			return numOccur >= 3
		})
		s.SetPaylines(paylinesFilter)
	}
	chipsMcb := s.Bet().Chips

	for _, payline := range s.Paylines() {
		payline.Rate = e.RatioPayline(payline)
		payline.Chips = int64(payline.Rate * float64(chipsMcb))
	}
	return s, nil
}

func (e *normalEngine) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	slotDesk := &pb.SlotDesk{}
	if !s.IsSpinChange {
		return s.LastResult, entity.ErrorSpinNotChange
	}
	slotDesk.GameReward = &pb.GameReward{}
	s.IsSpinChange = false
	// set matrix spin
	{
		sm := s.Matrix
		slotDesk.Matrix = sm.ToPbSlotMatrix()
	}
	slotDesk.ChipsMcb = s.Bet().GetChips()
	// set matrix spread matrix wild symbol
	{
		sm := s.WildMatrix
		slotDesk.SpreadMatrix = sm.ToPbSlotMatrix()
	}
	// add payline result
	if s.WinJp == pb.WinJackpot_WIN_JACKPOT_GRAND {
		slotDesk.GameReward.ChipsWin = int64(s.WinJp) * s.Bet().Chips
		slotDesk.BigWin = pb.BigWin_BIG_WIN_MEGA
		slotDesk.WinJp = s.WinJp
	} else {
		totalRate := float64(0)
		slotDesk.Paylines = s.Paylines()
		for _, payline := range slotDesk.Paylines {
			slotDesk.GameReward.ChipsWin += payline.GetChips()
			totalRate += payline.Rate
		}
		slotDesk.BigWin = e.TotalRateToTypeBigWin(totalRate)
	}
	// check if win bonus game
	{
		nextSiXiangGame := e.GetNextSiXiangGame(s)
		s.NextSiXiangGame = nextSiXiangGame

	}
	slotDesk.CurrentSixiangGame = s.CurrentSiXiangGame
	slotDesk.NextSixiangGame = s.NextSiXiangGame
	slotDesk.IsFinishGame = true
	slotDesk.NumSpinLeft = int64(s.NumSpinLeft)
	slotDesk.GameReward.TotalChipsWinByGame = slotDesk.GameReward.ChipsWin
	s.LastResult = slotDesk
	return slotDesk, nil
}

func (e *normalEngine) Loop(s interface{}) (interface{}, error) {
	return s, nil
}

func (e *normalEngine) SpinMatrix(matrix entity.SlotMatrix) entity.SlotMatrix {
	// matrix := matchState.GetMatrix()
	mapColExistScatter := make(map[int]bool)
	spinMatrix := entity.NewSlotMatrix(matrix.Rows, matrix.Cols)
	spinMatrix.List = make([]pb.SiXiangSymbol, spinMatrix.Size)
	matrix.ForEeach(func(idx, row, col int, _ pb.SiXiangSymbol) {
		for {
			// numRandom := e.Random(0, len(entity.ListSymbol))
			// symbol := entity.ListSymbol[numRandom]
			randSymbol := entity.ListSymbolSpinInSixiangNormal[e.Random(0, len(entity.ListSymbolSpinInSixiangNormal))]
			if randSymbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER {
				// Scatter only allow appear in list RellsAllowScatter
				if !AllowScatter(col) {
					continue
				}
				// check if symbol scatter already exist in this row
				if mapColExistScatter[col] {
					continue
				}
				mapColExistScatter[col] = true
			}
			spinMatrix.List[idx] = randSymbol
			break
		}
	})
	return spinMatrix
}

func (e *normalEngine) SpreadWildInMatrix(matrix entity.SlotMatrix) entity.SlotMatrix {
	// matrix := matchState.GetMatrix()
	spreadMatrix := entity.SlotMatrix{
		List: make([]pb.SiXiangSymbol, len(matrix.List)),
		Cols: matrix.Cols,
		Rows: matrix.Rows,
	}

	matrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		if symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_WILD {
			cols := spreadMatrix.Cols
			spreadMatrix.List[idx] = symbol
			for row := 0; row < matrix.Rows; row++ {
				id := row*cols + col
				spreadMatrix.List[id] = symbol
			}
			return
		}
		if spreadMatrix.List[idx] != pb.SiXiangSymbol_SI_XIANG_SYMBOL_WILD {
			spreadMatrix.List[idx] = symbol
		}
	})
	return spreadMatrix
}

func (e *normalEngine) CheckJpMatrix(matrix entity.SlotMatrix) bool {
	for _, symbol := range matrix.List {
		if symbol != pb.SiXiangSymbol_SI_XIANG_SYMBOL_WILD {
			return false
		}
	}
	return true
}

// return payline, and check jackpot if win
func (e *normalEngine) PaylineMatrix(matrix entity.SlotMatrix) []*pb.Payline {
	paylines := make([]*pb.Payline, 0)
	// idx := 0
	for pair := entity.MapPaylineIdx.Oldest(); pair != nil; pair = pair.Next() {
		payline := &pb.Payline{
			// Id: int32(idx),
		}
		payline.Id = int32(pair.Key)
		// idx++
		symbols := matrix.ListFromIndexs(pair.Value)
		payline.Indices = make([]int32, 0)
		// for _, val := range entity.ListSymbol {
		// 	numOccur := 0
		// 	if val == pb.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER {
		// 		continue
		// 	}
		// 	if val == pb.SiXiangSymbol_SI_XIANG_SYMBOL_WILD {
		// 		continue
		// 	}
		// 	for idx, symbol := range symbols {
		// 		if (symbol & val) > 0 {
		// 			numOccur++
		// 			payline.Indices = append(payline.Indices, int32(pair.Value[idx]))
		// 			continue
		// 		}
		// 		break
		// 	}
		// 	if numOccur > int(payline.NumOccur) {
		// 		payline.NumOccur = int32(numOccur)
		// 		payline.Symbol = val
		// 	}
		// 	if numOccur >= 3 {
		// 		break
		// 	}
		// }
		compareSym := pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED
		for _, symbol := range symbols {
			if symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER {
				continue
			}
			if symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_WILD {
				continue
			}
			compareSym = symbol
			break
		}
		if compareSym == pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED {
			continue
		}
		payline.NumOccur = 0
		for idx, symbol := range symbols {
			if symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_WILD || symbol == compareSym {
				payline.NumOccur++
				payline.Indices = append(payline.Indices, int32(pair.Value[idx]))
				continue
			}
			break
		}
		if payline.NumOccur == 0 {
			continue
		}
		payline.Symbol = compareSym
		paylines = append(paylines, payline)
	}
	return paylines
}

func (e *normalEngine) FilterPayline(paylines []*pb.Payline, fn func(numOccur int) bool) []*pb.Payline {
	list := make([]*pb.Payline, 0)
	for _, payline := range paylines {
		if fn(int(payline.GetNumOccur())) {
			list = append(list, payline)
		}
	}
	return list
}

func (e *normalEngine) RatioPayline(payline *pb.Payline) float64 {
	return entity.RatioPaylineMap[payline.Symbol][payline.NumOccur]
}

func (e *normalEngine) TotalRateToTypeBigWin(totalRate float64) pb.BigWin {
	bigWin := pb.BigWin_BIG_WIN_UNSPECIFIED
	if int(totalRate) >= int(pb.BigWin_BIG_WIN_MEGA.Number()) {
		bigWin = pb.BigWin_BIG_WIN_MEGA
	} else if int(totalRate) >= int(pb.BigWin_BIG_WIN_HUGE.Number()) {
		bigWin = pb.BigWin_BIG_WIN_HUGE
	} else if int(totalRate) >= int(pb.BigWin_BIG_WIN_BIG.Number()) {
		bigWin = pb.BigWin_BIG_WIN_BIG
	} else if int(totalRate) >= int(pb.BigWin_BIG_WIN_NICE.Number()) {
		bigWin = pb.BigWin_BIG_WIN_NICE
	}
	return bigWin
}

func (e *normalEngine) GetNextSiXiangGame(s *entity.SlotsMatchState) pb.SiXiangGame {
	matrix := s.Matrix
	numScatter := 0
	matrix.ForEachCol(func(col int, symbols []pb.SiXiangSymbol) {
		if !entity.RowsAllowScatter[col] {
			return
		}
		for _, symbol := range symbols {
			if symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER {
				numScatter++
			}
		}
	})
	if numScatter >= 3 {
		return pb.SiXiangGame_SI_XIANG_GAME_BONUS
	}
	return pb.SiXiangGame_SI_XIANG_GAME_NORMAL
}

func (e *normalEngine) PrintMatrix(matrix entity.SlotMatrix) {
	// matrix := matchState.GetMatrix()
	matrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		if idx != 0 && col == 0 {
			fmt.Println("")
		}
		fmt.Printf("%8d", symbol.Number())
	})
	fmt.Println("")
}
