package token

import (
	"github.com/nbkit/mdf/framework/reg"
	"github.com/nbkit/mdf/gin"
	"github.com/nbkit/mdf/utils"
	"net/http"
	"regexp"
	"strings"
)

type handler struct {
}

func newHandler(config Config) *handler {
	h := &handler{}
	return h
}
func (g *handler) isEmptyContext(context *utils.TokenContext) bool {
	if context == nil || context.ID() == "" {
		return true
	}
	return false

}
func (g *handler) Handle(c *gin.Context) {
	var (
		context *utils.TokenContext
		isHint  bool
		ignore  bool
	)
	ignorePattern := strings.Split(utils.Config.GetValue("oauth.ignore"), ",")
	for _, item := range ignorePattern {
		match := regexp.MustCompile(item)
		if ok := match.MatchString(c.Request.RequestURI); ok {
			ignore = true
			break
		}
	}
	context = Get(c)
	//authorization
	if g.isEmptyContext(context) {
		context, isHint = g.tryFromHeaderAuthorization(c)
		if isHint && !ignore && g.isEmptyContext(context) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
	if g.isEmptyContext(context) {
		context, isHint = g.tryFromCookieAuthorization(c)
		if isHint && !ignore && g.isEmptyContext(context) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
	//token
	if !ignore && g.isEmptyContext(context) {
		context, isHint = g.tryFromHeaderToken(c)
		if isHint && g.isEmptyContext(context) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
	if !ignore && g.isEmptyContext(context) {
		context, isHint = g.tryFromURLParamToken(c)
		if isHint && g.isEmptyContext(context) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
	if !ignore && g.isEmptyContext(context) {
		context, isHint = g.tryFromCookieToken(c)
		if isHint && g.isEmptyContext(context) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
	if g.isEmptyContext(context) && !ignore {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if g.isEmptyContext(context) {
		context = utils.NewTokenContext()
	}
	c.Set(DefaultContextKey, context)

	c.Next()
}
func (g *handler) tryFromURLParamToken(c *gin.Context) (*utils.TokenContext, bool) {
	if token := c.Param(AuthTokenKey); token != "" {
		if ct, err := reg.GetTokenContext(token); ct != nil && err == nil {
			return ct, true
		}
		return nil, true
	}
	return nil, false
}

func (g *handler) tryFromHeaderToken(c *gin.Context) (*utils.TokenContext, bool) {
	if token := c.GetHeader(AuthTokenKey); token != "" {
		if ct, err := reg.GetTokenContext(token); ct != nil && err == nil {
			return ct, true
		}
		return nil, true
	}
	return nil, false
}

func (g *handler) tryFromCookieToken(c *gin.Context) (*utils.TokenContext, bool) {
	if token, _ := c.Cookie(AuthTokenKey); token != "" {
		if ct, err := reg.GetTokenContext(token); ct != nil && err == nil {
			return ct, true
		}
		return nil, true
	}
	return nil, false
}
func (g *handler) tryFromCookieAuthorization(c *gin.Context) (*utils.TokenContext, bool) {
	if ck, _ := c.Cookie(AuthSessionKey); ck != "" {
		context := g.parseJWTToken(ck)
		if context != nil && context.IsValid() {
			return context, true
		}
		return nil, true
	}
	return nil, false
}
func (g *handler) tryFromHeaderAuthorization(c *gin.Context) (*utils.TokenContext, bool) {
	token := c.GetHeader("Authorization")
	if token == "" {
		return nil, false
	}
	authHeaderParts := strings.Split(token, " ")
	if len(authHeaderParts) == 2 && strings.ToLower(authHeaderParts[0]) == "bearer" && authHeaderParts[1] != "" {
		return g.parseJWTToken(authHeaderParts[1]), true
	}
	return nil, false
}
func (g *handler) parseJWTToken(token string) *utils.TokenContext {
	if context, err := utils.NewTokenContext().FromTokenString(token); err != nil || context == nil || !context.IsValid() {
		return nil
	} else {
		return context
	}
}
