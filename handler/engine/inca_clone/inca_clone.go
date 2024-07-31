package incaclone

import (
	"errors"

	"github.com/nakamaFramework/cgb-slots-game-module/entity"
	"github.com/nakamaFramework/cgb-slots-game-module/handler/engine/inca"
	"github.com/nakamaFramework/cgp-common/define"
	"github.com/nakamaFramework/cgp-common/lib"
	pb "github.com/nakamaFramework/cgp-common/proto"
)

type engine struct {
	incaEngine       lib.Engine
	useFnSymbolAlias func(source pb.SiXiangSymbol) pb.SiXiangSymbol
}

func NewEngine(gameCode define.GameName) lib.Engine {
	e := &engine{
		incaEngine: inca.NewEngine(),
	}
	switch gameCode {
	case define.NoelGameName:
		e.useFnSymbolAlias = entity.GetSymbolNoelFromInca
	case define.FruitGameName:
		e.useFnSymbolAlias = entity.GetSymbolFruitFromInca
	default:
		panic("not implement clone game")
	}
	return e
}

// NewGame implements lib.Engine
func (e *engine) NewGame(matchState interface{}) (interface{}, error) {
	return e.incaEngine.NewGame(matchState)
}

// Process implements lib.Engine
func (e *engine) Process(matchState interface{}) (interface{}, error) {
	return e.incaEngine.Process(matchState)
}
func (e *engine) Loop(matchState interface{}) (interface{}, error) {
	return e.incaEngine.Loop(matchState)
}

// Random implements lib.Engine
func (e *engine) Random(min int, max int) int {
	// return e.Random(min, max)
	return e.incaEngine.Random(min, max)
}

// Finish implements lib.Engine.
func (e *engine) Finish(matchState interface{}) (interface{}, error) {
	result, err := e.incaEngine.Finish(matchState)
	if err != nil {
		return nil, err
	}
	slotDesk, ok := result.(*pb.SlotDesk)
	if !ok {
		return nil, errors.New("invalid result")
	}
	return e.convertSymbol(slotDesk), nil
}

// Info implements lib.Engine.
func (e *engine) Info(matchState interface{}) (interface{}, error) {
	result, err := e.incaEngine.Info(matchState)
	if err != nil {
		return nil, err
	}
	slotDesk, ok := result.(*pb.SlotDesk)
	if !ok {
		return nil, errors.New("invalid result")
	}
	return e.convertSymbol(slotDesk), nil
}

func (e *engine) convertSymbol(slotDesk *pb.SlotDesk) *pb.SlotDesk {
	e.convertSlotMaxtrix(slotDesk.Matrix)
	e.convertSlotMaxtrix(slotDesk.SpreadMatrix)
	e.convertSpinSymbol(slotDesk.Matrix.SpinLists)
	e.convertSpinSymbol(slotDesk.SpreadMatrix.SpinLists)
	e.convertSpinSymbol(slotDesk.SpinSymbols)
	e.convertPayline(slotDesk.Paylines)
	return slotDesk
}

func (e *engine) convertSlotMaxtrix(sm *pb.SlotMatrix) {
	entity.PbSlotMatrixForEeach(sm, func(idx int, row, col int32, symbol pb.SiXiangSymbol) {
		newSymbol := e.useFnSymbolAlias(symbol)
		sm.Lists[idx] = newSymbol
	})
}

func (e *engine) convertSpinSymbol(sm []*pb.SpinSymbol) {
	for _, spin := range sm {
		newSymbol := e.useFnSymbolAlias(spin.GetSymbol())
		spin.Symbol = newSymbol
	}
}

func (e *engine) convertPayline(sm []*pb.Payline) {
	for _, payline := range sm {
		newSymbol := e.useFnSymbolAlias(payline.GetSymbol())
		payline.Symbol = newSymbol
	}
}
