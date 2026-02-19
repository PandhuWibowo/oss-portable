.PHONY: dev

dev:
	@trap 'kill 0' INT; \
	  (cd backend && go run main.go) & \
	  (until nc -z localhost 8080 2>/dev/null; do sleep 0.3; done && cd frontend && bun run dev) & \
	  wait
