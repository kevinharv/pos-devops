CR := docker
AppName := pos-server

run: build
	$(CR) compose up -d
	$(CR) compose logs -f

build:
	$(CR) compose build

stop:
	$(CR) compose down

restart: stop run

clean:
	$(CR) compose down -v
	$(CR) system prune
