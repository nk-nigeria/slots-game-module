module github.com/ciaolink-game-platform/cgb-slots-game-module

go 1.18

require (
	github.com/bwmarrin/snowflake v0.3.0
	github.com/heroiclabs/nakama-common v1.22.0
)

require (
	github.com/ciaolink-game-platform/cgp-common v0.0.0-20230704085817-ae1a5627d091
	github.com/qmuntal/stateless v1.5.3
	github.com/stretchr/testify v1.7.2
	github.com/wk8/go-ordered-map/v2 v2.1.3
	google.golang.org/grpc v1.42.0
	google.golang.org/protobuf v1.27.1
)

replace github.com/ciaolink-game-platform/cgp-common => ./cgp-common

require (
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	google.golang.org/genproto v0.0.0-20211118181313-81c1377c94b1 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
