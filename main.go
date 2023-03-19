package main

import (
	"context"
	"database/sql"
	"time"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/api"
	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"

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
	if err := initializer.RegisterMatch(entity.ModuleName, func(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) (runtime.Match, error) {
		return api.NewMatchHandler(marshaler, unmarshaler), nil
	}); err != nil {
		return err
	}

	logger.Info("Plugin loaded in '%d' msec.", time.Since(initStart).Milliseconds())
	return nil
}
