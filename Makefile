.PHONY: generate generate-mocks build-docker-image start-docker-compose stop-docker-compose update-dependencies test

generate:
	go generate ./...

generate-mocks:
	mockery

build-docker-image:
	docker build -t sokol111/ecommerce-product-query-service:latest .

start-docker-compose:
	docker-compose up -d

stop-docker-compose:
	docker-compose down

update-dependencies:
	go get -u ./...

test:
	go test ./... -v -cover
