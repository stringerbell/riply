up:
	docker compose up -d
down:
	docker compose down

test:
	go test ./... -v

integ:
	INTEG=true go test ./... -v
