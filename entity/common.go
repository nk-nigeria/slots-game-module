package entity

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/bwmarrin/snowflake"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

var DefaultMarshaler = &protojson.MarshalOptions{
	UseEnumNumbers:  true,
	EmitUnpopulated: false,
}
var DefaulUnmarshaler = &protojson.UnmarshalOptions{
	DiscardUnknown: false,
}

const (
	MaxPresenceCard   = 13
	JackpotPercentTax = 1 // 1%
	TickRate          = 2
)

var (
	ErrorSpinReachMax         = errors.New("spin reach max")
	ErrorMissingSpinSymbol    = errors.New("missing spin symbol")
	ErrorNoGameEngine         = errors.New("no game engine")
	ErrorSpinNotChange        = errors.New("spin not change")
	ErrorChipNotEnough        = errors.New("chip is not enough")
	ErrorSpinIndexAleadyTaken = errors.New("spin index already taken")
	ErrorSpinIndexRequired    = errors.New("spin index required")
	ErrorInfoBetInvalid       = errors.New("info bet invalid")
	ErrorInvalidRequestGame   = errors.New("invalid request game")
	ErrorInternal             = errors.New("internal error")
)

// free game by lv
// [level]=%
// https://docs.google.com/spreadsheets/d/1OKPtCzTGe5Da-HRUKe37rS3bIGYw4F_B/edit#gid=1754766987
var feeGameByLvPercent = map[int]int{0: 7, 1: 7, 2: 6, 3: 5, 4: 5, 6: 4, 7: 4, 8: 4, 9: 4, 10: 4}

func GetFeeGameByLevel(lv int) int {
	val, exist := feeGameByLvPercent[lv]
	if exist {
		return val
	}
	return 5
}

var SnowlakeNode, _ = snowflake.NewNode(1)

type WalletAction string

const (
	WalletActionWinGameJackpot WalletAction = "win_game_jackpot"
	// WalletActionGameFee        WalletAction = "game_fee"
)

func InterfaceToString(inf interface{}) string {
	if inf == nil {
		return ""
	}
	str, ok := inf.(string)
	if !ok {
		return ""
	}
	return str
}

func ToInt64(inf interface{}, def int64) int64 {
	if inf == nil {
		return def
	}
	switch v := inf.(type) {
	case int:
		return int64(inf.(int))
	case int64:
		return inf.(int64)
	case string:
		str := inf.(string)
		i, _ := strconv.ParseInt(str, 10, 64)
		return i
	case float64:
		return int64(inf.(float64))
	default:
		fmt.Printf("I don't know about type %T!\n", v)
	}
	return def
}

func MinInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func MaxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func SliceRepeat[T any](size int, v T) []T {
	retval := make([]T, 0, size)
	for i := 0; i < size; i++ {
		retval = append(retval, v)
	}
	return retval
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

func GoldPickSymbolToReward(symbol pb.SiXiangSymbol) (pb.BigWin, pb.WinJackpot) {
	var bigWin pb.BigWin = pb.BigWin_BIG_WIN_UNSPECIFIED
	var winJp pb.WinJackpot = pb.WinJackpot_WIN_JACKPOT_UNSPECIFIED
	switch symbol {
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_GOLD_PICK_JP_MINOR:
		bigWin = pb.BigWin_BIG_WIN_NICE
		winJp = pb.WinJackpot_WIN_JACKPOT_MINOR
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_GOLD_PICK_JP_MAJOR:
		bigWin = pb.BigWin_BIG_WIN_MEGA
		winJp = pb.WinJackpot_WIN_JACKPOT_MAJOR
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_GOLD_PICK_JP_MEGA:
		bigWin = pb.BigWin_BIG_WIN_MEGA
		winJp = pb.WinJackpot_WIN_JACKPOT_MEGA
	}
	return bigWin, winJp
}

func DragonPearlSymbolToReward(symbol pb.SiXiangSymbol) (pb.BigWin, pb.WinJackpot) {
	var bigWin pb.BigWin = pb.BigWin_BIG_WIN_UNSPECIFIED
	var winJp pb.WinJackpot = pb.WinJackpot_WIN_JACKPOT_UNSPECIFIED
	switch symbol {
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_JP_MINOR:
		bigWin = pb.BigWin_BIG_WIN_NICE
		winJp = pb.WinJackpot_WIN_JACKPOT_MINOR
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_JP_MAJOR:
		bigWin = pb.BigWin_BIG_WIN_MEGA
		winJp = pb.WinJackpot_WIN_JACKPOT_MAJOR
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_JP_MEGA:
		bigWin = pb.BigWin_BIG_WIN_MEGA
		winJp = pb.WinJackpot_WIN_JACKPOT_MEGA
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_JP_GRAND:
		bigWin = pb.BigWin_BIG_WIN_MEGA
		winJp = pb.WinJackpot_WIN_JACKPOT_GRAND
	}
	return bigWin, winJp
}
