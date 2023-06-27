package entity

import (
	"encoding/json"
	"time"

	"github.com/ciaolink-game-platform/cgp-common/define"
	"github.com/ciaolink-game-platform/cgp-common/lib"

	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

const (
	MinPresences        = 1
	MaxPresences        = 1
	MinNumSpinLetter6th = 280
)

var BetLevels []int64 = []int64{100, 200, 500, 1000}

// [num_gem_already_collect]priceratio
var PriceBuySixiangGem map[int]int

func init() {
	PriceBuySixiangGem = make(map[int]int)
	PriceBuySixiangGem[0] = 90
	PriceBuySixiangGem[1] = 110
	PriceBuySixiangGem[2] = 150
	PriceBuySixiangGem[1] = 250
}

type SixiangSaveGame struct {
	GameEyePlayed map[int]map[pb.SiXiangGame]int `json:"game_eye_played,omitempty"`
	LastMcb       int64                          `json:"last_mcb,omitempty"`
}

type SlotsMatchState struct {
	// SaveGame map[string]*pb.SaveGame
	lib.MatchState
	// prevent calc reward multil time,
	IsSpinChange  bool
	allowSpin     bool // allow user submit new bet
	balanceResult *pb.BalanceResult
	// UserDataMatch  UserDataMatch
	Matrix         SlotMatrix
	WildMatrix     SlotMatrix
	paylines       []*pb.Payline
	WaitSpinMatrix bool

	CurrentSiXiangGame pb.SiXiangGame
	// next game in loop
	NextSiXiangGame pb.SiXiangGame
	bet             *pb.InfoBet
	WinJp           pb.WinJackpot

	MatrixSpecial SlotMatrix
	// ChipsWinInSpecialGame int64
	SpinSymbols []*pb.SpinSymbol
	NumSpinLeft int // gem using for spin in dragon perl
	// lần quay chắc chắn ra ngọc
	TurnSureSpin     int
	EyeSymbolRemains []pb.SiXiangSymbol
	// [mcb]gamebonus
	gameEyePlayed map[int]map[pb.SiXiangGame]int
	// Danh sach ngoc tứ linh spin được theo chip bet. [game][bet][symbol]qty_of_symbol
	CollectionSymbol map[pb.SiXiangGame]map[int]map[pb.SiXiangSymbol]int
	SpinList         []*pb.SpinSymbol

	// tarzan
	// List idx of free symbol index
	TrackIndexFreeSpinSymbol map[int]bool
	// ChipWinByGame            map[pb.SiXiangGame]int64
	// LineWinByGame map[pb.SiXiangGame]int
	ChipStat *chipStat
	// so luong payline di qua free spin tarzan
	CountLineCrossFreeSpinSymbol int
	// ngoc rung xanh
	PerlGreenForest int
	// chip tich luy
	ChipsBonus       int64
	NumScatterSeq    int
	NumFruitBasket   int
	RatioFruitBasket int

	LastSpinTime            time.Time
	DurationTriggerAutoSpin time.Duration
	NumSpinRemain6thLetter  int
}

func NewSlotsMathState(label *lib.MatchLabel) *SlotsMatchState {
	m := SlotsMatchState{
		IsSpinChange:   false,
		MatchState:     lib.NewMathState(label, NewMyPrecense),
		balanceResult:  nil,
		WaitSpinMatrix: false,
		bet: &pb.InfoBet{
			Chips: 0,
		},
		CurrentSiXiangGame: pb.SiXiangGame_SI_XIANG_GAME_NORMAL,
		NextSiXiangGame:    pb.SiXiangGame_SI_XIANG_GAME_NORMAL,
		CollectionSymbol:   make(map[pb.SiXiangGame]map[int]map[pb.SiXiangSymbol]int, 0),
		// ChipWinByGame:      make(map[pb.SiXiangGame]int64, 0),
		// LineWinByGame:    make(map[pb.SiXiangGame]int, 0),
		ChipStat:               NewChipStat(),
		RatioFruitBasket:       1,
		gameEyePlayed:          make(map[int]map[pb.SiXiangGame]int),
		NumSpinRemain6thLetter: MinNumSpinLetter6th,
	}
	// m.SaveGame = make(map[string]*pb.SaveGame)

	return &m
}

func (s *SlotsMatchState) SetMatrix(matrix SlotMatrix) {
	s.Matrix = matrix
}

func (s *SlotsMatchState) SetWildMatrix(matrix SlotMatrix) {
	s.WildMatrix = matrix
}

func (s *SlotsMatchState) Paylines() []*pb.Payline {
	return s.paylines
}

func (s *SlotsMatchState) SetPaylines(paylines []*pb.Payline) {
	s.paylines = paylines
}

func (s *SlotsMatchState) Bet() *pb.InfoBet {
	return s.bet
}

func (s *SlotsMatchState) SetBetInfo(bet *pb.InfoBet) {
	s.bet = bet
}

func (s *SlotsMatchState) InitNewRound() {
	s.WaitSpinMatrix = false
	s.WinJp = pb.WinJackpot_WIN_JACKPOT_UNSPECIFIED
	s.paylines = nil
	s.WildMatrix = SlotMatrix{}
}
func (s *SlotsMatchState) IsAllowSpin() bool {
	return s.allowSpin
}

func (s *SlotsMatchState) SetAllowSpin(val bool) {
	s.allowSpin = val
}

func (s *SlotsMatchState) SetUpCountDown(duration time.Duration) {
	s.CountDownReachTime = time.Now().Add(duration)
	s.LastCountDown = -1
}

func (s *SlotsMatchState) ResetBalanceResult() {
	s.SetBalanceResult(nil)
}

func (s *SlotsMatchState) GetBalanceResult() *pb.BalanceResult {
	return s.balanceResult
}

func (s *SlotsMatchState) SetBalanceResult(u *pb.BalanceResult) {
	s.balanceResult = u
}

//	func (s *SlotsMatchState) ResetTrackingPlayBonusGame() {
//		s.trackingPlaySiXiangGame = make(map[pb.SiXiangGame]int)
//	}
func (s *SlotsMatchState) AddCollectionSymbol(game pb.SiXiangGame, chipMcb int, sym pb.SiXiangSymbol) {
	if _, ok := s.CollectionSymbol[game]; !ok {
		s.CollectionSymbol[game] = make(map[int]map[pb.SiXiangSymbol]int)
	}
	collection, exist := s.CollectionSymbol[game][chipMcb]
	if !exist {
		s.CollectionSymbol[game][chipMcb] = make(map[pb.SiXiangSymbol]int)
		collection = s.CollectionSymbol[game][chipMcb]
	}
	num := collection[sym]
	num++
	collection[sym] = num
	s.CollectionSymbol[game][chipMcb] = collection
}

func (s *SlotsMatchState) ResetCollection(game pb.SiXiangGame, chipMcb int) {
	s.CollectionSymbol[game] = make(map[int]map[pb.SiXiangSymbol]int)
	s.CollectionSymbol[game][chipMcb] = make(map[pb.SiXiangSymbol]int)
}

func (s *SlotsMatchState) CollectionSymbolToSlice(game pb.SiXiangGame, chipMcb int) []*pb.CollectSymbol {
	ml := make([]*pb.CollectSymbol, 0, len(s.CollectionSymbol))
	if _, ok := s.CollectionSymbol[game]; !ok {
		return ml
	}
	collection, exist := s.CollectionSymbol[game][chipMcb]
	if !exist {
		return ml
	}
	for symbol, qty := range collection {
		v := &pb.CollectSymbol{
			Symbol: symbol,
			Qty:    int64(qty),
		}
		ml = append(ml, v)
	}
	return ml
}

func (s *SlotsMatchState) SizeCollectionSymbol(game pb.SiXiangGame, chipMcb int) int {
	if _, ok := s.CollectionSymbol[game]; !ok {
		return 0
	}
	return len(s.CollectionSymbol[game][chipMcb])
}

func (s *SlotsMatchState) AddGameEyePlayed(game pb.SiXiangGame) int {
	if _, ok := s.gameEyePlayed[int(s.bet.Chips)]; !ok {
		s.gameEyePlayed[int(s.bet.Chips)] = make(map[pb.SiXiangGame]int)
	}
	m := s.gameEyePlayed[int(s.bet.Chips)]
	num := m[game]
	num++
	m[game] = num
	s.gameEyePlayed[int(s.bet.Chips)] = m
	return num
}

func (s *SlotsMatchState) NumGameEyePlayed() int {
	return len(s.gameEyePlayed[int(s.bet.Chips)])
}

func (s *SlotsMatchState) ClearGameEyePlayed() {
	s.gameEyePlayed[int(s.Bet().Chips)] = make(map[pb.SiXiangGame]int, 0)
}

func (s *SlotsMatchState) GameEyePlayed() map[pb.SiXiangGame]int {
	return s.gameEyePlayed[int(s.Bet().Chips)]
}

func (s *SlotsMatchState) LoadSaveGame(saveGame *pb.SaveGame) {
	// save game expire
	if time.Now().Unix()-saveGame.LastUpdateUnix > 30*86400 {
		return
	}
	if len(saveGame.Data) == 0 {
		return
	}
	switch s.Label.Code {
	case define.SixiangGameName:
		sixiangSaveGame := &SixiangSaveGame{}
		err := json.Unmarshal([]byte(saveGame.Data), &sixiangSaveGame)
		if err != nil {
			return
		}
		s.gameEyePlayed = sixiangSaveGame.GameEyePlayed
		if s.gameEyePlayed == nil {
			s.gameEyePlayed = make(map[int]map[pb.SiXiangGame]int)
		}
		s.bet = &pb.InfoBet{
			Chips: sixiangSaveGame.LastMcb,
		}
	}
}

func (s *SlotsMatchState) SaveGameJson() string {
	// return "test"
	switch s.Label.Code {
	case define.SixiangGameName:
		sixiangSaveGame := &SixiangSaveGame{
			GameEyePlayed: s.gameEyePlayed,
			LastMcb:       s.bet.Chips,
		}
		data, _ := json.Marshal(sixiangSaveGame)
		return string(data)
	}
	return ""
}
