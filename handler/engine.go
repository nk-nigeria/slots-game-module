package handler

import (
	"fmt"
	"strconv"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgb-slots-game-module/handler/engine/sixiang"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
	"github.com/ciaolink-game-platform/cgp-common/utilities"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ lib.Engine = &slotsEngine{}

var RellsAllowScatter = map[int]bool{0: true, 2: true, 4: true}

func AllowScatter(col int) bool {
	return RellsAllowScatter[col]
}

type slotsEngine struct {
	engines map[pb.SiXiangGame]lib.Engine
}

func newEngine(game pb.SiXiangGame) lib.Engine {
	switch game {
	case pb.SiXiangGame_SI_XIANG_GAME_NORMAL:
		return sixiang.NewNormalEngine()
	case pb.SiXiangGame_SI_XIANG_GAME_BONUS:
		return sixiang.NewBonusEngine(nil)
	case pb.SiXiangGame_SI_XIANG_GAME_DRAGON_PEARL:
		return sixiang.NewDragonPearlEngine(nil, nil)

	case pb.SiXiangGame_SI_XIANG_GAME_LUCKDRAW:
		return sixiang.NewLuckyDrawEngine(nil, nil)
	case pb.SiXiangGame_SI_XIANG_GAME_GOLDPICK:
		return sixiang.NewGoldPickEngine(nil, nil)
	case pb.SiXiangGame_SI_XIANG_GAME_RAPIDPAY:
		return sixiang.NewRapidPayEngine(nil, nil)
	case pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS:
		return sixiang.NewSixiangBonusEngine()
	case pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_DRAGON_PEARL,
		pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_LUCKDRAW,
		pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_GOLDPICK,
		pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_RAPIDPAY:
		return sixiang.NewSixiangBonusInGameEngine(4)
	}
	return sixiang.NewNormalEngine()
}

func NewSlotsEngine() lib.Engine {
	slotEngine := slotsEngine{}
	slotEngine.engines = make(map[pb.SiXiangGame]lib.Engine)
	i := 1
	for {
		game := pb.SiXiangGame(i)
		if game == pb.SiXiangGame_SI_XIANG_GAME_UNSPECIFIED || i > 100 ||
			game.String() == strconv.Itoa(i) {
			break
		}
		slotEngine.engines[game] = newEngine(game)
		i++
	}
	return &slotEngine
}

func (e *slotsEngine) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SixiangMatchState)
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
	s := matchState.(*entity.SixiangMatchState)
	engine, ok := e.engines[s.CurrentSiXiangGame]
	if !ok {
		return nil, status.Error(codes.Unimplemented, "not implement process game "+s.CurrentSiXiangGame.String())
	}
	return engine.Process(matchState)
}

func (e *slotsEngine) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SixiangMatchState)
	engine, ok := e.engines[s.CurrentSiXiangGame]
	if !ok {
		return nil, status.Error(codes.Unimplemented, "not implement fisnish game "+s.CurrentSiXiangGame.String())
	}
	return engine.Finish(matchState)
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
