DIST_DIR = .dist

$(DIST_DIR):
	mkdir -p $(DIST_DIR)

$(DIST_DIR)/beremiz.exe: cmd/beremiz/main.go $(wildcard internal/**/*.go) | $(DIST_DIR)
	go build -o $@ cmd/beremiz/main.go

all: .dist/beremiz.exe

clean:
	rm -rf $(DIST_DIR)

.PHONY: all clean
