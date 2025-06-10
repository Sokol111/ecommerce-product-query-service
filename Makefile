include .env
DOMAIN = $(PROJECT_NAME).localhost

.PHONY: generate-mocks build-docker-image start-docker-compose stop-docker-compose update-dependencies test init-git show-container-logs ensure-network

generate-mocks:
	mockery

add-host:
	@echo "Adding domain to /etc/hosts..."
	@if ! grep -q "$(DOMAIN)" /etc/hosts; then \
		echo "127.0.0.1 $(DOMAIN)" | sudo tee -a /etc/hosts > /dev/null; \
		echo "Added: $(DOMAIN)"; \
	else \
		echo "Already exists: $(DOMAIN)"; \
	fi

ensure-network:
	docker network inspect shared-network > /dev/null 2>&1 || docker network create shared-network

build-docker-image:
	docker build -t sokol111/$(PROJECT_NAME):latest .

start-docker-compose: ensure-network stop-docker-compose
	docker compose up -d

stop-docker-compose:
	docker compose down

show-container-logs:
	docker compose logs -f

update-dependencies:
	go get -u ./...

test:
	go test ./... -v -cover

init-git:
	git config user.name "Sokol111"
	git config user.email "igorsokol111@gmail.com"
	git config commit.gpgSign false