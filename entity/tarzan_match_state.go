package entity

import (
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

type TarzanMatchState struct {
	lib.MatchState
	CurrentSiXiangGame pb.SiXiangGame
	// next game in loop
	NextSiXiangGame  pb.SiXiangGame
	collectionSymbol map[pb.SiXiangSymbol]int

	Matrix        SlotMatrix
	SwingMatrix   SlotMatrix
	MatrixSpecial SlotMatrix
	SpinSymbols   []*pb.SpinSymbol

	Bet      *pb.InfoBet
	Paylines []*pb.Payline
	// List idx of free symbol index
	TrackIndexFreeSpinSymbol     map[int]bool
	ChipWinByGame                map[pb.SiXiangGame]int64
	CountLineCrossFreeSpinSymbol int
	GemSpin                      int // gem using for spin in freex9
	// ngoc rung xanh
	PerlGreeForest int
	// chip tich luy
	ChipsAccumulation int64
}

func NewTarzanMatchState(label *lib.MatchLabel) *TarzanMatchState {
	m := TarzanMatchState{
		MatchState: lib.NewMathState(label, NewMyPrecense),
		Bet: &pb.InfoBet{
			Chips: 0,
		},
		CurrentSiXiangGame: pb.SiXiangGame_SI_XIANG_GAME_NORMAL,
		NextSiXiangGame:    pb.SiXiangGame_SI_XIANG_GAME_NORMAL,
	}
	return &m
}

func (m *TarzanMatchState) AddCollectionSymbol(sym pb.SiXiangSymbol) {
	num := m.collectionSymbol[sym]
	num++
	m.collectionSymbol[sym] = num
}

func (m *TarzanMatchState) CollectionSymbolToSlice() []pb.SiXiangSymbol {
	ml := make([]pb.SiXiangSymbol, 0, len(m.collectionSymbol))
	for k := range m.collectionSymbol {
		ml = append(ml, k)
	}
	return ml
}

func (m *TarzanMatchState) SizeCollectionSymbol() int {
	return len(m.collectionSymbol)
}
