go ?= GO111MODULE=on go
tag_prefix ?=

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

go_files = go.mod go.sum $(shell $(go) list -f '{{$$Dir := .Dir}}{{range .GoFiles}}{{$$Dir}}/{{.}} {{end}}' ./...)
go_files_last_commit = $(shell git log -1 --format=%H -- $(go_files))

## Show module version as expected by the Go toolchain
go-version: $(go_files)
	@{ git describe --tags --match '$(tag_prefix)v*.*.*' --exact-match 2>/dev/null $(shell git log -1 --format=%H -- $^ ) || TZ=UTC git log -1 '--date=format-local:%Y%m%d%H%M%S' --abbrev=12 '--pretty=tformat:%(describe:tags,match=$(tag_prefix)v*,abbrev=0)-%cd-%h' $^ | perl -pE 's/(\d+)(?=-)/$$1+1/e' ; } | sed -e 's!.*/!!'

## Show "go get" command to upgrade the module in a downstream project
go-get:
	@echo $(go) get -d $(shell $(go) list .)@$(shell $(MAKE) "go=$(go)" go-version)

.PHONY: next.minor next.patch tag.minor tag.patch

## Show next minor tag to create: prefix/vX.Y.Z -> prefix/vX.(Y+1).0
next.minor:
	@git tag -l --sort=-v:refname $(tag_prefix)'v*' | perl -E '$$_=<>; s/\.([0-9]+)\..*$$/".".($$1+1).".0"/e; print'

## Show next patch tag to create: prefix/vX.Y.Z -> prefix/vX.Y.(Z+1)
next.patch:
	@git tag -l --sort=-v:refname $(tag_prefix)'v*' | perl -E '$$_=<>; s/\.([0-9]+)$$/".".($$1+1)/e; print'

## Tag a new release, increasing the minor version: prefix/vX.Y.Z -> prefix/vX.(Y+1).0
tag.minor:
	git tag -a $$(git tag -l --sort=-v:refname $(tag_prefix)'v*' | perl -E '$$_=<>; s/\.([0-9]+)\..*$$/".".($$1+1).".0"/e; print') $(go_files_last_commit)

## Tag a new release, increasing the patch version: prefix/vX.Y.Z -> prefix/vX.Y.(Z+1)
tag.patch:
	git tag -a $$(git tag -l --sort=-v:refname $(tag_prefix)'v*' | perl -E '$$_=<>; s/\.([0-9]+)$$/".".($$1+1)/e; print') $(go_files_last_commit)

.PHONY:	upgrade-blobqueue

## After a release of the main blobqueue module, upgrade blobqueue in sub modules and tag a minor release
upgrade-blobqueue:
	git push --tags
	V=$$($(MAKE) go-get); for mod in typedqueue queueredis queuemsgpack; do ( cd $$mod; go get github.com/blueboardio/go-blobqueue@$$V && go mod tidy ); done
	git commit -a -m "typedqueue queueredis queuemsgpack: upgrade blobqueue"
	git tag -a $$(git tag -l --sort=-v:refname $(tag_prefix)'v*' | perl -E '$$_=<>; s/\.([0-9]+)$$/".".($$1+1)/e; print')

.PHONY: bump-tag edit-tag changelog

## Bump last non-pushed tag to HEAD
bump-tag:
	# Check if tag has already been pushed...
	t=$$(git tag -l --sort=-v:refname $(tag_prefix)'v*' | head -n1); ! git ls-remote --exit-code --tags origin $$t
	t=$$(git tag -l --sort=-v:refname $(tag_prefix)'v*' | head -n1); git tag -f -a -m "$$(git tag -l '--format=%(contents)' $$t)" $$t

## Edit the message attached to the last tag
edit-tag:
	# Check if tag has already been pushed...
	t=$$(git tag -l --sort=-v:refname $(tag_prefix)'v*' | head -n1); ! git ls-remote --exit-code --tags origin $$t
	t=$$(git tag -l --sort=-v:refname $(tag_prefix)'v*' | head -n1); git tag -f -a $$t $$t^{}


## Dump changelog from Git tags
changelog:
	@git tag -l --sort=-v:refname "--format=[%(refname:short)] %(contents)*****************************" $(tag_prefix)'v*'
