PROJECT_NAME=github.com/nk-nigeria/slots-game-module
APP_NAME=slots_plugin.so
APP_PATH=$(PWD)
NAKAMA_VER=3.27.0

build:
	go mod vendor
	docker run --rm -w "/app" -v "${APP_PATH}:/app" "heroiclabs/nakama-pluginbuilder:${NAKAMA_VER}" build -buildvcs=false --trimpath --buildmode=plugin -o ./bin/${APP_NAME} . && cp ./bin/${APP_NAME} ../bin/
	
sync:
	rsync -aurv --delete ./bin/${APP_NAME} root@cgpdev:/root/cgp-server/dev/data/modules/
	# ssh root@cgpdev 'cd /root/cgp-server && docker restart nakama'

bsync: build sync
proto:
	protoc -I ./ --go_out=$(pwd)/proto  ./proto/chinese_poker_game_api.proto

local:
	git submodule update --init
	git submodule update --remote
	go get github.com/nk-nigeria/cgp-common@main
	go mod tidy
	go mod vendor
	rm ./bin/* || true
	go build -buildvcs=false --trimpath --mod=vendor --buildmode=plugin -o ./bin/${APP_NAME}