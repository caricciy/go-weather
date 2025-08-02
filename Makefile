run:
	@docker compose up --build --force-recreate --no-deps -d

test:
	@go test -count=1 ./...

.PHONY: run test