# Public variables
PREFIX ?= ghcr.io/pojntfx/hydrapp-build-
TAG ?= latest

# Private variables
obj = apk deb dmg flatpak msi rpm
all: $(addprefix build/,$(obj))

# Build
build: $(addprefix build/,$(obj))
$(addprefix build/,$(obj)):
	docker build -t $(PREFIX)$(subst build/,,$@):$(TAG) ./pkg/builders/$(subst build/,,$@)/
