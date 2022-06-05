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
	m := strings.Join(cors.Methods, ",")
	h := strings.Join(cors.Headers, ",")

	if len(cors.Origins) > 0 {
		o := make([]string, (len(cors.Origins) * 2))

		for i := range cors.Origins {
			o = append(o, ("http://" + cors.Origins[i]), ("https://" + cors.Origins[i]))
		}

		target = func(ctx *fasthttp.RequestCtx) {
			headers := &ctx.Response.Header

			headers.Set(aheaders.AccessControlAllowMethods, m)
			headers.Set(aheaders.AccessControlAllowHeaders, h)

			origin := aconversion.BytesToStringNoCopy(ctx.Request.Header.Peek(aheaders.Origin))

			for i := range o {
				if o[i] == origin {
					headers.Set(aheaders.AccessControlAllowOrigin, origin)

					source(ctx)
					return
				}
			}
		}

		return
	}

	target = func(ctx *fasthttp.RequestCtx) {
		headers := &ctx.Response.Header

		headers.Set(aheaders.AccessControlAllowMethods, m)
		headers.Set(aheaders.AccessControlAllowHeaders, h)
		headers.Set(aheaders.AccessControlAllowOrigin, "*")

		source(ctx)
	}

	return
}

func (cors *CORS) Handler() (handler fasthttp.RequestHandler) {
	m := strings.Join(cors.Methods, ",")
	h := strings.Join(cors.Headers, ",")

	if len(cors.Origins) > 0 {
		o := make([]string, (len(cors.Origins) * 2))

		for i := range cors.Origins {
			o = append(o, ("http://" + cors.Origins[i]), ("https://" + cors.Origins[i]))
		}

		handler = func(ctx *fasthttp.RequestCtx) {
			headers := &ctx.Response.Header

			headers.Set(aheaders.AccessControlAllowMethods, m)
			headers.Set(aheaders.AccessControlAllowHeaders, h)

			origin := aconversion.BytesToStringNoCopy(ctx.Request.Header.Peek(aheaders.Origin))

			for i := range o {
				if o[i] == origin {
					headers.Set(aheaders.AccessControlAllowOrigin, origin)
					return
				}
			}
		}

		return
	}

	handler = func(ctx *fasthttp.RequestCtx) {
		headers := &ctx.Response.Header

		headers.Set(aheaders.AccessControlAllowMethods, m)
		headers.Set(aheaders.AccessControlAllowHeaders, h)
		headers.Set(aheaders.AccessControlAllowOrigin, "*")
	}

	return
}
