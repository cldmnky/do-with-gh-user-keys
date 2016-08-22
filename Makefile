package = github.com/cldmnky/do-with-gh-user-keys

.PHONY: release

release:
	mkdir -p release
	GOOS=linux GOARCH=amd64 go build -o release/do-with-gh-user-keys-linux-amd64 $(package)
