package middlewares

import (
	"strings"

	"github.com/valyala/fasthttp"

	aconversion "github.com/go-asphyxia/conversion"
	aheaders "github.com/go-asphyxia/http/headers"
)

type (
	CORS struct {
		Hosts, Methods, Headers []string
	}
)

func NewCORS(hosts []string, methods, headers []string) (cors *CORS) {
	cors = &CORS{
		Hosts:   hosts,
		Methods: methods,
		Headers: headers,
	}

	return
}

func (cors *CORS) Middleware(source fasthttp.RequestHandler) (target fasthttp.RequestHandler) {
	m := strings.Join(cors.Methods, ",")
	h := strings.Join(cors.Headers, ",")

	target = func(ctx *fasthttp.RequestCtx) {
		host := aconversion.BytesToStringNoCopy(ctx.Request.Host())

		for i := range cors.Hosts {
			if cors.Hosts[i] == host {
				headers := &ctx.Response.Header

				headers.Set(aheaders.AccessControlAllowOrigin, host)
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

	handler = func(ctx *fasthttp.RequestCtx) {
		host := aconversion.BytesToStringNoCopy(ctx.Request.Host())

		for i := range cors.Hosts {
			if cors.Hosts[i] == host {
				headers := &ctx.Response.Header

				headers.Set(aheaders.AccessControlAllowOrigin, host)
				headers.Set(aheaders.AccessControlAllowMethods, m)
				headers.Set(aheaders.AccessControlAllowHeaders, h)
			}
		}
	}

	return
}
