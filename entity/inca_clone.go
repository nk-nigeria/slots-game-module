package entity

import (
	pb "github.com/nakamaFramework/cgp-common/proto"
)

var (
	mapSymbolNoelFromInca  map[pb.SiXiangSymbol]pb.SiXiangSymbol
	mapSymbolFruitFromInca map[pb.SiXiangSymbol]pb.SiXiangSymbol
)

func init() {
	initMapSymbolNoelFromInca()
	initMapSymbolFruitFromInca()
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

func initMapSymbolFruitFromInca() {
	mapSymbolFruitFromInca[pb.SiXiangSymbol_SI_XIANG_SYMBOL_SUN] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STONE_GREEN
	mapSymbolFruitFromInca[pb.SiXiangSymbol_SI_XIANG_SYMBOL_EAGLE_GARUDA] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STONE_VIOLET
	mapSymbolFruitFromInca[pb.SiXiangSymbol_SI_XIANG_SYMBOL_ANTIQUE] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STONE_DIAMOND

	mapSymbolFruitFromInca[pb.SiXiangSymbol_SI_XIANG_SYMBOL_SUIT_CLUBS] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_STRAWBERRY
	mapSymbolFruitFromInca[pb.SiXiangSymbol_SI_XIANG_SYMBOL_SUIT_DIAMONDS] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_WATERMELON
	mapSymbolFruitFromInca[pb.SiXiangSymbol_SI_XIANG_SYMBOL_SUIT_HEARTS] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_PINAPPLE
	mapSymbolFruitFromInca[pb.SiXiangSymbol_SI_XIANG_SYMBOL_SUIT_SPADES] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_MANGOSTEEN

}

func GetSymbolNoelFromInca(sym pb.SiXiangSymbol) pb.SiXiangSymbol {
	noelSym, exist := mapSymbolNoelFromInca[sym]
	if !exist {
		noelSym = sym
	}
	return noelSym
}

func GetSymbolFruitFromInca(sym pb.SiXiangSymbol) pb.SiXiangSymbol {
	noelSym, exist := mapSymbolFruitFromInca[sym]
	if !exist {
		noelSym = sym
	}
	return noelSym
}
