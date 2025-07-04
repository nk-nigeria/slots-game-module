package main

import (
	"context"
	"database/sql"
	"time"

	"github.com/nk-nigeria/slots-game-module/api"
	"github.com/nk-nigeria/slots-game-module/entity"

	"github.com/heroiclabs/nakama-common/runtime"
	"google.golang.org/protobuf/proto"
)

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	initStart := time.Now()

	marshaler := &proto.MarshalOptions{}
	unmarshaler := &proto.UnmarshalOptions{
		DiscardUnknown: false,
	}
	if err := initializer.RegisterMatch(entity.ModuleName, func(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) (runtime.Match, error) {
		return api.NewMatchHandler(marshaler, unmarshaler), nil
	}); err != nil {
		return err
	}

	logger.Info("Plugin loaded in '%d' msec.", time.Now().Sub(initStart).Milliseconds())
	return nil
}
