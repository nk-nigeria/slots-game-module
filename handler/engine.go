package handler

import (
	"fmt"
	"strconv"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgb-slots-game-module/handler/engine"
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
		return engine.NewNormalEngine()
	case pb.SiXiangGame_SI_XIANG_GAME_BONUS:
		return engine.NewBonusEngine(nil)
	case pb.SiXiangGame_SI_XIANG_GAME_DRAGON_PEARL:
		return engine.NewDragonPearlEngine(nil, nil)
	case pb.SiXiangGame_SI_XIANG_GAME_LUCKDRAW:
		return engine.NewLuckyDrawEngine(nil, nil)
	case pb.SiXiangGame_SI_XIANG_GAME_RAPIDPAY:
		return engine.NewRapidPayEngine()
	case pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS:
		return engine.NewSixiangBonusEngine()
	}
	return engine.NewNormalEngine()
}

func NewSlotsEngine() lib.Engine {
	slotEngine := slotsEngine{}
	slotEngine.engines = make(map[pb.SiXiangGame]lib.Engine)
	i := 1
	for {
		game := pb.SiXiangGame(i)
		if game == pb.SiXiangGame_SI_XIANG_GAME_UNSPECIFIED ||
			game.String() == strconv.Itoa(i) {
			break
		}
		slotEngine.engines[game] = newEngine(game)
		i++
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
	s.RatioBonus = 1
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
	return engine.Finish(matchState)
}

func (e *slotsEngine) PrintMatrix(matrix entity.SlotMatrix) {
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
