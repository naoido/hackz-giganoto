include .env

DIRS := $(wildcard microservices/*/)
TARGETS := $(notdir $(patsubst %/,%,$(DIRS)))

.PHONY: $(TARGETS) $(addprefix run-, $(TARGETS)) $(addprefix gen-, $(TARGETS)) build
$(addprefix gen-, $(TARGETS)): gen-%:
	@echo "Running goa gen: $*"
	cd microservices/$* && goa gen object-t.com/hackz-giganoto/microservices/$*/design

$(addprefix example-, $(TARGETS)): example-%:
	@echo "Running goa example: $*"
	cd microservices/$* && goa example object-t.com/hackz-giganoto/microservices/$*/design

$(addprefix run-, $(TARGETS)): run-%:
	@echo "Running server: $*"
	cd microservices/$* && go build -o server ./cmd/$* && ./server

build:
	docker compose build auth
	docker compose build bff
	docker compose build profile
	docker compose build chat

start:
	docker compose up