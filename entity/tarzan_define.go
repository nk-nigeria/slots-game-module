package entity

import pb "github.com/ciaolink-game-platform/cgp-common/proto"

const (
	ColsTarzanMatrix = 5
	RowsTarzanMatrix = 3
)

var TarzanSymbols = []pb.SiXiangSymbol{
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_J,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_Q,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_K,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_A,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_GORILLE,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_ELEPHANT,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JAGUAR,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_SNACK,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JANE,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JANE_FATHER,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_CLAYTON,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_FREE_SPIN,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_LETTER_J,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_LETTER_U,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_LETTER_N,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_LETTER_G,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_LETTER_L,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_LETTER_E,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_WILD,
}

var TarzanLowSymbol = map[pb.SiXiangSymbol]bool{
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_J:        true,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_Q:        true,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_K:        true,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_A:        true,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_GORILLE:  true,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_ELEPHANT: true,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JAGUAR:   true,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_SNACK:    true,
}

var TarzanMidSymbol = map[pb.SiXiangSymbol]bool{
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_GORILLE:  true,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_ELEPHANT: true,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JAGUAR:   true,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_SNACK:    true,
}

var TarzanHighSymbol = map[pb.SiXiangSymbol]bool{
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JANE:        true,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JANE_FATHER: true,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_CLAYTON:     true,
}

func NewTarzanMatrix() SlotMatrix {
	m := SlotMatrix{
		List: make([]pb.SiXiangSymbol, 0, ColsTarzanMatrix*RowsTarzanMatrix),
		Cols: ColsTarzanMatrix,
		Rows: RowsTarzanMatrix,
		Size: RowsTarzanMatrix * ColsTarzanMatrix,
	}

	return m
}
