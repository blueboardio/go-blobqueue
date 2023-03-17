go ?= GO111MODULE=on go

.PHONY: help

## Display this help screen
help: $(MAKEFILE_LIST)
	@printf "\e[36m%-35s %s\e[0m\n" Target Description
	@sed -n -e '/^## /{'\
		-e 's/## //g;'\
		-e 'h;'\
		-e 'n;'\
		-e 's/:.*//g;'\
		-e 'G;'\
		-e 's/\n/ /g;'\
		-e 'p;}' $^ | sort | awk '{printf "\033[33m%-35s\033[0m%s\n", $$1, substr($$0,length($$1)+1)}'


.PHONY: go-version go-get

## Show module version as expected by the Go toolchain
go-version: go.mod go.sum $(shell $(go) list -f '{{$$Dir := .Dir}}{{range .GoFiles}}{{$$Dir}}/{{.}} {{end}}' ./...)
	@git describe --tags --match 'v*.*.*' --exact-match 2>/dev/null $(shell git log -1 --format=%H -- $^ ) || TZ=UTC git log -1 '--date=format-local:%Y%m%d%H%M%S' --abbrev=12 '--pretty=tformat:%(describe:tags,match=v*,abbrev=0)-%cd-%h' $^ | perl -pE 's/(\d+)(?=-)/$$1+1/e'

## Show "go get" command to upgrade the module in a downstream project
go-get:
	@echo $(go) get -d $(shell $(go) list .)@$(shell $(MAKE) "go=$(go)" go-version)

.PHONY: tag.minor tag.patch upgrade-blobqueue

## Tag a new release, increasing the minor version: x.y.z -> x.(y+1).0
tag.minor:
	git tag -a $$(git tag -l --sort=-v:refname 'v*' | perl -E '$$_=<>; s/\.([0-9]+)\..*$$/".".($$1+1).".0"/e; print')

## Tag a new release, increasing the patch version: x.y.z -> x.y.(z+1)
tag.patch:
	git tag -a $$(git tag -l --sort=-v:refname 'v*' | perl -E '$$_=<>; s/\.([0-9]+)$$/".".($$1+1)/e; print')

## After a release of the main blobqueue module, upgrade blobqueue in sub modules and tag a minor release
upgrade-blobqueue:
	git push --tags
	V=$$($(MAKE) go-get); for mod in typedqueue queueredis queuemsgpack; do ( cd $$mod; go get github.com/blueboardio/go-blobqueue@$$V && go mod tidy ); done
	git commit -a -m "typedqueue queueredis queuemsgpack: upgrade blobqueue"
	git tag -a $$(git tag -l --sort=-v:refname 'v*' | perl -E '$$_=<>; s/\.([0-9]+)$$/".".($$1+1)/e; print')

.PHONY: bump-tag edit-tag changelog

## Bump last non-pushed tag to HEAD
bump-tag:
	# Check if tag has already been pushed...
	t=$$(git tag -l --sort=-v:refname 'v*' | head -n1); ! git ls-remote --exit-code --tags origin $$t
	t=$$(git tag -l --sort=-v:refname 'v*' | head -n1); git tag -f -a -m "$$(git tag -l '--format=%(contents)' $$t)" $$t

## Edit the message attached to the last tag
edit-tag:
	# Check if tag has already been pushed...
	t=$$(git tag -l --sort=-v:refname 'v*' | head -n1); ! git ls-remote --exit-code --tags origin $$t
	t=$$(git tag -l --sort=-v:refname 'v*' | head -n1); git tag -f -a $$t $$t^{}


## Dump changelog from Git tags
changelog:
	@git tag -l --sort=-v:refname "--format=[%(refname:short)] %(contents)*****************************" 'v*'
