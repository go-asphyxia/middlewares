package middlewares

import (
	"bytes"
	"strings"

	"github.com/valyala/fasthttp"

	aconversion "github.com/go-asphyxia/conversion"
	aheaders "github.com/go-asphyxia/http/headers"
)

type (
	CORS struct {
		Origins, Methods, Headers []string
	}
)

func NewCORS(origins []string, methods, headers []string) (cors *CORS) {
	cors = &CORS{
		Origins: origins,
		Methods: methods,
		Headers: headers,
	}

	return
}

func (cors *CORS) Middleware(source fasthttp.RequestHandler) (target fasthttp.RequestHandler) {
	m := strings.Join(cors.Methods, ",")
	h := strings.Join(cors.Headers, ",")

	prefix := []byte("https://")

	target = func(ctx *fasthttp.RequestCtx) {
		origin := aconversion.BytesToStringNoCopy(bytes.TrimPrefix(ctx.Request.Header.Peek(aheaders.Origin), prefix))

		for i := range cors.Origins {
			if cors.Origins[i] == origin {
				headers := &ctx.Response.Header

				headers.Set(aheaders.AccessControlAllowOrigin, origin)
				headers.Set(aheaders.AccessControlAllowMethods, m)
				headers.Set(aheaders.AccessControlAllowHeaders, h)

				source(ctx)
			}
		}
	}

	return
}

func (cors *CORS) Handler() (handler fasthttp.RequestHandler) {
	m := strings.Join(cors.Methods, ",")
	h := strings.Join(cors.Headers, ",")

	prefix := []byte("https://")

	handler = func(ctx *fasthttp.RequestCtx) {
		origin := aconversion.BytesToStringNoCopy(bytes.TrimPrefix(ctx.Request.Header.Peek(aheaders.Origin), prefix))

		for i := range cors.Origins {
			if cors.Origins[i] == origin {
				headers := &ctx.Response.Header

				headers.Set(aheaders.AccessControlAllowOrigin, origin)
				headers.Set(aheaders.AccessControlAllowMethods, m)
				headers.Set(aheaders.AccessControlAllowHeaders, h)
			}
		}
	}

	return
}
