PROJECT_NAME=github.com/nakamaFramework/cgb-slots-game-module
APP_NAME=slots-game.so
APP_PATH=$(PWD)

update-submodule-dev:
	# git checkout develop && git pull
	git submodule update --init
	git submodule update --remote
	cd ./cgp-common && git checkout develop && git pull && cd ..
	go get github.com/nakamaFramework/cgp-common@develop
update-submodule-stg:
	git checkout staging && git pull
	git submodule update --init
	git submodule update --remote
	cd ./cgp-common && git checkout staging && git pull && cd ..
	go get github.com/nakamaFramework/cgp-common@staging

build:
	./sync_pkg_3.11.sh
	go mod tidy && 	go mod vendor
	docker run --rm -w "/app" -v "${APP_PATH}:/app" heroiclabs/nakama-pluginbuilder:3.11.0 build -buildvcs=false --trimpath --buildmode=plugin -o ./bin/${APP_NAME}

syncdev:
	rsync -aurv --delete ./bin/${APP_NAME} root@cgpdev:/root/cgp-server-dev/dist/data/modules/bin/
	ssh root@cgpdev 'cd /root/cgp-server-dev && docker restart nakama_dev'
syncstg:
	rsync -aurv --delete ./bin/${APP_NAME} root@cgpdev:/root/cgp-server/dist/data/modules/bin
	ssh root@cgpdev 'cd /root/cgp-server && docker restart nakama'

dev: update-submodule-dev build
stg: update-submodule-stg build

proto:
	protoc -I ./ --go_out=$(pwd)/proto  ./proto/chinese_poker_game_api.proto

local:
	# git submodule update --init
	# git submodule update --remote
	# go get github.com/nakamaFramework/cgp-common@main
	./sync_pkg_3.11.sh
	go mod tidy
	go mod vendor
	rm ./bin/* || true
	go build -buildvcs=false --trimpath --mod=vendor --buildmode=plugin -o ./bin/${APP_NAME}
