tidy:
	go fmt ./...
	go mod tidy -v

easyjs:
	easyjson -no_std_marshalers -all internal/entity

vendor:
	go mod tidy -v
	go mod vendor

run:
	docker compose build --no-cache
	docker compose up