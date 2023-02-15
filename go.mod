module github.com/ciaolink-game-platform/cgb-slots-game-module

go 1.19

require (
	github.com/bwmarrin/snowflake v0.3.0
	github.com/heroiclabs/nakama-common v1.26.0
)

replace github.com/ciaolink-game-platform/cgp-common => ./cgp-common

require (
	github.com/ciaolink-game-platform/cgp-common v0.0.0-00010101000000-000000000000
	github.com/qmuntal/stateless v1.6.2
	google.golang.org/grpc v1.42.0
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.0.0-20220630215102-69896b714898 // indirect
	golang.org/x/sys v0.0.0-20220704084225-05e143d24a9e // indirect
	google.golang.org/genproto v0.0.0-20211118181313-81c1377c94b1 // indirect
)
