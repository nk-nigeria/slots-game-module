package entity

import (
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

var mapSymbolNoelFromInca map[pb.SiXiangSymbol]pb.SiXiangSymbol

func init() {
	initMapSymbolNoelFromInca()
}

func initMapSymbolNoelFromInca() {
	mapSymbolNoelFromInca[pb.SiXiangSymbol_SI_XIANG_SYMBOL_SUN] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_CHRISMAS_GIFT
	mapSymbolNoelFromInca[pb.SiXiangSymbol_SI_XIANG_SYMBOL_EAGLE_GARUDA] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_CHRISMAS_CANDY
	mapSymbolNoelFromInca[pb.SiXiangSymbol_SI_XIANG_SYMBOL_ANTIQUE] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_CHRISMAS_RING

	mapSymbolNoelFromInca[pb.SiXiangSymbol_SI_XIANG_SYMBOL_A] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_Q
	mapSymbolNoelFromInca[pb.SiXiangSymbol_SI_XIANG_SYMBOL_J] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_A
	mapSymbolNoelFromInca[pb.SiXiangSymbol_SI_XIANG_SYMBOL_Q] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_J
	mapSymbolNoelFromInca[pb.SiXiangSymbol_SI_XIANG_SYMBOL_K] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_K

}

func GetSymbolNoelFromInca(sym pb.SiXiangSymbol) pb.SiXiangSymbol {
	noelSym, exist := mapSymbolNoelFromInca[sym]
	if !exist {
		noelSym = sym
	}
	return noelSym
}
