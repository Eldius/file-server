
start:
	go run ./cmd/file-server start --debug --log-format=json

get_root:
	curl -i http://localhost:8080/

get_text:
	curl -i http://localhost:8080/Makefile

get_img:
	curl -i http://localhost:8080/internal/sample_data/mic-drop.png
