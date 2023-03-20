package sixiangengine

import (
	"testing"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	api "github.com/ciaolink-game-platform/cgp-common/proto"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
	"github.com/stretchr/testify/assert"
)

func Test_sixiangBonusIngameEngine_NewGame(t *testing.T) {
	type fields struct {
		ratioBonus  int
		enginesGame map[pb.SiXiangGame]lib.Engine
	}
	type args struct {
		matchState interface{}
	}
	type test struct {
		name    string
		args    args
		want    lib.Engine
		wantErr bool
	}
	tests := make([]test, 0)
	arrGame := []pb.SiXiangGame{
		pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_DRAGON_PEARL,
		pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_LUCKDRAW,
		pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_GOLDPICK,
		pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_RAPIDPAY,
	}
	for _, game := range arrGame {
		matchState := entity.NewSlotsMathState(nil)
		matchState.SetBetInfo(&pb.InfoBet{
			Chips: 100,
		})
		matchState.CurrentSiXiangGame = game
		entity.NewSlotsMathState(nil)
		var expectEngine lib.Engine
		switch game {
		case pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_DRAGON_PEARL:
			expectEngine = NewDragonPearlEngine(nil, nil)
		case pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_LUCKDRAW:
			expectEngine = NewLuckyDrawEngine(nil, nil)
		case pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_GOLDPICK:
			expectEngine = NewGoldPickEngine(nil, nil)
		case pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_RAPIDPAY:
			expectEngine = NewRapidPayEngine(nil, nil)
		}
		tests = append(tests, test{
			name: "Test_sixiangBonusIngameEngine_NewGame_" + game.String(),
			args: args{
				matchState: matchState,
			},
			want: expectEngine,
		})
	}
	e := NewSixiangBonusInGameEngine(4)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := e.NewGame(tt.args.matchState)
			expectGot, expectErr := tt.want.NewGame(tt.args.matchState)
			assert.Equal(t, expectErr, err)
			assert.Equal(t, expectGot, got)
		})
	}
}

func Test_sixiangBonusIngameEngine_Process(t *testing.T) {
	type fields struct {
		ratioBonus  int
		enginesGame map[pb.SiXiangGame]lib.Engine
	}
	type args struct {
		matchState interface{}
	}
	type test struct {
		name    string
		args    args
		want    lib.Engine
		wantErr bool
	}
	tests := make([]test, 0)
	arrGame := []pb.SiXiangGame{
		api.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_DRAGON_PEARL,
		api.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_LUCKDRAW,
		api.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_GOLDPICK,
		api.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_RAPIDPAY,
	}
	for _, game := range arrGame {
		matchState := entity.NewSlotsMathState(nil)
		matchState.SetBetInfo(&pb.InfoBet{
			Chips: 100,
		})
		matchState.CurrentSiXiangGame = game
		entity.NewSlotsMathState(nil)
		var expectEngine lib.Engine
		switch game {
		case api.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_DRAGON_PEARL:
			expectEngine = NewDragonPearlEngine(nil, nil)
		case api.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_LUCKDRAW:
			expectEngine = NewLuckyDrawEngine(nil, nil)
		case api.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_GOLDPICK:
			expectEngine = NewGoldPickEngine(nil, nil)
		case api.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_RAPIDPAY:
			expectEngine = NewRapidPayEngine(nil, nil)
		}
		tests = append(tests, test{
			name: "Test_sixiangBonusIngameEngine_NewGame_" + game.String(),
			args: args{
				matchState: matchState,
			},
			want: expectEngine,
		})
	}
	e := NewSixiangBonusInGameEngine(4)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := e.NewGame(tt.args.matchState)
			expectGot, expectErr := tt.want.NewGame(tt.args.matchState)
			assert.Equal(t, expectErr, err)
			assert.Equal(t, expectGot, got)
			got, err = e.Process(tt.args.matchState)
			expectGot, expectErr = tt.want.Process(tt.args.matchState)
			assert.Equal(t, expectErr, err)
			assert.Equal(t, expectGot, got)
		})
	}
}
