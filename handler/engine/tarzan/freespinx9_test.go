package tarzan

import (
	"testing"

	"github.com/nakamaFramework/cgb-slots-game-module/entity"
	api "github.com/nakamaFramework/cgp-common/proto"
	"github.com/stretchr/testify/assert"
)

func Test_freespinx9_Process(t *testing.T) {
	name := "test Test_freespinx9_Process max spin"
	matchState := entity.NewSlotsMathState(nil)
	engine := NewFreeSpinX9(func(i1, i2 int) int {
		return i1
	})
	matchState.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_NORMAL
	engine.NewGame(matchState)
	t.Run(name, func(t *testing.T) {
		for i := 0; i < maxGemSpinFreeSpinX9; i++ {
			_, err := engine.Process(matchState)
			assert.NoError(t, err)
		}
		_, err := engine.Process(matchState)
		assert.ErrorIs(t, entity.ErrorSpinReachMax, err)
	})
}
