# LeanCloud Go SDK
Golang SDK for LeanCloud Storage and LeanEngine.

```go
import "github.com/leancloud/go-sdk/leancloud"
```

## Examples

- [LeanEngine Getting Started](https://github.com/leancloud/golang-getting-started)

## Documentation

- [Go SDK Setup](https://leancloud.cn/docs/sdk_setup-go.html)
- [API Reference](https://pkg.go.dev/github.com/leancloud/go-sdk/leancloud)

## Development

Release:

- Update `Version` in `leancloud/client.go`
- `git tag v<major>.<minor>.<patch>`
- Update pkg.go.dev via `GOPROXY=https://proxy.golang.org GO111MODULE=on go get github.com/leancloud/go-sdk@v<major>.<minor>.<patch>`
