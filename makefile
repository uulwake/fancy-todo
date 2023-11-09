all: hello

hello: 
	@printf "Welcome to the Fancy Todo written in 5 different languages.\n\
	They are Typescript, Go, Java, Rust, and Elixir.\n\
	Enjoy.\n"

# tear down
down:
	docker compose down

# database
db\:migrate:
	dbmate --wait up -v

db\:up:
	docker compose up -d pg es && make wait:es && make db:migrate 

db\:down:
	docker compose down

db\:up\:pg:
	docker compose up -d pg && make db:migrate

db\:up\:es:
	docker compose up -d es && make wait:es

wait\:es:
	@printf "waiting ElasticSeach...\n" && chmod +x ./wait-for-it.sh && ./wait-for-it.sh es:9200 -- echo "waiting is done"

# typescript ts
ts\:build:
	docker build . -f ts.Dockerfile -t ts

ts\:up:
	docker compose up -d ts

ts\:run:
	make db:up && make ts:build && make ts:up