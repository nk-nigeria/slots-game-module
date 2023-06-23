package main

import (
	"context"
	"database/sql"
	"time"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/api"
	"github.com/ciaolink-game-platform/cgp-common/define"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/heroiclabs/nakama-common/runtime"
)

func InitModule(_ context.Context, logger runtime.Logger, _ *sql.DB, _ runtime.NakamaModule, initializer runtime.Initializer) error {
	initStart := time.Now()

	marshaler := &protojson.MarshalOptions{
		UseEnumNumbers:  true,
		EmitUnpopulated: true,
	}
	unmarshaler := &protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}
	gameNames := []string{define.SixiangGameName, define.TarzanGameName, define.JuicyGarden}
	for _, gameName := range gameNames {
		name := gameName
		if err := initializer.RegisterMatch(name, func(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) (runtime.Match, error) {
			return api.NewMatchHandler(name, marshaler, unmarshaler), nil
		}); err != nil {
			return err
		}
	}
	logger.Info("Plugin loaded in '%d' msec.", time.Since(initStart).Milliseconds())
	return nil
}
