package sixiang

import (
	"fmt"
	"time"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
	"github.com/ciaolink-game-platform/cgp-common/utilities"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ lib.Engine = &slotsEngine{}

type slotsEngine struct {
	engines map[pb.SiXiangGame]lib.Engine
}

func newEngineByGame(game pb.SiXiangGame) lib.Engine {
	switch game {
	case pb.SiXiangGame_SI_XIANG_GAME_NORMAL:
		return NewNormalEngine(nil)
	case pb.SiXiangGame_SI_XIANG_GAME_BONUS:
		return NewBonusEngine(nil)
	case pb.SiXiangGame_SI_XIANG_GAME_DRAGON_PEARL:
		return NewDragonPearlEngine(1, nil, nil)

	case pb.SiXiangGame_SI_XIANG_GAME_LUCKDRAW:
		return NewLuckyDrawEngine(1, nil, nil)
	case pb.SiXiangGame_SI_XIANG_GAME_GOLDPICK:
		return NewGoldPickEngine(1, nil, nil)
	case pb.SiXiangGame_SI_XIANG_GAME_RAPIDPAY:
		return NewRapidPayEngine(1, nil, nil)
	case pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS:
		return NewSixiangBonusEngine()
	case pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_DRAGON_PEARL,
		pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_LUCKDRAW,
		pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_GOLDPICK,
		pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_RAPIDPAY:
		return NewSixiangBonusInGameEngine(4)
	}
	return NewNormalEngine(nil)
}

func NewEngine() lib.Engine {
	slotEngine := slotsEngine{}

	slotEngine.engines = make(map[pb.SiXiangGame]lib.Engine)
	// i := 1
	games := []pb.SiXiangGame{
		pb.SiXiangGame_SI_XIANG_GAME_NORMAL,
		pb.SiXiangGame_SI_XIANG_GAME_BONUS,
		pb.SiXiangGame_SI_XIANG_GAME_DRAGON_PEARL,
		pb.SiXiangGame_SI_XIANG_GAME_LUCKDRAW,
		pb.SiXiangGame_SI_XIANG_GAME_GOLDPICK,
		pb.SiXiangGame_SI_XIANG_GAME_RAPIDPAY,
		// bonus
		pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS,
		pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_DRAGON_PEARL,
		pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_LUCKDRAW,
		pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_GOLDPICK,
		pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_RAPIDPAY,
	}
	for _, game := range games {
		// game := pb.SiXiangGame(i)
		slotEngine.engines[game] = newEngineByGame(game)
	}
	return &slotEngine
}

func (e *slotsEngine) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	engine, ok := e.engines[s.CurrentSiXiangGame]
	if !ok {
		return nil, status.Error(codes.Unimplemented, "not implement new game "+s.CurrentSiXiangGame.String())
	}
	s.LastResult = nil
	s.LastSpinTime = time.Now()
	engine.NewGame(s)
	return nil, nil
}

func (e *slotsEngine) Random(min, max int) int {
	return utilities.RandomNumber(min, max)
}

func (e *slotsEngine) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	engine, ok := e.engines[s.CurrentSiXiangGame]
	if !ok {
		return nil, status.Error(codes.Unimplemented, "not implement process game "+s.CurrentSiXiangGame.String())
	}
	return engine.Process(matchState)
}

func (e *slotsEngine) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	engine, ok := e.engines[s.CurrentSiXiangGame]
	if !ok {
		return nil, status.Error(codes.Unimplemented, "not implement fisnish game "+s.CurrentSiXiangGame.String())
	}
	result, err := engine.Finish(matchState)
	if err != nil {
		return result, err
	}
	slotDesk, ok := result.(*pb.SlotDesk)
	if !ok {
		return result, err
	}
	ratio := entity.PriceBuySixiangGem[s.NumGameEyePlayed()]
	chips := ratio * int(s.Bet().Chips)
	slotDesk.ChipsBuyGem = int64(chips)
	slotDesk.SixiangGems = make([]pb.SiXiangGame, 0)
	for gem := range s.GameEyePlayed() {
		slotDesk.SixiangGems = append(slotDesk.SixiangGems, gem)
	}
	slotDesk.WinJpHistory = s.WinJPHistory()
	slotDesk.BetLevels = make([]int64, 0)
	slotDesk.BetLevels = append(slotDesk.BetLevels, entity.BetLevels...)
	return slotDesk, nil
}

