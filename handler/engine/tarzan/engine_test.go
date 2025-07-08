package tarzan

import (
	"reflect"
	"testing"

	"github.com/nk-nigeria/cgp-common/lib"
	pb "github.com/nk-nigeria/cgp-common/proto"
)

func TestNewEngine(t *testing.T) {
	tests := []struct {
		name string
		want lib.Engine
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEngine(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEngine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tarzanEngine_NewGame(t *testing.T) {
	type fields struct {
		engines map[pb.SiXiangGame]lib.Engine
	}
	type args struct {
		matchState interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &tarzanEngine{
				engines: tt.fields.engines,
			}
			got, err := e.NewGame(tt.args.matchState)
			if (err != nil) != tt.wantErr {
				t.Errorf("tarzanEngine.NewGame() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tarzanEngine.NewGame() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tarzanEngine_Process(t *testing.T) {
	type fields struct {
		engines map[pb.SiXiangGame]lib.Engine
	}
	type args struct {
		matchState interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &tarzanEngine{
				engines: tt.fields.engines,
			}
			got, err := e.Process(tt.args.matchState)
			if (err != nil) != tt.wantErr {
				t.Errorf("tarzanEngine.Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tarzanEngine.Process() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tarzanEngine_Random(t *testing.T) {
	type fields struct {
		engines map[pb.SiXiangGame]lib.Engine
	}
	type args struct {
		min int
		max int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &tarzanEngine{
				engines: tt.fields.engines,
			}
			if got := e.Random(tt.args.min, tt.args.max); got != tt.want {
				t.Errorf("tarzanEngine.Random() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tarzanEngine_Finish(t *testing.T) {
	type fields struct {
		engines map[pb.SiXiangGame]lib.Engine
	}
	type args struct {
		matchState interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &tarzanEngine{
				engines: tt.fields.engines,
			}
			got, err := e.Finish(tt.args.matchState)
			if (err != nil) != tt.wantErr {
				t.Errorf("tarzanEngine.Finish() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tarzanEngine.Finish() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tarzanEngine_transformLineWinToBigWin(t *testing.T) {
	type fields struct {
		engines map[pb.SiXiangGame]lib.Engine
	}
	type args struct {
		lineWin int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   pb.BigWin
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &tarzanEngine{
				engines: tt.fields.engines,
			}
			if got := e.transformLineWinToBigWin(tt.args.lineWin); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tarzanEngine.transformLineWinToBigWin() = %v, want %v", got, tt.want)
			}
		})
	}
}
