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
	PriceBuySixiangGem[3] = 250
}

type SixiangSaveGame struct {
	GameEyePlayed map[int]map[pb.SiXiangGame]int `json:"game_eye_played,omitempty"`
	LastMcb       int64                          `json:"last_mcb,omitempty"`
}

type GameConfig struct {
	*pb.GameConfig
	AddGiftSpin bool `json:"add_gift_spin,omitempty"`
}

type JuiceSaveGame struct {
	LastMcb        int64                   `json:"last_mcb,omitempty"`
	ChipsAccumByJp map[pb.WinJackpot]int64 `json:"chips_accum_by_jp,omitempty"`
	GamePlaying    pb.SiXiangGame          `json:"game_playing,omitempty"`
	GameConfig     *GameConfig             `json:"game_config,omitempty"`
	NumSpinLeft    int                     `json:"num_spin_left,omitempty"`
	TotalChipWin   int                     `json:"total_chip_win,omitempty"`
	SpinList       []*pb.SpinSymbol        `json:"spin_list,omitempty"`
	MatrixSpecial  *SlotMatrix             `json:"matrix_special,omitempty"`
}

type TarzanSaveGame struct {
	LastMcb         int64              `json:"last_mcb,omitempty"`
	LetterSymbol    []pb.SiXiangSymbol `json:"letter_symbol,omitempty"`
	PerlGreenForest int                `json:"perl_green_forest,omitempty"`
	// chip tich luy
	PerlGreenForestChipsCollect  int64            `json:"perl_green_forest_chips_collect,omitempty"`
	GamePlaying                  pb.SiXiangGame   `json:"game_playing,omitempty"`
	NumSpinLeft                  int              `json:"num_spin_left,omitempty"`
	TotalChipWin                 int              `json:"total_chip_win,omitempty"`
	TotalLineWin                 int              `json:"total_line_win,omitempty"`
	CountLineCrossFreeSpinSymbol int              `json:"count_line_cross_free_spin_sym,omitempty"`
	SpinList                     []*pb.SpinSymbol `json:"spin_list,omitempty"`
	MatrixSpecial                *SlotMatrix      `json:"matrix_special,omitempty"`
	TurnSureSpinSpecial          int              `json:"turn_sure_spin_special,omitempty"`
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

	MatrixSpecial *SlotMatrix
	// ChipsWinInSpecialGame int64
	SpinSymbols []*pb.SpinSymbol
	NumSpinLeft int // gem using for spin in dragon perl
	// lần quay chắc chắn ra ngọc
	TurnSureSpinSpecial int
	// EyeSymbolRemains []pb.SiXiangSymbol
	// [mcb]gamebonus
	gameEyePlayed map[int]map[pb.SiXiangGame]int
	// Danh sach ngoc tứ linh spin được theo chip bet. [game][bet][symbol]qty_of_symbol
	// CollectionSymbol map[pb.SiXiangGame]map[int]map[pb.SiXiangSymbol]int
	SpinList []*pb.SpinSymbol

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
	PerlGreenForestChipsCollect int64
	NumScatterSeq               int
	NumFruitBasket              int

	LastSpinTime            time.Time
	DurationTriggerAutoSpin time.Duration
	NumSpinRemain6thLetter  int
	LetterSymbol            map[pb.SiXiangSymbol]bool `json:"letter_symbol,omitempty"`
	winJPHistory            *pb.JackpotHistory
	LastResult              *pb.SlotDesk
	Rtp                     lib.Rtp
	NotDropEyeSymbol        bool
	// chip accum by bet
	ChipsAccumByJp map[pb.WinJackpot]int64
	GameConfig     *GameConfig
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
		CurrentSiXiangGame:     pb.SiXiangGame_SI_XIANG_GAME_NORMAL,
		NextSiXiangGame:        pb.SiXiangGame_SI_XIANG_GAME_NORMAL,
		ChipStat:               NewChipStat(),
		gameEyePlayed:          make(map[int]map[pb.SiXiangGame]int),
		NumSpinRemain6thLetter: MinNumSpinLetter6th,
		LetterSymbol:           make(map[pb.SiXiangSymbol]bool),
		Rtp: lib.Rtp{
			Id:            0,
			PercentExpect: 120,
		},
		ChipsAccumByJp: make(map[pb.WinJackpot]int64),
	}
	m.winJPHistory = &pb.JackpotHistory{
		Mini: &pb.JackpotReward{
			WinJackpot: pb.WinJackpot_WIN_JACKPOT_MINI,
			Ratio:      int64(pb.WinJackpot_WIN_JACKPOT_MINOR),
		},
		Minor: &pb.JackpotReward{
			WinJackpot: pb.WinJackpot_WIN_JACKPOT_MINOR,
			Ratio:      int64(pb.WinJackpot_WIN_JACKPOT_MINOR),
		},
		Major: &pb.JackpotReward{
			WinJackpot: pb.WinJackpot_WIN_JACKPOT_MAJOR,
			Ratio:      int64(pb.WinJackpot_WIN_JACKPOT_MAJOR),
		},
		Mega: &pb.JackpotReward{
			WinJackpot: pb.WinJackpot_WIN_JACKPOT_MEGA,
			Ratio:      int64(pb.WinJackpot_WIN_JACKPOT_MEGA),
		},
		Grand: &pb.JackpotReward{
			WinJackpot: pb.WinJackpot_WIN_JACKPOT_GRAND,
			Ratio:      int64(pb.WinJackpot_WIN_JACKPOT_GRAND),
		},
	}
	m.GameConfig = &GameConfig{}
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
	// s.WildMatrix = SlotMatrix{}
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

