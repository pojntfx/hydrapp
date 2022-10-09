# Public variables
PREFIX ?= ghcr.io/pojntfx/hydrapp-build-

# Private variables
obj = apk deb dmg flatpak msi rpm
all: $(addprefix build/,$(obj))

# Build
build: $(addprefix build/,$(obj))
$(addprefix build/,$(obj)):
	docker build -t ghcr.io/pojntfx/hydrapp-build-$(subst build/,,$@) ./pkg/builders/$(subst build/,,$@)/
