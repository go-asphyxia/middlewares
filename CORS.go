package middlewares

import (
	"bytes"
	"strings"

	"github.com/valyala/fasthttp"

	aheaders "github.com/go-asphyxia/http/headers"
)

type (
	CORS struct {
		Host string

		Methods, Headers []string
	}
)

func NewCORS(host string, methods, headers []string) (cors *CORS) {
	cors = &CORS{
		Host:    host,
		Methods: methods,
		Headers: headers,
	}
	
	return
}

func (cors *CORS) Middleware(source fasthttp.RequestHandler) (target fasthttp.RequestHandler) {
	host := []byte(cors.Host)

	m := strings.Join(cors.Methods, ",")
	h := strings.Join(cors.Headers, ",")

	target = func(ctx *fasthttp.RequestCtx) {
		if bytes.HasSuffix(ctx.Request.Host(), host) {
			headers := &ctx.Response.Header

			headers.SetBytesV(aheaders.AccessControlAllowOrigin, host)
			headers.Set(aheaders.AccessControlAllowMethods, m)
			headers.Set(aheaders.AccessControlAllowHeaders, h)

			source(ctx)
		}
	}

	return
}

func (cors *CORS) Handler() (handler fasthttp.RequestHandler) {
	host := []byte(cors.Host)

	m := strings.Join(cors.Methods, ",")
	h := strings.Join(cors.Headers, ",")

	handler = func(ctx *fasthttp.RequestCtx) {
		if bytes.HasSuffix(ctx.Request.Host(), host) {
			headers := &ctx.Response.Header

			headers.SetBytesV(aheaders.AccessControlAllowOrigin, host)
			headers.Set(aheaders.AccessControlAllowMethods, m)
			headers.Set(aheaders.AccessControlAllowHeaders, h)
		}
	}

	return
}
