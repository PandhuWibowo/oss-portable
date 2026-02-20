.PHONY: dev build

dev:
	@trap 'kill 0' INT; \
	  (cd server && go run main.go) & \
	  (until nc -z localhost 8080 2>/dev/null; do sleep 0.3; done && cd web && bun run dev) & \
	  wait

build:
	@cd server && go build -o ../bin/server .
	@cd web && bun run build
