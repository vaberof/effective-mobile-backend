PROJECT_NAME=car-catalog
APP_LOCAL_NAME=web-backend

DOCKER_LOCAL_IMAGE_NAME=$(PROJECT_NAME)/$(APP_LOCAL_NAME)

WORK_DIR_LINUX=./cmd/carcatalog
CONFIG_DIR_LINUX=./cmd/carcatalog/config

POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=admin
POSTGRES_DATABASE=car_catalog

DB_URL="postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=disable"

MIGRATIONS_PATH=migrations/

docker.run: docker.build
	docker compose -f cmd/carcatalog/docker-compose.yaml up -d

docker.build: build.linux
	docker build -t $(DOCKER_LOCAL_IMAGE_NAME) -f cmd/carcatalog/Dockerfile .

run.linux: build.linux
	go run $(WORK_DIR_LINUX)/*.go \
		-config.files $(CONFIG_DIR_LINUX)/application.yaml \
		-env.vars.file $(CONFIG_DIR_LINUX)/application.env \

build.linux: build.linux.clean
	mkdir -p $(WORK_DIR_LINUX)/build
	go build -o $(WORK_DIR_LINUX)/build/main $(WORK_DIR_LINUX)/*.go
	cp -R $(CONFIG_DIR_LINUX)/* $(WORK_DIR_LINUX)/build

build.linux.local: build.linux.clean
	mkdir -p $(WORK_DIR_LINUX)/build
	go build -o $(WORK_DIR_LINUX)/build/main $(WORK_DIR_LINUX)/*.go
	cp -R $(CONFIG_DIR_LINUX)/* $(WORK_DIR_LINUX)/build
	@echo "build.local: OK"

build.linux.clean:
	rm -rf $(WORK_DIR_LINUX)/build

migrate.up:
	migrate -path $(MIGRATIONS_PATH) -database $(DB_URL) up

migrate.down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down

swagger.gen:
	swag init --parseDependency --parseInternal -g ./cmd/carcatalog/main.go -o ./cmd/carcatalog/docs
