package noel

import (
	"errors"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgb-slots-game-module/handler/engine/inca"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

type engine struct {
	incaEngine lib.Engine
}

func NewEngine() lib.Engine {
	e := &engine{
		incaEngine: inca.NewEngine(),
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
		newSymbol := entity.GetSymbolNoelFromInca(symbol)
		sm.Lists[idx] = newSymbol
	})
}

func (e *engine) convertSpinSymbol(sm []*pb.SpinSymbol) {
	for _, spin := range sm {
		newSymbol := entity.GetSymbolNoelFromInca(spin.GetSymbol())
		spin.Symbol = newSymbol
	}
}

func (e *engine) convertPayline(sm []*pb.Payline) {
	for _, payline := range sm {
		newSymbol := entity.GetSymbolNoelFromInca(payline.GetSymbol())
		payline.Symbol = newSymbol
	}
}
