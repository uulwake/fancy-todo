all: hello

hello: 
	@printf "Welcome to the Fancy Todo written in 5 different languages.\n\
	They are Typescript, Go, Java, Rust, and Elixir.\n\
	Enjoy.\n"

# tear down
down:
	docker compose down && cd ./typescript && docker compose down

down\:db:
	docker compose down

down\:ts:
	cd ./typescript && docker compose down

# database
db\:migrate:
	dbmate --wait up -v

db\:up:
	docker compose up -d && make wait && make db:migrate

db\:down:
	docker compose down

db\:up\:pg:
	docker compose up -d pg && make db:migrate

db\:up\:es:
	docker compose up -d es && make wait

wait:
	chmod +x ./wait-for-it.sh && ./wait-for-it.sh es:9200 -- echo "ElasticSearch database is up"

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