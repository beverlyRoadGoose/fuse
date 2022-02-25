include project.mk

test-unit:
	go test ./... -coverprofile=coverage.txt -covermode=atomic
