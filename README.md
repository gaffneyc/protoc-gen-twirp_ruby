# Twirp RPC Ruby Client Generator

## Install

```bash
go get -u github.com/gaffneyc/protoc-gen-twirp_ruby
```

## Usage

```bash
protoc --proto_path=$GOPATH/src:. --twirp_ruby_out=. --ruby_out=. service.proto
```
