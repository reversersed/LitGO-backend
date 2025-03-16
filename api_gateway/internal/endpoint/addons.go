package endpoint

import (
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func cookieHeaderMatcher(key string) (string, bool) {
	if strings.Contains(key, "set-cookie") {
		return "Set-Cookie", true
	}
	return runtime.DefaultHeaderMatcher(key)
}
