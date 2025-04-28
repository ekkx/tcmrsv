.PHONY: up down shell

up:
	docker compose -f ./build/compose.yaml -p tcmrsv up -d

down:
	docker compose -f ./build/compose.yaml -p tcmrsv down

shell:
	docker compose -f ./build/compose.yaml -p tcmrsv exec -it go bash
