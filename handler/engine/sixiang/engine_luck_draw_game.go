package sixiang

import (
	"fmt"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

var _ lib.Engine = &luckyDrawEngine{}

type luckyDrawEngine struct {
	randomIntFn         func(min, max int) int
	randomFloat64       func(min, max float64) float64
	ratioInSixiangBonus int
}

func NewLuckyDrawEngine(ratioInSixiangBonus int, randomIntFn func(min, max int) int, randomFloat64 func(min, max float64) float64) lib.Engine {
	engine := luckyDrawEngine{
		ratioInSixiangBonus: ratioInSixiangBonus,
	}
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
	s.ChipStat.Reset(s.CurrentSiXiangGame)
	s.SpinList = make([]*pb.SpinSymbol, 0)
	s.MatrixSpecial.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		s.SpinList = append(s.SpinList, &pb.SpinSymbol{
			Symbol:    pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED,
			Row:       int32(row),
			Col:       int32(col),
			Index:     int32(idx),
			Ratio:     0,
			WinJp:     pb.WinJackpot_WIN_JACKPOT_UNSPECIFIED,
			WinAmount: 0,
		})
	})
	return s, nil
}

func (e *luckyDrawEngine) Random(min, max int) int {
	return entity.RandomInt(min, max)
}

func (e *luckyDrawEngine) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	if s.Bet().Id < 0 {
		return nil, entity.ErrorSpinIndexRequired
	}
	if int(s.Bet().Id) >= len(s.MatrixSpecial.List) {
		return nil, entity.ErrorInfoBetInvalid
	}

	if s.MatrixSpecial.IsFlip(int(s.Bet().Id)) {
		return nil, entity.ErrorSpinIndexAleadyTaken
	}
	// check if game already collect 3 jp symbol
	if e.GetNextSiXiangGame(s) == pb.SiXiangGame_SI_XIANG_GAME_NORMAL {
		return nil, entity.ErrorInvalidRequestGame

	}
	// bet := s.GetBetInfo()
	// check if all id already flip
	{
		idsNotFlip := make([]int, 0)
		for id := range s.MatrixSpecial.List {
			if !s.MatrixSpecial.IsFlip(id) {
				idsNotFlip = append(idsNotFlip, id)
			}
		}
		if len(idsNotFlip) == 0 {
			return s, entity.ErrorSpinReachMax
		}
	}
	s.IsSpinChange = true
	// id := e.Random(0, len(idsNotFlip))
	// id:=
	// idFlip := idsNotFlip[id]
	idFlip := int(s.Bet().Id)
	// s.MatrixSpecial.TrackFlip[idFlip] = true
	spinSymbol := &pb.SpinSymbol{
		Symbol: s.MatrixSpecial.Flip(idFlip),
		Index:  int32(idFlip),
	}
	row, col := s.MatrixSpecial.RowCol(idFlip)
	spinSymbol.Row = int32(row)
	spinSymbol.Col = int32(col)
	spinSymbol.Index = int32(idFlip)
	_, spinSymbol.WinJp = entity.LuckySymbolToReward(spinSymbol.Symbol)
	s.SpinSymbols = []*pb.SpinSymbol{spinSymbol}
	s.SpinList[spinSymbol.Index] = spinSymbol
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
		ChipsMcb:   s.Bet().Chips,
	}
	if !s.IsSpinChange {
		return s.LastResult, entity.ErrorSpinNotChange
	}
	s.IsSpinChange = false
	s.MatrixSpecial.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		if s.MatrixSpecial.IsFlip(idx) {
			slotDesk.Matrix.Lists = append(slotDesk.Matrix.Lists, symbol)
		} else {
			slotDesk.Matrix.Lists = append(slotDesk.Matrix.Lists, pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED)
		}
	})
	for _, spin := range s.SpinSymbols {
		sym := spin.Symbol
		// ignore jp symbol, 3 symbol jp -> end game and + chip win of this jp
		if sym < pb.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_GOLD_1 {
			continue
		}
		ratio, chips := e.calcRewardBySymbol(spin, s.Bet().Chips)
		spin.Ratio = ratio
		spin.WinAmount = chips
		s.SpinList[spin.Index].Ratio = spin.Ratio
		s.SpinList[spin.Index].WinAmount = spin.WinAmount
		slotDesk.GameReward.ChipsWin += spin.WinAmount
		slotDesk.GameReward.RatioWin += spin.Ratio
	}
	for _, spin := range s.SpinList {
		// _, chips := e.calcRewardBySymbol(spin, s.Bet().Chips)
		slotDesk.GameReward.TotalChipsWinByGame += spin.WinAmount
	}
	s.NextSiXiangGame = e.GetNextSiXiangGame(s)
	slotDesk.NextSixiangGame = s.NextSiXiangGame
	slotDesk.CurrentSixiangGame = s.CurrentSiXiangGame
	if s.NextSiXiangGame == pb.SiXiangGame_SI_XIANG_GAME_NORMAL {
		// calc chip in special game
		slotDesk.IsFinishGame = true
		symbolWin := s.SpinSymbols[0].Symbol
		slotDesk.BigWin, slotDesk.WinJp = entity.LuckySymbolToReward(symbolWin)
		chipsJpWin := int64(slotDesk.WinJp) * s.Bet().Chips * int64(e.ratioInSixiangBonus)
		slotDesk.GameReward.ChipsWin += chipsJpWin
		slotDesk.GameReward.TotalChipsWinByGame += chipsJpWin
	}
	slotDesk.NumSpinLeft = int64(s.NumSpinLeft)
	slotDesk.SpinSymbols = s.SpinSymbols
	slotDesk.Matrix.SpinLists = s.SpinList
	if slotDesk.IsFinishGame {
		s.AddGameEyePlayed(pb.SiXiangGame_SI_XIANG_GAME_LUCKDRAW)
	}
	s.LastResult = slotDesk
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
	return s.CurrentSiXiangGame
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

// return ratio, chip
func (e *luckyDrawEngine) calcRewardBySymbol(spin *pb.SpinSymbol, mcb int64) (float32, int64) {
	sym := spin.Symbol
	if sym < pb.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_GOLD_1 {
		return 0, 0
	}
	ratio := float64(spin.Ratio)
	if ratio <= 0 {
		ratio = e.randomFloat64(float64(entity.ListSymbolLuckyDraw[sym].Value.Min),
			float64(entity.ListSymbolLuckyDraw[sym].Value.Max))
	}
	ratio *= float64(e.ratioInSixiangBonus)
	chips := int64(ratio*100) * mcb / 100
	return float32(ratio), chips
}
