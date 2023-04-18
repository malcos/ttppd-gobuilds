# use always go mod init before working on any go projects, go modules are not perfect but better than the old strategy
export GO111MODULE=on

# folder name for output files
TARGET = target

# extension for the symbolic links created in the root folder, use something unique or "clean" may produce collateral damage
LINKEXT = .link

# the home or parent folder for the command(s). It must start with ./ otherwise the go compiler will complain
CMDHOME = ./cmd-single

# the name of command to compile (only used to override default name when in single command mode)
# CMDNAME =

# list of target platforms in the format "os_arch". Host OS and ARCH will be appended automatically, duplications will be removed
# PLATFORMS =
# PLATFORMS += linux_386
# PLATFORMS += linux_amd64
# PLATFORMS += darwin_amd64
# PLATFORMS += windows_386
# PLATFORMS += windows_amd64

# additional build arguments
# BUILDARGS = -mod=readonly -ldflags="-X 'main.Version=$(VERSION)' -X 'main.GoVersion=$(GOVERSION)'"
# BUILDARGS = -mod=readonly -ldflags="-X 'main.Version=$(VERSION)' -X 'main.GoVersion=$(GOVERSION)' -X 'main.GitCommit=$(shell git rev-parse HEAD)'"
# BUILDARGS =

# Other variables accessible from the included file
# GOVERSION       // info section: the version of the go compiler
# GOMODULE_PATH   // build section: the full name of the module as stated in the go.mod file
# GOMODULE_NAME   // build section: the project name extracted from GOMODULE_PATH
# HOST_OS         // build section: the OS of the machine where the builds are running
# HOST_ARCH       // build section: the architecture of the machine where the builds are running

# options are "zip", "tgz" and "plain" (no compression), it will be used as default target in the "package" rule
PACKAGE_COMPRESSION = zip

# name of the checksum application to use, it must accept "*.ext" as single parameter to work: like sha256sum, sha1sum or md5sum
PACKAGE_CHECKSUM = md5sum

# Assign this to version your artifacts (there is a nice additional module I wrote for this as well)
# VERSION = 1.7.42

# you can define your default goal here
.DEFAULT_GOAL := local

#
# :: Building artifacts
#

# get the general go module full nname and project name
GOMODULE_PATH := $(word 2, $(file < go.mod))
GOMODULE_NAME := $(lastword $(subst /, ,$(GOMODULE_PATH)))

# get the host OS and ARCH, append it to the existing platforms and remove duplicates
HOST_OS := $(shell go env GOHOSTOS)
HOST_ARCH := $(shell go env GOHOSTARCH)
PLATFORMS += $(HOST_OS)_$(HOST_ARCH)
PLATFORMS := $(sort $(PLATFORMS))

