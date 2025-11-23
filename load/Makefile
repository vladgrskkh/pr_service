include .env_db

# =========================================================================================== #
# LOAD TESTING
# =========================================================================================== #

.PHONY: load/up
load/up:
	docker-compose -f docker-compose.test.yml up --build -d

.PHONY: load/down
load/down:
	docker-compose -f docker-compose.test.yml down -v

.PHONY: load/seed
load/seed:
	cat synt.sql | docker compose -f docker-compose.test.yml exec -T db psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

.PHONY: load/start
load/start:
	vegeta attack -duration=100s -rate=30 -targets=target.list -output=attack-30.bin
	vegeta plot -title=Attack%30Results attack-30.bin > results30.html
	vegeta report -type=json attack-30.bin > attack-30.json