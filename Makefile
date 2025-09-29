clear:
	clear && rm -rf dist/*.csv

run: clear
	go run main.go
	