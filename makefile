# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

# Deploy First Mentality


run:
	go run apis/services/sales/main.go | go run apis/tooling/logfmt/main.go