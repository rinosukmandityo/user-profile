COMPOSE_RUN=docker-compose --no-ansi -f docker/docker-compose.yaml
run:
	$(COMPOSE_RUN) down
	$(COMPOSE_RUN) up -d --force-recreate --always-recreate-deps

COMPOSE_TEST=docker-compose --no-ansi -f docker/docker-compose-test.yaml
RUN_TEST=run -e CGO_ENABLED=0 -e PWD=$(CURDIR) \
	-e GO111MODULE=on
test:
	$(COMPOSE_TEST) down --remove-orphan
	$(COMPOSE_TEST) up -d
	$(COMPOSE_TEST) $(RUN_TEST) test sh coverage.sh
	$(COMPOSE_TEST) stop

db:
	$(COMPOSE_RUN) down
	$(COMPOSE_RUN) up -d --force-recreate --always-recreate-deps mysql

COMPOSE_RUN_SECRET=docker-compose --no-ansi -f docker/docker-compose-secret.yaml
run-secret:
	$(COMPOSE_RUN_SECRET) down
	$(COMPOSE_RUN_SECRET) up -d --force-recreate --always-recreate-deps