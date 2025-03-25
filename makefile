SRC = examples/example1.go main.go
DST = examples/example1_json.go

all: $(DST)

$(DST): $(SRC)
	go run . -dir examples -name Example1
