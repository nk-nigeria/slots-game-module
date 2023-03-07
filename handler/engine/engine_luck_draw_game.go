package engine

import (
	"errors"
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
		engine.randomIntFn = RandomInt
	}
	if randomFloat64 != nil {
		engine.randomFloat64 = randomFloat64
	} else {
		engine.randomFloat64 = RandomFloat64
	}
	return &engine
}

func (e *luckyDrawEngine) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	matrix := entity.NewSiXiangMatrixLuckyDraw()
	s.MatrixSpecial = ShuffleMatrix(matrix)
	s.SpinSymbols = []*pb.SpinSymbol{
		{Symbol: pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED},
	}
	// s.ChipsWinInSpecialGame = 0
	return s, nil
}

func (e *luckyDrawEngine) Random(min, max int) int {
	return RandomInt(min, max)
}

func (e *luckyDrawEngine) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	// bet := s.GetBetInfo()
	idsNotFlip := make([]int, 0)
	for id := range s.MatrixSpecial.List {
		if s.MatrixSpecial.TrackFlip[id] == false {
			idsNotFlip = append(idsNotFlip, id)
		}
	}
	if len(idsNotFlip) == 0 {
		return s, errors.New("Spin all")
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
	}
	mapUniqueSym := make(map[pb.SiXiangSymbol]pb.SiXiangSymbol)
	for id, symbol := range matrix.List {
		if s.MatrixSpecial.TrackFlip[id] {
			slotDesk.Matrix.Lists = append(slotDesk.Matrix.Lists, symbol)
			mapUniqueSym[symbol] = symbol
		} else {
			slotDesk.Matrix.Lists = append(slotDesk.Matrix.Lists, pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED)
		}
	}
	// calc chip win
	{
		totalRatio := float64(0)
		for _, symbol := range mapUniqueSym {
			rangeRatio := entity.ListSymbolLuckyDraw[symbol].Value
			totalRatio += e.randomFloat64(float64(rangeRatio.Min), float64(rangeRatio.Max))
		}
		slotDesk.ChipsWin += int64(totalRatio * float64(s.GetBetInfo().GetChips()))
	}
	s.NextSiXiangGame = e.GetNextSiXiangGame(s)
	slotDesk.NextSixiangGame = s.NextSiXiangGame
	slotDesk.CurrentSixiangGame = s.CurrentSiXiangGame
	if s.NextSiXiangGame == pb.SiXiangGame_SI_XIANG_GAME_NORMAL {
		// calc chip in special game
		slotDesk.IsFinishGame = true
		symbolWin := s.SpinSymbols[0].Symbol
		slotDesk.BigWin, slotDesk.WinJp = LuckySymbolToReward(symbolWin)
	}
	slotDesk.SpinSymbols = s.SpinSymbols
	slotDesk.ChipsMcb = s.GetBetInfo().Chips
	return slotDesk, nil
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
		return
	})
	fmt.Println("")
}
