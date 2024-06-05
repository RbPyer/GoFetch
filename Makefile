EXEC_PATH=./packages/releases
OBJ_PATH=./cmd/main.go
ADD=git add .
REPO=git@github.com:RbPyer/GoFetch.git develop
PUSH=git push $(REPO)

all:
	go build -o $(EXEC_PATH)/GoFetch $(OBJ_PATH)

