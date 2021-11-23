package utils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/nbkit/mdf/log"
	"io/ioutil"
	"strings"
)

/**
请求通用类
*/
type TokenContext struct {
	*BizObject
}

func NewTokenContext() *TokenContext {
	return &TokenContext{BizObject: NewBizObject()}
}
func (s *TokenContext) UserID() string {
	return s.GetString("UserID")
}

func (s *TokenContext) EntID() string {
	return s.GetString("EntID")
}

func (s *TokenContext) OrgID() string {
	return s.GetString("OrgID")
}

func (s *TokenContext) ID() string {
	return s.GetString("ID")
}

func (s TokenContext) ToTokenString() string {
	claim := jwt.MapClaims{}
	data := s.Data()
	for k, v := range data {
		claim[k] = v
	}
	token := jwt.NewWithClaims(s.getSigningMethod(), claim)
	tokenString, err := token.SignedString(s.getPrivateKey())
	if err != nil {
		return ""
	}
	return "bearer " + tokenString
}
func (s TokenContext) getSigningMethod() jwt.SigningMethod {
	v := Config.GetValue("oauth.alg")
	switch v {
	case jwt.SigningMethodRS256.Name:
		return jwt.SigningMethodRS256
	case jwt.SigningMethodES256.Name:
		return jwt.SigningMethodES256
	default:
		return jwt.SigningMethodHS256
	}
	return jwt.SigningMethodHS256
}
func (s TokenContext) getPrivateKey() interface{} {
	valueKey := "oauth.privatekey.value"
	key := "oauth.privatekey"
	if v := Config.GetObject(valueKey); v != nil {
		return v
	}
	switch s.getSigningMethod() {
	case jwt.SigningMethodRS256:
		keyByte, err := ioutil.ReadFile(JoinCurrentPath(Config.GetValue(key)))
		if err != nil {
			log.ErrorD(err)
			break
		}
		if v, err := jwt.ParseRSAPrivateKeyFromPEM(keyByte); err != nil {
			log.ErrorD(err)
			break
		} else {
			Config.SetValue(valueKey, v)
		}
		break
	default:
		Config.SetValue(valueKey, Config.GetValue(key))
	}
	return Config.GetObject(valueKey)
}
func (s TokenContext) getPublicKey() interface{} {
	valueKey := "oauth.publickey.value"
	key := "oauth.publickey"
	if v := Config.GetObject(valueKey); v != nil {
		return v
	}
	switch s.getSigningMethod() {
	case jwt.SigningMethodRS256:
		keyByte, err := ioutil.ReadFile(JoinCurrentPath(Config.GetValue(key)))
		if err != nil {
			log.ErrorD(err)
			break
		}
		if v, err := jwt.ParseRSAPublicKeyFromPEM(keyByte); err != nil {
			log.ErrorD(err)
			break
		} else {
			Config.SetValue(valueKey, v)
		}
		break
	default:
		Config.SetValue(valueKey, Config.GetValue(key))
	}
	return Config.GetObject(valueKey)
}
func (s TokenContext) FromTokenString(token string) (*TokenContext, error) {
	ctx := NewTokenContext()
	tokenParts := strings.Split(token, " ")
	if len(tokenParts) == 2 && strings.ToLower(tokenParts[0]) == "bearer" {
		token = tokenParts[1]
	}
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return s.getPublicKey(), nil
	})
	if err != nil {
		return ctx, err
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		for k, v := range claims {
			ctx.Set(k, v)
		}
		ctx.Set("id", GUID())
	}
	return ctx, nil
}
