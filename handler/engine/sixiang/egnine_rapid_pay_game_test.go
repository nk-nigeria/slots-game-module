package sixiang

import (
	"testing"

	"github.com/nakamaFramework/cgb-slots-game-module/entity"
	api "github.com/nakamaFramework/cgp-common/proto"
	"github.com/stretchr/testify/assert"
)

func Test_rapidPayEngine_NewGame(t *testing.T) {
	name := "Test_rapidPayEngine_NewGame"
	t.Run(name, func(t *testing.T) {
		e := NewRapidPayEngine(4, nil, nil)
		matchState := entity.NewSlotsMathState(nil)
		got, err := e.NewGame(matchState)
		assert.NotNil(t, got)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(matchState.SpinSymbols))
		assert.Equal(t, api.WinJackpot_WIN_JACKPOT_UNSPECIFIED, matchState.WinJp)
		assert.Equal(t, int(25), len(matchState.MatrixSpecial.List))
		assert.Equal(t, int(defaultRapidPayGemSpin), matchState.NumSpinLeft)
	})

}

func Test_rapidPayEngine_Process(t *testing.T) {
	type args struct {
		s          *entity.SlotsMatchState
		numProcess int
	}

	type want struct {
		gemSpin   []int
		trackFlip map[int]bool
	}
	type test struct {
		name string
		args args
		want want
	}

	engine := NewRapidPayEngine(4,
		func(min, max int) int { return min },
		func(min, max float64) float64 { return min },
	)
	tests := make([]test, 0)
	{
		// test process
		test := test{
			name: "Test_rapidPayEngine_Process",
		}
		s := entity.NewSlotsMathState(nil)
		s.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_RAPIDPAY
		engine.NewGame(s)
		test.args = args{
			s:          s,
			numProcess: 6,
		}
		want := want{
			gemSpin:   []int{4, 3, 2, 1, 0},
			trackFlip: map[int]bool{0: true, 5: true, 10: true, 15: true, 20: true},
		}
		test.want = want
		tests = append(tests, test)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < tt.args.numProcess; i++ {
				_, err := engine.Process(tt.args.s)
				if err != nil {
					assert.Equal(t, entity.ErrorSpinReachMax, err)
					continue
				}
				assert.Equal(t, tt.want.gemSpin[i], tt.args.s.NumSpinLeft)
				assert.Equal(t, 1, len(tt.args.s.SpinSymbols))
				assert.Equal(t, i+1, len(tt.args.s.MatrixSpecial.TrackFlip))
			}

		})
	}
}

func Test_rapidPayEngine_ProcessPlayGame(t *testing.T) {
	name := "Test_rapidPayEngine_Process"
	s := entity.NewSlotsMathState(nil)
	s.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_RAPIDPAY
	engine := NewRapidPayEngine(4, nil, nil)
	engine.NewGame(s)
	s.Bet().Chips = 1000
	t.Run(name, func(t *testing.T) {
		for {
			// s.Bet().Id = int32(engine.Random(0, 4))
			s.Bet().Id = 4
			res, err := engine.Process(s)
			if err != nil {
				break
			}
			assert.NotNil(t, res)
		}
	})
}

func Test_rapidPayEngine_Finish(t *testing.T) {

	type test struct {
		name    string
		args    *entity.SlotsMatchState
		want    *api.SlotDesk
		wantErr bool
	}
	engine := NewRapidPayEngine(4,
		func(min, max int) int { return min },
		func(min, max float64) float64 { return min },
	)
	tests := make([]test, 0)

	arrSym := []api.SiXiangSymbol{
		api.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X2,
		api.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X2,
		api.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X2,
		api.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X3,
		api.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X4,
	}
	// mapIdxFlip := map[int]int{0: 20, 1: 15, 2: 10, 3: 5, 4: 0}
	for _, sym := range arrSym {
		s := entity.NewSlotsMathState(nil)
		s.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_RAPIDPAY
		s.SetBetInfo(&api.InfoBet{Chips: 100})
		engine.NewGame(s)
		// sym := api.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X2
		// s.MatrixSpecial.TrackFlip[mapIdxFlip[idx]] = true
		idFlip := 0
		s.MatrixSpecial.ForEeach(func(idx, row, col int, symbol api.SiXiangSymbol) {
			if sym == symbol {
				idFlip = idx
			}
		})
		s.MatrixSpecial.Flip(idFlip)
		s.SpinSymbols = []*api.SpinSymbol{{Symbol: sym}}
		slotDesk := &api.SlotDesk{
			ChipsMcb:           s.Bet().Chips,
			CurrentSixiangGame: api.SiXiangGame_SI_XIANG_GAME_RAPIDPAY,
			NextSixiangGame:    api.SiXiangGame_SI_XIANG_GAME_RAPIDPAY,
		}
		slotDesk.GameReward.ChipsWin = int64((defaultAddRatioMcb + float64(entity.ListSymbolRapidPay[sym].Value.Min)) * float64(s.Bet().Chips))

		test := test{
			name: "Test_rapidPayEngine_Finish" + sym.String(),
			args: s,
			want: slotDesk,
		}
		tests = append(tests, test)
	}
	// test end
	{
		s := entity.NewSlotsMathState(nil)
		s.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_RAPIDPAY
		s.SetBetInfo(&api.InfoBet{Chips: 100})
		engine.NewGame(s)
		sym := api.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_END
		s.MatrixSpecial.TrackFlip[18] = true
		s.SpinSymbols = []*api.SpinSymbol{{Symbol: sym}}
		slotDesk := &api.SlotDesk{
			ChipsMcb:           s.Bet().Chips,
			CurrentSixiangGame: api.SiXiangGame_SI_XIANG_GAME_RAPIDPAY,
			NextSixiangGame:    api.SiXiangGame_SI_XIANG_GAME_NORMAL,
		}
		slotDesk.GameReward.ChipsWin = int64((defaultAddRatioMcb + float64(entity.ListSymbolRapidPay[sym].Value.Min)) * float64(s.Bet().Chips))

		test := test{
			name: "Test_rapidPayEngine_Finish" + sym.String(),
			args: s,
			want: slotDesk,
		}
		tests = append(tests, test)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := engine.Finish(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("rapidPayEngine.Finish() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, got)
			slotDesk := got.(*api.SlotDesk)
			assert.Equal(t, tt.want.ChipsMcb, slotDesk.ChipsMcb)
			assert.Equal(t, tt.want.CurrentSixiangGame, slotDesk.CurrentSixiangGame)
			assert.Equal(t, tt.want.NextSixiangGame, slotDesk.NextSixiangGame)
			assert.Equal(t, tt.want.GameReward.ChipsWin, slotDesk.GameReward.ChipsWin)
		})
	}
}
