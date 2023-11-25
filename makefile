default: hello

hello: 
	@printf "Welcome to the Fancy Todo written in 5 different languages.\n\
	They are Typescript, Go, Java, Rust, and Elixir.\n\
	Enjoy.\n"

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
	@printf "waiting ElasticSeach...\n" && chmod +x ./wait-for-it.sh && ./wait-for-it.sh es:9200 -- echo "waiting is done"

# ts
build\:ts:
	docker compose build ts

run\:ts:
	make db:up && docker compose up -d ts

test\:ts:
	cd ./tests && make ts:test

test\:unit\:ts:
	cd ./typescript && npm run test:coverage

# go
build\:go:
	docker compose build go

run\:go:
	make db:up && docker compose up -d go

test\:go:
	cd ./tests && make go:test