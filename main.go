package main

import (
	"context"
	"database/sql"
	"time"

	"github.com/nakamaFramework/cgb-slots-game-module/api"
	"github.com/nakamaFramework/cgp-common/define"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/heroiclabs/nakama-common/runtime"
)

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	initStart := time.Now()

	marshaler := &protojson.MarshalOptions{
		UseEnumNumbers:  true,
		EmitUnpopulated: true,
	}
	unmarshaler := &protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}
	gameNames := []define.GameName{
		define.SixiangGameName,
		define.TarzanGameName,
		define.JuicyGardenName,
		define.IncaGameName,
		//clone game
		define.CryptoRush,
	}
	for _, gameName := range gameNames {
		name := gameName.String()
		if err := initializer.RegisterMatch(name, func(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) (runtime.Match, error) {
			return api.NewMatchHandler(name, marshaler, unmarshaler), nil
		}); err != nil {
			return err
		}
	}
	logger.Info("Plugin loaded in '%d' msec.", time.Since(initStart).Milliseconds())
	return nil
}
