package entity

import (
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

type TarzanMatchState struct {
	lib.MatchState
	CurrentSiXiangGame pb.SiXiangGame
	// next game in loop
	NextSiXiangGame pb.SiXiangGame
}
