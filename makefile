# service
down:
	docker compose down && cd ./typescript && docker compose down

# database
db\:migrate:
	dbmate --wait up -v

db\:up:
	docker compose up -d && make db:migrate

db\:down:
	docker compose down

db\:up\:pg:
	docker compose up -d pg && make db:migrate

db\:up\:es:
	docker compose up -d es

# typescript ts
ts\:install:
	cd ./typescript && npm install
	
ts\:build:
	cd ./typescript && npm run build

ts\:dev:
	cd ./typescript && npm run dev

ts\:run\:skip-es:
	make db:up:pg && cd ./typescript && npm run start

ts\:run:
	make db:up && cd ./typescript && npm run start