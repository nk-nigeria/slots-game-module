package engine

import (
	"errors"
	"math/rand"
	"time"

	pb "github.com/ciaolink-game-platform/cgp-common/proto"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
)

var (
	ErrorSpinReadMax       = errors.New("Spin reach max")
	ErrorMissingSpinSymbol = errors.New("Missing spin symbol")
	ErrorNoGameEngine      = errors.New("No game engine")
)

func RandomInt(min, max int) int {
	if min < 0 {
		min = 0
	}
	if max <= min {
		max = min + 1
	}
	n := max - min
	rand.Seed(time.Now().UTC().UnixNano())
	num := rand.Intn(n)
	return num + min
}

func RandomFloat64(min, max float64) float64 {
	if min < 0 {
		min = 0.0
	}
	var ratio float64
	var m1 float64
	var m2 float64
	if min < 1.0 || max < 1.0 {
		ratio = 1000
		m1 = min * ratio
		m2 = max * ratio
	} else {
		ratio = 100
		m1 = min * ratio
		m2 = max * ratio
	}
	n := RandomInt(int(m1), int(m2))

	num := float64(n) / float64(ratio)
	if num > max {
		num = max
	}
	return num
}

func ShuffleMatrix(matrix entity.SlotMatrix) entity.SlotMatrix {
	list := matrix.List
	matrix.List = ShuffleSlice(list)
	return matrix
}

func ShuffleSlice[T any](slice []T) []T {
	rand.Seed(time.Now().UTC().UnixNano())
	ml := make([]T, len(slice), len(slice))
	copy(ml, slice)
	rand.Shuffle(len(ml), func(i, j int) { ml[i], ml[j] = ml[j], ml[i] })
	return ml
}

func LuckySymbolToReward(symbol pb.SiXiangSymbol) (pb.BigWin, pb.WinJackpot) {
	var bigWin pb.BigWin
	var winJp pb.WinJackpot
	switch symbol {
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_MINOR:
		bigWin = pb.BigWin_BIG_WIN_NICE
		winJp = pb.WinJackpot_WIN_JACKPOT_MINOR
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_MAJOR:
		bigWin = pb.BigWin_BIG_WIN_MEGA
		winJp = pb.WinJackpot_WIN_JACKPOT_MAJOR
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_MEGA:
		bigWin = pb.BigWin_BIG_WIN_MEGA
		winJp = pb.WinJackpot_WIN_JACKPOT_MEGA
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_LUCKYDRAW_GRAND:
		bigWin = pb.BigWin_BIG_WIN_MEGA
		winJp = pb.WinJackpot_WIN_JACKPOT_GRAND
	}
	return bigWin, winJp
}
