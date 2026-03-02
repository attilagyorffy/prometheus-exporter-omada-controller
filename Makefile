.PHONY: build run query

build:
	go build -o omada-controller-exporter .

run:
	OMADA_URL=https://10.0.0.3:30077 \
	OMADA_USER=$$(op read "op://Private/Omada/username") \
	OMADA_PASS=$$(op read "op://Private/Omada/password") \
	OMADA_INSECURE=true \
	go run .

query:
	curl -s http://localhost:6779/metrics
