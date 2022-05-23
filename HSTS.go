package middlewares

import (
	"strconv"

	"github.com/valyala/fasthttp"

	aheaders "github.com/go-asphyxia/http/headers"
)

type (
	HSTS struct {
		MaxAge int
	}
)

func NewHSTS(maxAge int) (hsts *HSTS) {
	hsts = &HSTS{
		MaxAge: maxAge,
	}

	return
}

func (hsts *HSTS) Middleware(source fasthttp.RequestHandler) (target fasthttp.RequestHandler) {
	strictTransportSecurity := "max-age=" + strconv.Itoa(hsts.MaxAge) + "; includeSubDomains; preload"

	target = func(ctx *fasthttp.RequestCtx) {
		headers := &ctx.Response.Header

		if ctx.IsTLS() {
			headers.Set(aheaders.StrictTransportSecurity, strictTransportSecurity)

			source(ctx)
			return
		}

		uri := ctx.Request.URI()
		uri.SetScheme("https")

		ctx.Response.SetStatusCode(fasthttp.StatusMovedPermanently)
		headers.SetBytesV(aheaders.Location, uri.FullURI())
	}

	return
}
