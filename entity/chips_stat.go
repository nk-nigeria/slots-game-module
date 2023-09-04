package entity

import (
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

type chipStat struct {
	chipsWin      map[pb.SiXiangGame]int64
	lineWin       map[pb.SiXiangGame]int64
	totalChipsWin map[pb.SiXiangGame]int64
	totalLineWin  map[pb.SiXiangGame]int64
}

func NewChipStat() *chipStat {
	s := chipStat{
		chipsWin:      make(map[pb.SiXiangGame]int64),
		lineWin:       make(map[pb.SiXiangGame]int64),
		totalChipsWin: make(map[pb.SiXiangGame]int64),
		totalLineWin:  make(map[pb.SiXiangGame]int64),
	}
	return &s
}

func (s *chipStat) ResetAll() {
	s.chipsWin = make(map[pb.SiXiangGame]int64)
	s.lineWin = make(map[pb.SiXiangGame]int64)
	s.totalChipsWin = make(map[pb.SiXiangGame]int64)
	s.totalLineWin = make(map[pb.SiXiangGame]int64)
}

func (s *chipStat) Reset(game pb.SiXiangGame) {
	s.ResetChipWin(game)
	s.ResetLineWin(game)
	s.ResetTotalChipWin(game)
	s.ResetTotalLineWin(game)
}

func (s *chipStat) AddChipWin(game pb.SiXiangGame, chips int64) {
	{
		v := s.chipsWin[game]
		v += chips
		s.chipsWin[game] = v
	}
	{
		v := s.totalChipsWin[game]
		v += chips
		s.totalChipsWin[game] = v
	}
}

func (s *chipStat) AddLineWin(game pb.SiXiangGame, lineWin int64) {
	{
		v := s.lineWin[game]
		v += lineWin
		s.lineWin[game] = v
	}
	{
		v := s.totalLineWin[game]
		v += lineWin
		s.totalLineWin[game] = v
	}
}

func (s *chipStat) ResetChipWin(game pb.SiXiangGame) {
	s.chipsWin[game] = 0
}

func (s *chipStat) ResetLineWin(game pb.SiXiangGame) {
	s.lineWin[game] = 0
}

func (s *chipStat) ChipWin(game pb.SiXiangGame) int64 {
	return s.chipsWin[game]
}

func (s *chipStat) LineWin(game pb.SiXiangGame) int64 {
	return s.lineWin[game]
}

func (s *chipStat) TotalChipWin(game pb.SiXiangGame) int64 {
	return s.totalChipsWin[game]
}

func (s *chipStat) TotalLineWin(game pb.SiXiangGame) int64 {
	return s.totalLineWin[game]
}

func (s *chipStat) ResetTotalChipWin(game pb.SiXiangGame) {
	s.totalChipsWin[game] = 0
}

func (s *chipStat) ResetTotalLineWin(game pb.SiXiangGame) {
	s.totalLineWin[game] = 0
}
