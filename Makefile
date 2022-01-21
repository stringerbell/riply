up:
	docker compose up -d
down:
	docker compose down

test:
	go test ./... -v

integ:
	INTEG=true go test ./... -v

# clean up any containers that might be laying around
clean:
	docker stop `docker ps | grep range | tr -d '' | cut -d ' ' -f 1`
