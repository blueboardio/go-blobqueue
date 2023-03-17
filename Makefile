

go ?= GO111MODULE=on go

go-version: go.mod $(shell $(go) list -f '{{$$Dir := .Dir}}{{range .GoFiles}}{{$$Dir}}/{{.}} {{end}}' ./...)
	@git describe --tags --match 'v*.*.*' --exact-match 2>/dev/null $(shell git log -1 --format=%H -- $^ ) || TZ=UTC git log -1 '--date=format-local:%Y%m%d%H%M%S' --abbrev=12 '--pretty=tformat:%(describe:tags,match=v*,abbrev=0)-%cd-%h' $^ | perl -pE 's/(\d+)(?=-)/$$1+1/e'

go-get:
	@echo $(go) get -d $(shell $(go) list .)@$(shell $(MAKE) "go=$(go)" go-version)

