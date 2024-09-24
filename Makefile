.PHONY: run-singlenode
run-singlenode: down
	docker compose up

.PHONY: run-multinode
run-multinode: down
	docker compose up --build --scale api=3

.PHONY: down
down:
	docker compose down
