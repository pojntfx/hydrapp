# Public variables
DESTDIR ?=
PREFIX ?= /usr/local
OUTPUT_DIR ?= out
DST ?=

# Private variables
obj = ./hydrapp-builder/cmd/hydrapp
all: $(addprefix build/,$(obj))

# Build
build: $(addprefix build/,$(obj))
$(addprefix build/,$(obj)):
ifdef DST
	go build -o $(DST) $(subst build/,,$@)
else
	go build -o $(OUTPUT_DIR)/$(subst build/,,$@) $(subst build/,,$@)
endif

# Install
install: $(addprefix install/,$(obj))
$(addprefix install/,$(obj)):
	install -D -m 0755 $(OUTPUT_DIR)/$(subst install/,,$@) $(DESTDIR)$(PREFIX)/bin/$(shell basename $(subst install/,,$@))

# Uninstall
uninstall: $(addprefix uninstall/,$(obj))
$(addprefix uninstall/,$(obj)):
	rm $(DESTDIR)$(PREFIX)/bin/$(shell basename $(subst uninstall/,,$@))

# Run
$(addprefix run/,$(obj)):
	$(subst run/,,$@) $(ARGS)

# Test
test: $(addprefix test/,$(obj))
$(addprefix test/,$(obj)):
	go test -timeout 3600s -parallel $(shell nproc) $(shell dirname $(dir $(subst test/,,$@)))/...

# Benchmark
benchmark: $(addprefix benchmark/,$(obj))
$(addprefix benchmark/,$(obj)):
	go test -timeout 3600s -bench=$(shell dirname $(dir $(subst benchmark/,,$@)))/... $(shell dirname $(dir $(subst benchmark/,,$@)))/...

# Clean
clean:
	rm -rf out

# Dependencies
depend:
	true
