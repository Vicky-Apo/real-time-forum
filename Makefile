.PHONY: backend frontend dev

backend:
	go -C server run ./cmd

frontend:
	go -C client run main.go

dev:
	$(MAKE) backend & \
	$(MAKE) frontend & \
	wait
