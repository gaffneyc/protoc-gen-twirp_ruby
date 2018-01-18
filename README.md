# Twirp RPC Ruby Client Generator

## Install

```bash
go get -u github.com/gaffneyc/protoc-gen-twirp_ruby
```

## Usage

Using the [Haberdasher example](https://github.com/twitchtv/twirp/wiki/Usage-Example:-Haberdasher):

```bash
protoc --proto_path=$GOPATH/src:. --twirp_ruby_out=. --ruby_out=. path/to/service.proto
```

```ruby
require "service_pb"
require "service_twirp"

client = Haberdasher::HaberdasherClient.new("http://localhost:8081")
puts client.make_hat(Haberdasher::Size.new(inches: 5)).inspect
```
