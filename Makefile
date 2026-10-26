.PHONY: help status check hint ticket reset run test

help:
	@echo "Дежурный по проду: инциденты в Go-сервисе справочника брендов."
	@echo ""
	@echo "  make status     доска инцидентов"
	@echo "  make level-1    начать инцидент 1 (так же level-2, level-3, ...)"
	@echo "  make ticket     показать тикет текущего инцидента"
	@echo "  make check      проверить, починил ли"
	@echo "  make hint       следующая подсказка"
	@echo "  make reset      начать текущий инцидент заново"
	@echo "  make run        запустить сервис на :8080"
	@echo "  make test       только твои тесты"

level-%:
	@./scripts/level.sh $*

status:
	@./scripts/status.sh

check:
	@./scripts/check.sh

hint:
	@./scripts/hint.sh

ticket:
	@./scripts/ticket.sh

reset:
	@./scripts/level.sh $$(cat .gopher/current)

run:
	go run ./cmd/api

test:
	go test ./...
