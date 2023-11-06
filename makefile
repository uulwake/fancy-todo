# database
db\:up:
	docker-compose up -d && dbmate up

db\:down:
	docker-compose down

db\:up\:pg:
	docker-compose up -d pg && dbmate up

db\:up\:es:
	docker-compose up -d es && dbmate up

# node ts
node\:build:
	cd ./ts && npm run build

node\:dev:
	cd ./ts && npm run dev