# Public variables
DESTDIR ?=
PREFIX ?= /usr/local
OUTPUT_DIR ?= out
DST ?=

# Private variables
obj = hydrapp-cli hydrapp-example-rest hydrapp-example-forms hydrapp-example-dudirekta
all: $(addprefix build/,$(obj))

# Build
build: $(addprefix build/,$(obj))
$(addprefix build/,$(obj)):
ifdef DST
	go build -o $(DST) ./$(subst build/,,$@)
else
	go build -o $(OUTPUT_DIR)/$(subst build/,,$@) ./$(subst build/,,$@)
endif

# Install
install: $(addprefix install/,$(obj))
$(addprefix install/,$(obj)):
	install -D -m 0755 $(OUTPUT_DIR)/$(subst install/,,$@) $(DESTDIR)$(PREFIX)/bin/$(basename $(subst install/,,$@))

# Uninstall
uninstall: $(addprefix uninstall/,$(obj))
$(addprefix uninstall/,$(obj)):
	rm -f $(DESTDIR)$(PREFIX)/bin/$(basename $(subst uninstall/,,$@))

# Run
$(addprefix run/,$(obj)):
	$(subst run/,,$@) $(ARGS)

# Test
test: $(addprefix test/,$(obj))
$(addprefix test/,$(obj)):
	go test -timeout 3600s -parallel $(shell nproc) ./$(shell echo $(subst test/,,$@) | cut -d / -f1)/...

# Benchmark
benchmark: $(addprefix benchmark/,$(obj))
$(addprefix benchmark/,$(obj)):
	go test -timeout 3600s -bench=./$(shell echo $(subst benchmark/,,$@) | cut -d / -f1)/... ./$(shell echo $(subst benchmark/,,$@) | cut -d / -f1)/...

# Clean
clean:
	rm -rf out

# Dependencies
depend: $(addprefix depend/,$(obj))
$(addprefix depend/,$(obj)):
	cd ./$(shell echo $(subst depend/,,$@) | cut -d / -f1) && go generate ./...