func (s *SlotsMatchState) PriceBuySixiangGem() (int64, error) {
	if s.Label.Code != define.SixiangGameName.String() {
		return 0, ErrorInvalidRequestGame
	}

	numGemCollect := s.NumGameEyePlayed()
	if numGemCollect < 0 || numGemCollect > 4 {
		return 0, nil
	}
	ratio := PriceBuySixiangGem[numGemCollect]
	chips := int64(ratio) * s.Bet().Chips
	return chips, nil
}

func (s *SlotsMatchState) WinJPHistory() *pb.JackpotHistory {
	ratio := int64(1)
	switch s.CurrentSiXiangGame {
	case pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_DRAGON_PEARL,
		pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_LUCKDRAW,
		pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_GOLDPICK,
		pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_RAPIDPAY:
		ratio = 4
	}
	s.winJPHistory.Minor.Chips = s.winJPHistory.Minor.Ratio * s.bet.Chips * ratio
	s.winJPHistory.Major.Chips = s.winJPHistory.Major.Ratio * s.bet.Chips * ratio
	s.winJPHistory.Mega.Chips = s.winJPHistory.Mega.Ratio * s.bet.Chips * ratio
	s.winJPHistory.Grand.Chips = s.winJPHistory.Grand.Ratio * s.bet.Chips * ratio
	return s.winJPHistory
}

func (s *SlotsMatchState) WinJPHistoryJuice() *pb.JackpotHistory {
	{
		s.winJPHistory.Mini.Ratio = int64(JuiceJpRatio[pb.WinJackpot_WIN_JACKPOT_MINI])
		s.winJPHistory.Mini.Chips = s.Bet().Chips * s.winJPHistory.Mini.Ratio
		s.winJPHistory.Mini.ChipsAccum = 0
	}
	{
		s.winJPHistory.Minor.Ratio = int64(JuiceJpRatio[pb.WinJackpot_WIN_JACKPOT_MINOR])
		s.winJPHistory.Minor.Chips = s.Bet().Chips * s.winJPHistory.Minor.Ratio
		s.winJPHistory.Minor.ChipsAccum = 0
	}
	{
		s.winJPHistory.Major.Ratio = int64(JuiceJpRatio[pb.WinJackpot_WIN_JACKPOT_MAJOR])
		s.winJPHistory.Major.Chips = s.Bet().Chips * s.winJPHistory.Major.Ratio
		s.winJPHistory.Major.ChipsAccum = s.ChipsAccumByJp[pb.WinJackpot_WIN_JACKPOT_MAJOR] / 1000
	}
	{
		s.winJPHistory.Grand.Ratio = int64(JuiceJpRatio[pb.WinJackpot_WIN_JACKPOT_GRAND])
		s.winJPHistory.Grand.Chips = s.Bet().Chips * s.winJPHistory.Grand.Ratio
		s.winJPHistory.Grand.ChipsAccum = s.ChipsAccumByJp[pb.WinJackpot_WIN_JACKPOT_GRAND] / 200
	}
	return s.winJPHistory
}

