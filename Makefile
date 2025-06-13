.PHONY:
.SILENT:
.DEFAULT_GOAL := run


build:
	GOOS=linux GOARCH=amd64 go build -o server