# determine the list of commands based on the contents of the CMDHOME folder
ifeq ($(wildcard $(CMDHOME)/*.go),)
# no go files in folder, assuming subfolders as commands
COMMANDS = $(notdir $(wildcard $(CMDHOME)/*))
COMMAND_SOURCE = $(CMDHOME)/$(CMDNAME)
else ifeq ($(strip $(CMDNAME)),)
# go files in folder and CMDNAME is empty, single command named after go module
COMMANDS = $(GOMODULE_NAME)
COMMAND_SOURCE = $(CMDHOME)
else
# go files in folder and CMDNAME is set, single command named after CMDNAME
COMMANDS = $(CMDNAME)
COMMAND_SOURCE = $(CMDHOME)
endif

# all possible combinations between COMMANDS and PLATFORMS in the format "command_os_arch"
MATRIX = $(foreach command,$(COMMANDS),$(foreach platform,$(PLATFORMS),$(command)_$(platform)))

# list of dynamically generated goals
GOALS_BUILD = $(addprefix build.,$(MATRIX))
GOALS_COMMANDS = $(addprefix command.,$(COMMANDS))
GOALS_PLATFORMS = $(addprefix platform.,$(PLATFORMS))
GOALS_LOCAL = $(addprefix local.,$(COMMANDS))

all build: $(GOALS_BUILD)

# build.CMDNAME_OS_ARCH
$(GOALS_BUILD):
	$(eval PARTS=$(subst _, , $(lastword $(subst ., ,$@))))
	$(eval CMDNAME=$(word 1, $(PARTS)))
	$(eval GOOS=$(word 2, $(PARTS)))
	$(eval GOARCH=$(word 3, $(PARTS)))
	$(eval EXT=$(if $(filter windows,$(GOOS)),.exe,))
	$(eval TARGET_BINARY=$(TARGET)/bin/$(CMDNAME)$(addprefix _,$(VERSION))_$(GOOS)_$(GOARCH)$(EXT))
	env GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(BUILDARGS) -o $(TARGET_BINARY) $(COMMAND_SOURCE)
	$(if $(SYMBOLIC), ln -sf $(TARGET_BINARY) $(CMDNAME)$(LINKEXT))

# command.CMDNAME
$(GOALS_COMMANDS):
	$(eval CMDNAME=$(lastword $(subst ., ,$@)))
	$(MAKE) $(foreach platform, $(PLATFORMS), build.$(CMDNAME)_$(platform))

# platform.OS_ARCH
$(GOALS_PLATFORMS):
	$(eval OS_ARCH=$(lastword $(subst ., ,$@)))
	$(MAKE) $(foreach command, $(COMMANDS), build.$(command)_$(OS_ARCH))

local: .symbolic $(foreach command, $(COMMANDS), build.$(command)_$(HOST_OS)_$(HOST_ARCH))

.symbolic:
	$(eval SYMBOLIC=1)

# latest.CMDNAME
$(GOALS_LOCAL):
	$(eval CMDNAME=$(lastword $(subst ., ,$@)))
	$(MAKE) SYMBOLIC=1 build.$(CMDNAME)_$(HOST_OS)_$(HOST_ARCH)

release: clean test.coverage build package

clean:
	@rm -f *$(LINKEXT)
	@rm -fr $(TARGET)

#
# :: Testing
#

COVERPROFILE = $(TARGET)/coverage.out

test:
	go test -count=1 ./...

test.coverage:
	@mkdir -p $(TARGET)
	go test -cover -coverprofile $(COVERPROFILE) -count=1 ./...

test.coverage.view: test.coverage
	go tool cover -html=$(COVERPROFILE)

#
# :: Packaging
#

# uppercase name of the software used for checksums
PACKAGE_CHECKSUM_UC = $(call __uppercase,$(PACKAGE_CHECKSUM))
CHECKSUM_FILE = $(GOMODULE_NAME)$(addprefix _,$(VERSION))_$(PACKAGE_CHECKSUM_UC)
CHECKSUM_CREATE = cd $(TARGET)/release; $(PACKAGE_CHECKSUM) $1 > $(CHECKSUM_FILE); $(PACKAGE_CHECKSUM) --check $(CHECKSUM_FILE)

package: package.$(PACKAGE_COMPRESSION)

.package.init:
	@rm -rf $(TARGET)/release
	@mkdir -p $(TARGET)/release
	$(eval ARTIFACT_NAMES:=$(notdir $(wildcard $(TARGET)/bin/*)))

package.zip: .package.init
	@$(foreach artifact,$(ARTIFACT_NAMES),zip -j $(TARGET)/release/$(artifact).zip $(TARGET)/bin/$(artifact);)
	@$(call CHECKSUM_CREATE,*.zip)

package.tgz: .package.init
	@$(foreach artifact,$(ARTIFACT_NAMES),tar -cvzf $(TARGET)/release/$(artifact).tgz -C $(TARGET)/bin $(artifact);)
	@$(call CHECKSUM_CREATE,*.tgz)

package.plain: .package.init
	@$(foreach artifact,$(ARTIFACT_NAMES),cp -av $(TARGET)/bin/$(artifact) $(TARGET)/release;)
	@$(call CHECKSUM_CREATE,*)

__lowercase = $(subst A,a,$(subst B,b,$(subst C,c,$(subst D,d,$(subst E,e,$(subst F,f,$(subst G,g,$(subst H,h,$(subst I,i,$(subst J,j,$(subst K,k,$(subst L,l,$(subst M,m,$(subst N,n,$(subst O,o,$(subst P,p,$(subst Q,q,$(subst R,r,$(subst S,s,$(subst T,t,$(subst U,u,$(subst V,v,$(subst W,w,$(subst X,x,$(subst Y,y,$(subst Z,z,$1))))))))))))))))))))))))))
__uppercase = $(subst a,A,$(subst b,B,$(subst c,C,$(subst d,D,$(subst e,E,$(subst f,F,$(subst g,G,$(subst h,H,$(subst i,I,$(subst j,J,$(subst k,K,$(subst l,L,$(subst m,M,$(subst n,N,$(subst o,O,$(subst p,P,$(subst q,Q,$(subst r,R,$(subst s,S,$(subst t,T,$(subst u,U,$(subst v,V,$(subst w,W,$(subst x,X,$(subst y,Y,$(subst z,Z,$1))))))))))))))))))))))))))

#
# :: Misc
#

GOVERSION := $(shell go env GOVERSION)

print.gitignore:
	@echo "$(TARGET)/"
	@echo "*$(LINKEXT)"

print.info:
	@echo "MODULE    : $(GOMODULE_PATH) $(if $(VERSION),$(VERSION),(no version))"
	@echo "COMMANDS  : $(COMMANDS)"
	@echo "PLATFORMS : $(PLATFORMS)"
	@echo "HOST      : $(HOST_OS)_$(HOST_ARCH)"
	@echo "COMPILER  : $(GOVERSION)"
	@echo "RULES     : $(words $(GOALS_COMMANDS)) commands, $(words $(GOALS_PLATFORMS)) platforms, $(words $(GOALS_BUILD)) build artifacts"
	@echo "PACKAGING : $(PACKAGE_COMPRESSION) files with $(PACKAGE_CHECKSUM_UC) checks"
