# Diretório de saída
DIST_DIR = .dist

# Cria o diretório se não existir
$(DIST_DIR):
	mkdir -p $(DIST_DIR)

# Sua regra: o target depende do main.go e de todos os .go em internal/**
$(DIST_DIR)/beremiz.exe: cmd/beremiz/main.go $(wildcard internal/**/*.go) | $(DIST_DIR)
	go build -o $@ cmd/beremiz/main.go

all: .dist/beremiz.exe

clean:
	rm -rf $(DIST_DIR)

.PHONY: all clean
