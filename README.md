# Puma Proxy

A really simple proxy for Ruby's [Puma web server](https://puma.io/) written in golang. It accepts HTTP requests on a given port, and proxies them to a Puma unix socket. The main benefit is that it also serves static assets, so Ruby does not have to handle that.

_Not meant for production_. Use in development or testing environments only.

# Install

Option 1: Install with homebrew -> `brew install ianks/tap/puma-proxy`
Option 2: Install with go -> `go get github.com/ianks/puma-proxy`
Option 3: Download the tarball from the releases page 
# Usage

```sh
$ puma-proxy -- bundle exec puma -b unix:///tmp/puma.sock

# Same as 👆, just writing the arguments out to show what's configurable
$ puma-proxy -listen=localhost:3000 -sock=/tmp/puma.sock -- bundle exec puma -b unix:///tmp/puma.sock

# Running with no command on the end, puma wont be launch but the proxy will still run
$ puma-proxy
```
