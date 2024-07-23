
test:
	go test ./... -v

lint:
	golangci-lint run

run:
	docker-compose up -d	
.PHONY: run	

