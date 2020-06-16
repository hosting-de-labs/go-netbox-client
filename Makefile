.PHONY: test

test:
	docker-compose -f test/docker/docker-compose.yml up -d
	sleep 5
	go test ./...
	docker-compose -f test/docker/docker-compose.yml down
