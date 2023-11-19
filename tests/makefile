default: hello

hello: 
	@printf "Integration Test"

# tear down
down:
	docker compose down

# database
db\:up:
	make db:up:pg && make db:up:es 

db\:down:
	docker compose down

pg\:migrate:
	dbmate --wait up -v

db\:up\:pg:
	docker compose up -d pg && make pg:migrate

db\:up\:es:
	docker compose up -d es && make wait:es

wait\:es:
	@printf "waiting ElasticSeach...\n" && chmod +x ../wait-for-it.sh && ../wait-for-it.sh es:9200 -- echo "waiting is done"

# ts
ts\:build:
	docker compose build ts

ts\:run:
	make db:up && docker compose up -d ts

ts\:test:
	make ts:run && npm run test && make down && docker volume rm tests_pg-data tests_elastic-data  

# go
go\:build:
	docker compose build go

go\:run:
	make db:up && docker compose up -d go

go\:test:
	make go:run && npm run test && make down && docker volume rm tests_pg-data tests_elastic-data 