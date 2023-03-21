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

	Bet      *pb.InfoBet
	Paylines []*pb.Payline
	// List idx of free symbol index
	FreeSpinSymbolIndexs []int
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
