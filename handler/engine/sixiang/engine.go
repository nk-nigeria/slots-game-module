package sixiang

import (
	"fmt"

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
		return NewNormalEngine()
	case pb.SiXiangGame_SI_XIANG_GAME_BONUS:
		return NewBonusEngine(nil)
	case pb.SiXiangGame_SI_XIANG_GAME_DRAGON_PEARL:
		return NewDragonPearlEngine(nil, nil)

	case pb.SiXiangGame_SI_XIANG_GAME_LUCKDRAW:
		return NewLuckyDrawEngine(nil, nil)
	case pb.SiXiangGame_SI_XIANG_GAME_GOLDPICK:
		return NewGoldPickEngine(nil, nil)
	case pb.SiXiangGame_SI_XIANG_GAME_RAPIDPAY:
		return NewRapidPayEngine(nil, nil)
	case pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS:
		return NewSixiangBonusEngine()
	case pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_DRAGON_PEARL,
		pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_LUCKDRAW,
		pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_GOLDPICK,
		pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_RAPIDPAY:
		return NewSixiangBonusInGameEngine(4)
	}
	return NewNormalEngine()
}

func NewEngine() lib.Engine {
	slotEngine := slotsEngine{}

	slotEngine.engines = make(map[pb.SiXiangGame]lib.Engine)
	// i := 1
	for _, i := range pb.SiXiangGame_value {
		if i == 0 {
			continue
		}
		game := pb.SiXiangGame(i)
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
	for gem := range s.GameEyePlayed() {
		slotDesk.SixiangGems = append(slotDesk.SixiangGems, gem)
	}
	var jpReward *pb.JackpotReward
	switch s.WinJp {
	case pb.WinJackpot_WIN_JACKPOT_MINOR:
		jpReward = s.WinJPHistory().Minor
	case pb.WinJackpot_WIN_JACKPOT_MAJOR:
		jpReward = s.WinJPHistory().Major
	case pb.WinJackpot_WIN_JACKPOT_MEGA:
		jpReward = s.WinJPHistory().Mega
	case pb.WinJackpot_WIN_JACKPOT_GRAND:
		jpReward = s.WinJPHistory().Grand
	}
	if jpReward != nil {
		jpReward.Count++
		jpReward.Chips = jpReward.Count * jpReward.Ratio * s.Bet().Chips
	}
	slotDesk.WinJpHistory = s.WinJPHistory()
	return slotDesk, nil
}

func (e *slotsEngine) Loop(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	engine := e.engines[s.CurrentSiXiangGame]
	return engine.Loop(s)
}

func (e *slotsEngine) PrintMatrix(matrix entity.SlotMatrix) {
	// matrix := matchState.GetMatrix()
	matrix.ForEeach(func(idx, _, col int, symbol pb.SiXiangSymbol) {
		if idx != 0 && col == 0 {
			fmt.Println("")
		}
		fmt.Printf("%8d", symbol.Number())
	})
	fmt.Println("")
}
