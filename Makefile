#====================
AUTHOR         ?= The sacloud/sacloud-usage-lib Authors
COPYRIGHT_YEAR ?= 2023-2025

GO_FILES       ?= $(shell find . -name '*.go')

include includes/go/common.mk
#====================

default: $(DEFAULT_GOALS)
tools: dev-tools
