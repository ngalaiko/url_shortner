package params

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/ngalayko/url_shortner/server/logger"
)

// Parse args parses http arguments into dest map by rules
func ParseArgs(logger *logger.Logger, args *fasthttp.Args, parseRules map[string]parseRule, dest map[string]interface{}) error {

	args.VisitAll(func(key, value []byte) {

		parseFunc, ok := parseRules[string(key)]
		if !ok {
			logger.Error("no rule to parse key",
				zap.ByteString("key", key),
				zap.ByteString("value", value),
			)
		}

		dest[string(key)] = parseFunc(value)
	})

	return nil
}
