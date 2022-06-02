package middlewares

import (
	"strings"

	"github.com/valyala/fasthttp"

	aconversion "github.com/go-asphyxia/conversion"
	aheaders "github.com/go-asphyxia/http/headers"
)

type (
	CORS struct {
		Scheme                    string
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
	o := make([]string, len(cors.Origins))

	for i := range cors.Origins {
		o = append(o, ("http://" + cors.Origins[i]), ("https://" + cors.Origins[i]))
	}

	m := strings.Join(cors.Methods, ",")
	h := strings.Join(cors.Headers, ",")

	target = func(ctx *fasthttp.RequestCtx) {
		origin := aconversion.BytesToStringNoCopy(ctx.Request.Header.Peek(aheaders.Origin))

		for i := range o {
			if o[i] == origin {
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
	o := make([]string, len(cors.Origins))

	for i := range cors.Origins {
		o = append(o, ("http://" + cors.Origins[i]), ("https://" + cors.Origins[i]))
	}

	m := strings.Join(cors.Methods, ",")
	h := strings.Join(cors.Headers, ",")

	handler = func(ctx *fasthttp.RequestCtx) {
		origin := aconversion.BytesToStringNoCopy(ctx.Request.Header.Peek(aheaders.Origin))

		for i := range o {
			if o[i] == origin {
				headers := &ctx.Response.Header

				headers.Set(aheaders.AccessControlAllowOrigin, origin)
				headers.Set(aheaders.AccessControlAllowMethods, m)
				headers.Set(aheaders.AccessControlAllowHeaders, h)
			}
		}
	}

	return
}
