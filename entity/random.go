package entity

import (
	"crypto/rand"
	"math/big"
	mrand "math/rand"
	"time"
)

func RandomInt(min, max int) int {
	if min < 0 {
		min = 0
	}
	if min == max {
		return min
	}
	if max <= min {
		max = min + 1
	}
	n := max - min
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		mrand.Seed(time.Now().UTC().UnixNano())
		return mrand.Intn(n) + min
	}
	return int(nBig.Int64()) + min
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

func ShuffleMatrix(matrix SlotMatrix) SlotMatrix {
	list := matrix.List
	matrix.List = ShuffleSlice(list)
	return matrix
}

func ShuffleSlice[T any](slice []T) []T {
	// mrand.Seed(time.Now().UTC().UnixNano())
	ml := make([]T, len(slice))
	copy(ml, slice)
	// mrand.NewSource(time.Now().UTC().UnixNano())
	// mrand.Shuffle(len(ml), func(i, j int) { ml[i], ml[j] = ml[j], ml[i] })
	lenSl := len(slice)
	for i := 0; i < lenSl; i++ {
		j := RandomInt(0, lenSl)
		ml[i], ml[j] = ml[j], ml[i]
	}
	return ml
}
