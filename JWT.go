package middlewares

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/valyala/fasthttp"

	aconversion "github.com/go-asphyxia/conversion"
	aheaders "github.com/go-asphyxia/http/headers"
	arandom "github.com/go-asphyxia/random"
)

type (
	TokenGeneratorConfiguration struct {
		Key string

		SignatureLength int
		RefreshLength   int

		TokenExpires   int
		RefreshExpires int
	}

	TokenGeneratorClaims[T any] struct {
		jwt.StandardClaims

		Claims T
	}

	TokenGenerator[T any] struct {
		Key string

		TokenSignature   []byte
		RefreshSignature []byte

		TokenExpires   time.Duration
		RefreshExpires time.Duration
	}
)

func NewTokenGenerator[T any](generator *arandom.Generator, c *TokenGeneratorConfiguration) (tg *TokenGenerator[T]) {
	tg = &TokenGenerator[T]{
		TokenSignature:   generator.Bytes(c.SignatureLength, arandom.DefaultCharset),
		RefreshSignature: generator.Bytes(c.RefreshLength, arandom.DefaultCharset),

		TokenExpires:   time.Minute * time.Duration(c.TokenExpires),
		RefreshExpires: time.Minute * time.Duration(c.RefreshExpires),
	}

	return
}

func (tg *TokenGenerator[T]) JWT(source fasthttp.RequestHandler) (target fasthttp.RequestHandler) {
	parser := jwt.NewParser()

	tgce := &TokenGeneratorClaims[T]{}

	target = func(ctx *fasthttp.RequestCtx) {
		authorization := aconversion.BytesToStringNoCopy((&ctx.Request.Header).Peek(aheaders.Authorization))

		t, err := parser.ParseWithClaims(authorization, tgce, tg.TokenFunction)

		if err != nil {
			ctx.Error(err.Error(), fasthttp.StatusUnauthorized)
			return
		}

		claims, ok := t.Claims.(*TokenGeneratorClaims[T])

		if !ok {
			ctx.Error(jwt.ErrTokenInvalidClaims.Error(), fasthttp.StatusUnauthorized)
			return
		}

		ctx.SetUserValue(tg.Key, claims.Claims)

		source(ctx)
	}

	return
}

func (tg *TokenGenerator[T]) TokenFunction(t *jwt.Token) (itf interface{}, err error) {
	itf = tg.TokenSignature
	return
}

func (tg *TokenGenerator[T]) RefreshFunction(t *jwt.Token) (itf interface{}, err error) {
	itf = tg.RefreshSignature
	return
}

func (tg *TokenGenerator[T]) Generate(claims T) (token, refresh string, err error) {
	tgc := &TokenGeneratorClaims[T]{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tg.TokenExpires).Unix(),
		},
		Claims: claims,
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, tgc).SignedString(tg.TokenSignature)
	if err != nil {
		return
	}

	tgr := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(tg.RefreshExpires).Unix(),
	}

	refresh, err = jwt.NewWithClaims(jwt.SigningMethodHS256, tgr).SignedString(tg.RefreshSignature)
	return
}
