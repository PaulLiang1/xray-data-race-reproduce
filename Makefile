SHELL := /bin/bash
TOP_DIR := $(shell pwd)

.PHONY: run-data-race
run-data-race:
	docker-compose -f $(TOP_DIR)/data_race/docker-compose.yml up --build --force-recreate

.PHONY: run-safe
run-safe:
	docker-compose -f $(TOP_DIR)/safe/docker-compose.yml up --build --force-recreate
