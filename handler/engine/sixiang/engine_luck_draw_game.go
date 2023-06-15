package sixiang

import (
	"fmt"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

var _ lib.Engine = &luckyDrawEngine{}

type luckyDrawEngine struct {
	randomIntFn   func(min, max int) int
	randomFloat64 func(min, max float64) float64
}

func NewLuckyDrawEngine(randomIntFn func(min, max int) int, randomFloat64 func(min, max float64) float64) lib.Engine {
	engine := luckyDrawEngine{}
	if randomIntFn != nil {
		engine.randomIntFn = randomIntFn
	} else {
		engine.randomIntFn = entity.RandomInt
	}
	if randomFloat64 != nil {
		engine.randomFloat64 = randomFloat64
	} else {
		engine.randomFloat64 = entity.RandomFloat64
	}
	return &engine
}

func (e *luckyDrawEngine) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	matrix := entity.NewMatrixLuckyDraw()
	s.MatrixSpecial = entity.ShuffleMatrix(matrix)
	s.SpinSymbols = []*pb.SpinSymbol{}
	s.NumSpinLeft = -1
	// s.ChipsWinInSpecialGame = 0
	s.ChipStat.ResetChipWin(s.CurrentSiXiangGame)
	s.ResetCollection(s.CurrentSiXiangGame, int(s.Bet().Chips))
	return s, nil
}

func (e *luckyDrawEngine) Random(min, max int) int {
	return entity.RandomInt(min, max)
}

func (e *luckyDrawEngine) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	s.IsSpinChange = true
	// bet := s.GetBetInfo()
	idsNotFlip := make([]int, 0)
	for id := range s.MatrixSpecial.List {
		if !s.MatrixSpecial.TrackFlip[id] {
			idsNotFlip = append(idsNotFlip, id)
		}
	}
	if len(idsNotFlip) == 0 {
		return s, entity.ErrorSpinReachMax
	}
	id := e.Random(0, len(idsNotFlip))
	idFlip := idsNotFlip[id]
	s.MatrixSpecial.TrackFlip[idFlip] = true
	spinSymbol := &pb.SpinSymbol{
		Symbol: s.MatrixSpecial.List[idFlip],
	}
	row, col := s.MatrixSpecial.RowCol(idFlip)
	spinSymbol.Row = int32(row)
	spinSymbol.Col = int32(col)
	s.SpinSymbols = []*pb.SpinSymbol{spinSymbol}
	return s, nil
}

func (e *luckyDrawEngine) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)

	matrix := s.MatrixSpecial
	slotDesk := &pb.SlotDesk{
		Matrix: &pb.SlotMatrix{
			Lists: make([]pb.SiXiangSymbol, 0),
			Rows:  int32(matrix.Rows),
			Cols:  int32(matrix.Cols),
		},
		GameReward: &pb.GameReward{},
	}
	if !s.IsSpinChange {
		return slotDesk, entity.ErrorSpinNotChange
	}
	s.IsSpinChange = false
	// mapUniqueSym := make(map[pb.SiXiangSymbol]pb.SiXiangSymbol)
	// for id, symbol := range matrix.List {
	// 	if s.MatrixSpecial.TrackFlip[id] {
	// 		slotDesk.Matrix.Lists = append(slotDesk.Matrix.Lists, symbol)
	// 		mapUniqueSym[symbol] = symbol
	// 	} else {
	// 		slotDesk.Matrix.Lists = append(slotDesk.Matrix.Lists, pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED)
	// 	}
	// }
	// // calc chip win
	// {
	// 	totalRatio := float64(0)
	// 	for _, symbol := range mapUniqueSym {
	// 		rangeRatio := entity.ListSymbolLuckyDraw[symbol].Value
	// 		totalRatio += e.randomFloat64(float64(rangeRatio.Min), float64(rangeRatio.Max))
	// 	}
	// 	slotDesk.GameReward.ChipsWin += int64(totalRatio * float64(s.Bet().GetChips()))
	// }
	ratioWin := float32(0)
	for _, spin := range s.SpinSymbols {
		sym := spin.Symbol
		if sym < pb.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_GOLD_1 {
			s.AddCollectionSymbol(s.CurrentSiXiangGame, int(s.Bet().Chips), sym)
			continue
		}
		spin.Ratio = float32(e.randomFloat64(float64(entity.ListSymbolLuckyDraw[sym].Value.Min), float64(entity.ListSymbolLuckyDraw[sym].Value.Max)))
		ratioWin += spin.Ratio
	}
	slotDesk.GameReward.ChipsWin = int64(float64(ratioWin) * float64(s.Bet().GetChips()))
	s.NextSiXiangGame = e.GetNextSiXiangGame(s)
	slotDesk.NextSixiangGame = s.NextSiXiangGame
	slotDesk.CurrentSixiangGame = s.CurrentSiXiangGame
	if s.NextSiXiangGame == pb.SiXiangGame_SI_XIANG_GAME_NORMAL {
		// calc chip in special game
		slotDesk.IsFinishGame = true
		symbolWin := s.SpinSymbols[0].Symbol
		slotDesk.BigWin, slotDesk.WinJp = entity.LuckySymbolToReward(symbolWin)
		slotDesk.GameReward.ChipsWin += int64(slotDesk.WinJp) * s.Bet().Chips
	}
	slotDesk.NumSpinLeft = int64(s.NumSpinLeft)
	slotDesk.SpinSymbols = s.SpinSymbols
	slotDesk.ChipsMcb = s.Bet().Chips
	s.ChipStat.AddChipWin(s.CurrentSiXiangGame, slotDesk.GameReward.ChipsWin)
	slotDesk.GameReward.TotalChipsWinByGame = s.ChipStat.TotalChipWin(s.CurrentSiXiangGame)
	slotDesk.GameReward.RatioWin = ratioWin
	slotDesk.CollectionSymbols = s.CollectionSymbolToSlice(s.CurrentSiXiangGame, int(s.Bet().Chips))
	return slotDesk, nil
}

func (e *luckyDrawEngine) Loop(s interface{}) (interface{}, error) {
	return s, nil
}

func (e *luckyDrawEngine) GetNextSiXiangGame(s *entity.SlotsMatchState) pb.SiXiangGame {
	matrix := s.MatrixSpecial
	trackFlipSameSymbol := make(map[pb.SiXiangSymbol]int)
	for id, symbol := range matrix.List {
		if s.MatrixSpecial.TrackFlip[id] {
			num := trackFlipSameSymbol[symbol]
			num++
			trackFlipSameSymbol[symbol] = num
		}
	}
	isFinishGame := false
	for k, v := range trackFlipSameSymbol {
		if int(k) < int(pb.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_GOLD_1) && v >= 3 {
			isFinishGame = true
			break
		}
	}
	if isFinishGame {
		return pb.SiXiangGame_SI_XIANG_GAME_NORMAL
	}
	return pb.SiXiangGame_SI_XIANG_GAME_LUCKDRAW
}

func (e *luckyDrawEngine) PrintMatrix(matrix entity.SlotMatrix) {
	// matrix := matchState.GetMatrix()
	matrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		if idx != 0 && col == 0 {
			fmt.Println("")
		}
		fmt.Printf("%8d", symbol.Number())
	})
	fmt.Println("")
}
