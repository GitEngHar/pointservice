.PHONY: tidy build run up down

APP_NAME=point-service
IMAGE_NAME=point-service
BASE_IMAGE_NAME=point-service-base
IMAGE_TAG=1.0
BASE_IMAGE_TAG=1.0

tidy:
	go mod tidy

set-build-context:
	docker context use default

build-app:
	docker build --target final -t $(IMAGE_NAME):$(IMAGE_TAG) -f DockerfileApp .

build-base:
	docker build --load -t $(BASE_IMAGE_NAME):$(BASE_IMAGE_TAG) -f DockerfileBase .

run-debug:
	docker compose up --build

run: set-build-context build-base
	docker compose -f docker-compose.yaml up --build

run-middleware:
	docker compose -f docker-compose.middleware.yml up -d

# keploy で record 実行
run-keploy: run-middleware build-base build-app
	keploy record -c "docker run -p 1323:1323 --name point-app --network pointservice_keploy-network -e POINT_MYSQL_HOST=mysql -e POINT_MYSQL_DATABASE=appdb -e POINT_MYSQL_PORT=3306 -e POINT_MYSQL_USER=appuser -e POINT_MYSQL_PASSWORD=apppass -e POINT_MYSQL_MAX_OPEN_CONNECTIONS=2 -e POINT_MYSQL_MAX_IDLE_CONNECTIONS=1 point-service:1.0"

test-keploy: run-middleware build-base build-app
	keploy test -c "docker run -p 1323:1323 --name point-app --network pointservice_keploy-network -e POINT_MYSQL_HOST=mysql -e POINT_MYSQL_DATABASE=appdb -e POINT_MYSQL_PORT=3306 -e POINT_MYSQL_USER=appuser -e POINT_MYSQL_PASSWORD=apppass -e POINT_MYSQL_MAX_OPEN_CONNECTIONS=2 -e POINT_MYSQL_MAX_IDLE_CONNECTIONS=1 -e ENVIRONMENT=keploy point-service:1.0" --delay 15

down:
	docker compose down -v || true

# OpenAPI
openapi-lint:
	npx --package @redocly/cli@latest redocly lint openapi/openapi.yaml

openapi-docs:
	npx --package @redocly/cli@latest redocly build-docs openapi/openapi.yaml --output html/openapi/index.html