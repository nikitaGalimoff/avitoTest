COMPOSE_FILE=docker-compose.yml


.PHONY: up down test e2e

up:
	docker-compose -f $(COMPOSE_FILE) up -d

down:
	docker-compose -f $(COMPOSE_FILE) down


e2e:
	go test ./e2e -v