func (s *SlotsMatchState) AddChipAccum(chips int64) {
	for _, jp := range []pb.WinJackpot{pb.WinJackpot_WIN_JACKPOT_MAJOR, pb.WinJackpot_WIN_JACKPOT_GRAND} {
		val := s.ChipsAccumByJp[jp]
		val += chips
		s.ChipsAccumByJp[jp] = val
	}
}

func (s *SlotsMatchState) GetAndResetChipAccumt(jp pb.WinJackpot) int64 {
	val := s.ChipsAccumByJp[jp]
	s.ChipsAccumByJp[jp] = 0
	return val
}

// func (s *SlotsMatchState) ChipAccum() int64 {
// 	v := s.chipsAccum[s.bet.Chips]
// 	return v
// }

func (s *SlotsMatchState) LoadSaveGame(saveGame *pb.SaveGame, suggestMcb func(mcbInSaveGame int64) int64) {
	defer func() {
		if s.bet.Chips < 0 && suggestMcb != nil {
			s.bet.Chips = suggestMcb(0)
		}
	}()
	// save game expire
	if saveGame.LastUpdateUnix == 0 || time.Now().Unix()-saveGame.LastUpdateUnix > 30*86400 {
		return
	}
	switch s.Label.Code {
	case define.SixiangGameName.String():
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
		if suggestMcb != nil {
			s.bet.Chips = suggestMcb(sixiangSaveGame.LastMcb)
		}
		if len(s.gameEyePlayed[int(s.Bet().Chips)]) == len(ListEyeSiXiang) {
			s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS
		}
	case define.TarzanGameName.String():
		tarzanSg := &TarzanSaveGame{}
		err := json.Unmarshal([]byte(saveGame.Data), &tarzanSg)
		if err != nil {
			return
		}
		s.bet = &pb.InfoBet{
			Chips: tarzanSg.LastMcb,
		}
		if suggestMcb != nil {
			s.bet.Chips = suggestMcb(tarzanSg.LastMcb)
		}
		if s.LetterSymbol == nil {
			s.LetterSymbol = make(map[pb.SiXiangSymbol]bool)
		}
		for _, sym := range tarzanSg.LetterSymbol {
			s.LetterSymbol[sym] = true
		}
		s.PerlGreenForest = tarzanSg.PerlGreenForest
		s.PerlGreenForestChipsCollect = tarzanSg.PerlGreenForestChipsCollect
		if tarzanSg.GamePlaying == pb.SiXiangGame_SI_XIANG_GAME_UNSPECIFIED {
			tarzanSg.GamePlaying = pb.SiXiangGame_SI_XIANG_GAME_NORMAL
		}
		s.NextSiXiangGame = tarzanSg.GamePlaying
		s.NumSpinLeft = tarzanSg.NumSpinLeft
		s.ChipStat.Reset(s.NextSiXiangGame)
		s.ChipStat.AddChipWin(s.NextSiXiangGame, int64(tarzanSg.TotalChipWin))
		s.ChipStat.AddLineWin(s.NextSiXiangGame, int64(tarzanSg.TotalLineWin))
		s.CountLineCrossFreeSpinSymbol = tarzanSg.CountLineCrossFreeSpinSymbol
		s.TurnSureSpinSpecial = 0
		s.SpinList = tarzanSg.SpinList
		switch s.NextSiXiangGame {
		case pb.SiXiangGame_SI_XIANG_GAME_TARZAN_FREESPINX9:
			s.MatrixSpecial = nil
			s.LastResult = nil
		case pb.SiXiangGame_SI_XIANG_GAME_TARZAN_JUNGLE_TREASURE:
			s.NumSpinLeft = tarzanSg.NumSpinLeft
			s.MatrixSpecial = tarzanSg.MatrixSpecial
			s.LastResult = nil
		default:
			s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_NORMAL
		}
		if len(s.LetterSymbol) == len(TarzanLetterSymbol) {
			s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_TARZAN_JUNGLE_TREASURE
		}
	case define.JuicyGardenName.String(),
		define.CrytoRush.String():
		{
			juiceSg := &JuiceSaveGame{}
			err := json.Unmarshal([]byte(saveGame.Data), &juiceSg)
			if err != nil {
				return
			}
			s.bet = &pb.InfoBet{
				Chips: juiceSg.LastMcb,
			}
			if suggestMcb != nil {
				s.bet.Chips = suggestMcb(juiceSg.LastMcb)
			}
			s.ChipsAccumByJp = juiceSg.ChipsAccumByJp
			if s.ChipsAccumByJp == nil {
				s.ChipsAccumByJp = make(map[pb.WinJackpot]int64)
			}
			s.NextSiXiangGame = juiceSg.GamePlaying
			if s.NextSiXiangGame == pb.SiXiangGame_SI_XIANG_GAME_UNSPECIFIED {
				s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_NORMAL
			}
			s.NumSpinLeft = juiceSg.NumSpinLeft
			s.ChipStat.Reset(s.NextSiXiangGame)
			s.ChipStat.AddChipWin(s.NextSiXiangGame, int64(juiceSg.TotalChipWin))
			s.MatrixSpecial = juiceSg.MatrixSpecial
			s.SpinList = juiceSg.SpinList
			s.LastResult = nil
			s.GameConfig = juiceSg.GameConfig
			if s.GameConfig == nil {
				s.GameConfig = &GameConfig{}
			}
		}
	}
}

