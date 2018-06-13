PKG := github.com/guggero/docker-wallet-control

DEP_PKG := github.com/golang/dep/cmd/dep

GO_BIN := ${GOPATH}/bin
DEP_BIN := $(GO_BIN)/dep

HAVE_DEP := $(shell command -v $(DEP_BIN) 2> /dev/null)

COMMIT := $(shell git rev-parse HEAD)
LDFLAGS := -ldflags "-X main.Commit=$(COMMIT)"

GOBUILD := go build -v
GOINSTALL := go install -v
GOTEST := go test -v

CGO_STATUS_QUO := ${CGO_ENABLED}

RM := rm -f
CP := cp
MAKE := make
XARGS := xargs -L 1

GREEN := "\\033[0;32m"
NC := "\\033[0m"
define print
	echo $(GREEN)$1$(NC)
endef

default: scratch

all: scratch unit install

# ============
# DEPENDENCIES
# ============

$(DEP_BIN):
	@$(call print, "Fetching dep.")
	go get -u $(DEP_PKG)

dep: $(DEP_BIN)
	@$(call print, "Compiling dependencies.")
	dep ensure -v

# ============
# INSTALLATION
# ============

build:
	@$(call print, "Building debug.")
	$(GOBUILD) $(LDFLAGS) $(PKG)

install:
	@$(call print, "Installing.")
	go install -v $(LDFLAGS) $(PKG)

scratch: dep build

clean:
	@$(call print, "Cleaning source.$(NC)")
	$(RM) ./docker-wallet-control
	$(RM) -r ./vendor

# =======
# TESTING
# =======

unit:
	@$(call print, "Running unit tests.")
	$(GOTEST) .
	$(GOTEST) ./util

# ======
# TRAVIS
# ======

travis: scratch unit

.PHONY: all \
	default \
	dep \
	build \
	install \
	scratch \
	unit \
	clean
