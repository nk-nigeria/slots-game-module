package api

import (
	"context"
	"database/sql"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/nk-nigeria/cgp-common/lib"
	"github.com/nk-nigeria/cgp-common/presenter"
	"github.com/nk-nigeria/slots-game-module/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (m *MatchHandler) MatchJoinAttempt(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presence runtime.Presence, metadata map[string]string) (interface{}, bool, string) {
	s := state.(*entity.SlotsMatchState)
	logger.Info("match join attempt, state=%v, meta=%v", s, metadata)

	// check password
	if !s.Label.Open {
		logger.Info("match protect with password, check password")
		joinPassword := metadata["password"]
		if joinPassword != s.Label.Password {
			return s, false, "wrong password"
		}
	}
	// Check if it's a user attempting to rejoin after a disconnect.
	if p, _ := s.Presences.Get(presence.GetUserId()); p != nil {
		// 	// User rejoining after a disconnect.
		logger.Info("user %s rejoin after disconnect", presence.GetUserId())
		s.RemoveLeavePresence(presence.GetUserId())

		s.JoinsInProgress++
		return s, true, ""

	}

	// join as new user

	// Check if match is full.
	if s.Presences.Size()+s.JoinsInProgress >= entity.MaxPresences {
		return s, false, "match full"
	}
	// check chip balance in wallet before allow join
	wallet, err := entity.ReadWalletUser(ctx, nk, logger, presence.GetUserId())
	if err != nil {
		return s, false, status.Error(codes.Internal, "read chip balance failed").Error()
	}
	if wallet.Chips < int64(s.Label.Bet.GetMarkUnit()) {
		logger.Warn("[Reject] reject allow user %s join game, not enough chip join game, balance user chip %d , game bet %d",
			presence.GetUserId(), wallet.Chips, s.Label.Bet)
		return s, false, status.Error(codes.Internal, "chip balance not enough").Error()
	}

	// New player attempting to connect.
	s.JoinsInProgress++
	return s, true, ""
}

func (m *MatchHandler) MatchJoin(
	ctx context.Context,
	logger runtime.Logger,
	db *sql.DB,
	nk runtime.NakamaModule,
	dispatcher runtime.MatchDispatcher,
	tick int64,
	state interface{},
	presences []runtime.Presence) interface{} {
	s := state.(*entity.SlotsMatchState)
	logger.Info("match join, state=%v, presences=%v", s, presences)

	m.processor.ProcessPresencesJoin(ctx,
		logger,
		nk,
		db,
		dispatcher,
		presences,
		s,
	)
	return s
}

func (m *MatchHandler) MatchLeave(ctx context.Context,
	logger runtime.Logger,
	db *sql.DB,
	nk runtime.NakamaModule,
	dispatcher runtime.MatchDispatcher,
	tick int64,
	state interface{},
	presences []runtime.Presence,
) interface{} {
	s := state.(*entity.SlotsMatchState)

	logger.Info("match leave, state=%v, presences=%v", s, presences)

	if m.machine.IsPlayingState() || m.machine.IsReward() {
		m.processor.ProcessPresencesLeavePending(ctx,
			logger,
			nk,
			dispatcher,
			presences,
			s,
		)
		return s
	}
	m.processor.ProcessPresencesLeave(ctx,
		logger,
		nk, db,
		dispatcher,
		presences,
		s,
	)
	return s
}

func (m *MatchHandler) MatchLoop(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, messages []runtime.MatchData) interface{} {
	s := state.(*entity.SlotsMatchState)

	err := m.machine.FireProcessEvent(lib.GetContextWithProcessorPackager(
		lib.NewProcessorPackage(
			s, m.processor,
			logger,
			nk,
			db,
			dispatcher,
			messages,
			ctx),
	))
	if err == presenter.ErrGameFinish {
		logger.Info("match need finish")

		return nil
	}

	return s
}

func (m *MatchHandler) MatchTerminate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, graceSeconds int) interface{} {
	logger.Info("match terminate, state=%v")
	return state
}