func (s *SlotsMatchState) SaveGameJson() string {
	// return "test"
	var saveGameInf interface{}
	switch s.Label.Code {
	case define.SixiangGameName.String():
		sixiangSaveGame := &SixiangSaveGame{
			GameEyePlayed: s.gameEyePlayed,
			LastMcb:       s.bet.Chips,
		}
		saveGameInf = sixiangSaveGame
	case define.TarzanGameName.String():
		tarzanSg := &TarzanSaveGame{
			LetterSymbol:                 make([]pb.SiXiangSymbol, 0),
			PerlGreenForest:              s.PerlGreenForest,
			PerlGreenForestChipsCollect:  s.PerlGreenForestChipsCollect,
			LastMcb:                      s.bet.Chips,
			GamePlaying:                  s.NextSiXiangGame,
			NumSpinLeft:                  s.NumSpinLeft,
			TotalChipWin:                 int(s.ChipStat.TotalChipWin(s.NextSiXiangGame)),
			TotalLineWin:                 int(s.ChipStat.TotalLineWin(s.NextSiXiangGame)),
			CountLineCrossFreeSpinSymbol: s.CountLineCrossFreeSpinSymbol,
			SpinList:                     s.SpinList,
			TurnSureSpinSpecial:          s.TurnSureSpinSpecial,
			MatrixSpecial:                s.MatrixSpecial,
		}
		if len(tarzanSg.SpinList) == 0 {
			tarzanSg.SpinList = nil
		}
		if tarzanSg.GamePlaying == pb.SiXiangGame_SI_XIANG_GAME_NORMAL || tarzanSg.NumSpinLeft <= 0 {
			tarzanSg.TotalChipWin = 0
			tarzanSg.TotalLineWin = 0
			tarzanSg.GamePlaying = pb.SiXiangGame_SI_XIANG_GAME_NORMAL
			tarzanSg.NumSpinLeft = 0
			tarzanSg.CountLineCrossFreeSpinSymbol = 0
			tarzanSg.TurnSureSpinSpecial = 0
		}
		for sym := range s.LetterSymbol {
			tarzanSg.LetterSymbol = append(tarzanSg.LetterSymbol, sym)
		}
		saveGameInf = tarzanSg
	case define.JuicyGardenName.String(),
		define.CrytoRush.String():
		{
			saveGame := JuiceSaveGame{
				LastMcb:        s.bet.Chips,
				ChipsAccumByJp: s.ChipsAccumByJp,
				GameConfig:     s.GameConfig,
				GamePlaying:    s.NextSiXiangGame,
				NumSpinLeft:    s.NumSpinLeft,
				TotalChipWin:   int(s.ChipStat.TotalChipWin(s.NextSiXiangGame)),
				SpinList:       s.SpinList,
				MatrixSpecial:  s.MatrixSpecial,
			}
			if s.NextSiXiangGame == pb.SiXiangGame_SI_XIANG_GAME_NORMAL {
				saveGame.TotalChipWin = 0
			}
			saveGameInf = saveGame
		}
	}
	data, _ := json.Marshal(saveGameInf)
	return string(data)
}
