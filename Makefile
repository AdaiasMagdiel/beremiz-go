.PHONY: all clean

.dist/beremiz.exe: cmd/beremiz/main.go
	go build -o .dist/beremiz.exe cmd/beremiz/main.go

all: .dist/beremiz.exe

clean:
	rm -rf .dist
	mkdir .dist
