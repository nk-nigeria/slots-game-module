package entity

import pb "github.com/ciaolink-game-platform/cgp-common/proto"

const (
	ColsJuiceMatrix  = 5
	RowsJuicenMatrix = 4
)

var JuiceLowSymbols = []pb.SiXiangSymbol{
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_J,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_Q,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_K,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_A,
}

var JuiceMidSymbol = []pb.SiXiangSymbol{
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STRAWBERRY,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_WATERMELON,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_PINAPPLE,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_MANGOSTEEN,
}

var JuiceHighSymbol = []pb.SiXiangSymbol{
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STONE_1,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STONE_2,
	pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STONE_3,
}
var JuiceSymbols = []pb.SiXiangSymbol{}
