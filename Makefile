clear:
	clear && rm -rf dist/*.csv

run: clear
	go run main.go

asg_list:
	cd cmd/asg_list && go run main.go
