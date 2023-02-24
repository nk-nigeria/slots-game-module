PROJECT_NAME=github.com/ciaolink-game-platform/cgb-slots-game-module
APP_NAME=slots-game.so
APP_PATH=$(PWD)

build:
	git submodule update --init
	git submodule update --remote
	go get github.com/ciaolink-game-platform/cgp-common@main
	go mod tidy
	go mod vendor
	docker run --rm -w "/app" -v "${APP_PATH}:/app" heroiclabs/nakama-pluginbuilder:3.11.0 build -buildvcs=false --trimpath --buildmode=plugin -o ./bin/${APP_NAME}
	
sync:
	rsync -aurv --delete ./bin/${APP_NAME} root@cgpdev:/root/cgp-server/dev/data/modules/
	# ssh root@cgpdev 'cd /root/cgp-server && docker restart nakama'

bsync: build sync
proto:
	protoc -I ./ --go_out=$(pwd)/proto  ./proto/chinese_poker_game_api.proto

local:
	go build --trimpath --mod=vendor --buildmode=plugin -o ./bin/slot-game.so