func (e *slotsEngine) Loop(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	engine := e.engines[s.CurrentSiXiangGame]
	return engine.Loop(s)
}

func (e *slotsEngine) Info(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	var matrix *pb.SlotMatrix
	var spreadMatrix *pb.SlotMatrix
	switch s.CurrentSiXiangGame {
	case pb.SiXiangGame_SI_XIANG_GAME_NORMAL:
		spreadMatrix = s.WildMatrix.ToPbSlotMatrix()
		matrix = spreadMatrix
	case pb.SiXiangGame_SI_XIANG_GAME_TARZAN_FREESPINX9,
		pb.SiXiangGame_SI_XIANG_GAME_JUICE_FRUIT_RAIN,
		pb.SiXiangGame_SI_XIANG_GAME_JUICE_FREE_GAME:
		matrix = s.Matrix.ToPbSlotMatrix()
		spreadMatrix = s.WildMatrix.ToPbSlotMatrix()
	case
		pb.SiXiangGame_SI_XIANG_GAME_DRAGON_PEARL,
		pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_DRAGON_PEARL,
		pb.SiXiangGame_SI_XIANG_GAME_LUCKDRAW,
		pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_LUCKDRAW,
		pb.SiXiangGame_SI_XIANG_GAME_GOLDPICK,
		pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_GOLDPICK:
		matrix = s.MatrixSpecial.ToPbSlotMatrix()
		for idx, symbol := range s.MatrixSpecial.List {
			if s.MatrixSpecial.IsFlip(idx) {
				matrix.Lists[idx] = symbol
			} else {
				matrix.Lists[idx] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED
			}
		}
	default:
		matrix = s.MatrixSpecial.ToPbSlotMatrix()
		spreadMatrix = s.MatrixSpecial.ToPbSlotMatrix()
	}
	matrix.SpinLists = s.SpinList
	slotdesk := &pb.SlotDesk{
		Matrix:             matrix,
		SpreadMatrix:       spreadMatrix,
		ChipsMcb:           s.Bet().Chips,
		CurrentSixiangGame: s.CurrentSiXiangGame,
		NextSixiangGame:    s.NextSiXiangGame,
		TsUnix:             time.Now().Unix(),
		SpinSymbols:        s.SpinSymbols,
		NumSpinLeft:        int64(s.NumSpinLeft),
		InfoBet:            s.Bet(),
		WinJpHistory:       s.WinJPHistory(),
		BetLevels:          entity.BetLevels[:],
	}
	slotdesk.ChipsBuyGem, _ = s.PriceBuySixiangGem()
	// slotdesk.LetterSymbols = make([]pb.SiXiangSymbol, 0)
	// for k := range s.LetterSymbol {
	// 	slotdesk.LetterSymbols = append(slotdesk.LetterSymbols, k)
	// }
	slotdesk.SixiangGems = make([]pb.SiXiangGame, 0)
	for gem := range s.GameEyePlayed() {
		slotdesk.SixiangGems = append(slotdesk.SixiangGems, gem)
	}
	return slotdesk, nil
}

func (e *slotsEngine) PrintMatrix(matrix entity.SlotMatrix) {
	matrix.ForEeach(func(idx, _, col int, symbol pb.SiXiangSymbol) {
		if idx != 0 && col == 0 {
			fmt.Println("")
		}
		fmt.Printf("%8d", symbol.Number())
	})
	fmt.Println("")
}
