test:
	go test ./... -v

integ:
	INTEG=true go test ./... -v
