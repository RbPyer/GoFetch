EXEC_PATH=./bin
OBJ_PATH=./cmd/main.go
ADD=git add .
REPO=git@github.com:RbPyer/GoFetch.git develop
PUSH=git push $(REPO)

all:
	go build  -ldflags "-s -w" -o $(EXEC_PATH)/GoFetch $(OBJ_PATH)

run: all
	$(EXEC_PATH)/GoFetch

bench:
	go test -bench=.