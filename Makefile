.PHONY: backend frontend dev

backend:
	go -C server run ./cmd

frontend:
	go -C frontend run ./cmd

dev:
	$(MAKE) backend & \
	$(MAKE) frontend & \
	wait
