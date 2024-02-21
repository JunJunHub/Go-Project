package token

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcache"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/guid"

	"encoding/base64"
	"time"
)

type MpuClaims struct {
	//【JWT ID】     该jwt的唯一ID编号
	//【issuer】     发布者的url地址
	//【issued at】  该jwt的发布时间；unix 时间戳
	//【subject】    该JWT所面向的用户，用于处理特定应用，不是常用的字段
	//【audience】   接受者的url地址
	//【expiration】 该jwt销毁的时间；unix时间戳
	//【not before】 该jwt的使用时间不能早于该时间；unix时间戳
	StandardClaims *jwt.StandardClaims
	RefreshTime    int64  //【The refresh time】 该jwt刷新的时间；unix时间戳
	UserId         string `json:"user_id"`
	LoginName      string `json:"login_name"`
}

type Token struct {
	Claim    *MpuClaims
	Token    string
	NewToken string
}

const CacheKey = "GJWT"

//New 创建Claims
func New() *MpuClaims {
	timeOut := g.Cfg().GetInt("jwt.timeout")
	if timeOut <= 0 {
		timeOut = 3600
	}
	refresh := g.Cfg().GetInt("jwt.refresh")
	if refresh <= 0 {
		refresh = timeOut / 2
	}
	var claims MpuClaims
	standardClaims := new(jwt.StandardClaims)
	standardClaims.Id = guid.S()
	standardClaims.ExpiresAt = time.Now().Add(time.Second * time.Duration(timeOut)).Unix()
	standardClaims.IssuedAt = time.Now().Unix()
	claims.RefreshTime = time.Now().Add(time.Second * time.Duration(refresh)).Unix()
	claims.StandardClaims = standardClaims
	return &claims
}

func (c *MpuClaims) SetIss(issuer string) *MpuClaims {
	c.StandardClaims.Issuer = issuer
	return c
}

func (c *MpuClaims) SetSub(subject string) *MpuClaims {
	c.StandardClaims.Subject = subject
	return c
}

func (c *MpuClaims) SetAud(audience string) *MpuClaims {
	c.StandardClaims.Audience = audience
	return c
}

func (c *MpuClaims) SetNbf(notBefore int64) *MpuClaims {
	c.StandardClaims.NotBefore = notBefore
	return c
}
func (c *MpuClaims) SetUserId(userId string) *MpuClaims {
	c.UserId = userId
	return c
}
func (c *MpuClaims) SetLoginName(loginName string) *MpuClaims {
	c.LoginName = loginName
	return c
}

func (c *MpuClaims) Valid() error {
	//标准验证
	return c.StandardClaims.Valid()
}

//CreateToken 创建token
func (c *MpuClaims) CreateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	mySignKeyBytes, err := base64.URLEncoding.DecodeString(g.Cfg().GetString("api.jwt.encryptKey"))
	if err != nil {
		return "", err
	}
	c.SetCache(CacheKey+c.LoginName, c.UserId)
	return token.SignedString(mySignKeyBytes)
}

//VerifyAuthToken 验证token
func VerifyAuthToken(token string) (*Token, error) {
	var MpuClaims = new(MpuClaims)
	MpuClaims, err := MpuClaims.DecryptToken(token)
	if err != nil {
		return nil, err
	}
	// 从缓存获取
	userId, err := MpuClaims.GetCache(CacheKey + MpuClaims.LoginName)
	if err != nil {
		return nil, err
	}
	rs := new(Token)
	rs.Claim = MpuClaims
	rs.Token = token
	//判断是否需要刷新
	if MpuClaims.RefreshTime > time.Now().Unix() {
		//生成新token
		newToken, err := New().SetUserId(userId).SetLoginName(rs.Claim.LoginName).CreateToken()
		if err == nil {
			rs.NewToken = newToken
		}
	}
	return rs, nil
}

func (c *MpuClaims) DecryptToken(token string) (*MpuClaims, error) {
	mySignKey := g.Cfg().GetString("api.jwt.encryptKey")
	mySignKeyBytes, err := base64.URLEncoding.DecodeString(mySignKey) //需要用和加密时同样的方式转化成对应的字节数组
	if err != nil {
		return nil, err
	}
	parseAuth, err := jwt.ParseWithClaims(token, c, func(*jwt.Token) (interface{}, error) {
		return mySignKeyBytes, nil
	})
	if err != nil {
		return nil, err
	}
	//验证claims
	if err := parseAuth.Claims.Valid(); err != nil {
		return nil, err
	}
	return c, nil
}

//SetCache 设置缓存
func (c *MpuClaims) SetCache(cacheKey string, userId string) string {
	cacheType := g.Cfg().GetString("jwt.cache")
	switch cacheType {
	case "default":
		err := gcache.Set(cacheKey, userId, gconv.Duration(c.RefreshTime)*time.Millisecond*2)
		if err != nil {
			g.Log().Error(err)
		}
	case "redis":
		_, err := g.Redis().DoVar("SETEX", cacheKey, c.RefreshTime-time.Now().Unix(), userId)
		if err != nil {
			g.Log().Error(err)
		}
	}
	return userId
}

//GetCache 获取缓存
func (c *MpuClaims) GetCache(cacheKey string) (string, error) {

	cacheType := g.Cfg().GetString("jwt.cache")
	var userCacheValue interface{}
	var err error
	switch cacheType {
	case "default":
		userCacheValue, _ = gcache.Get(cacheKey)
	case "redis":
		userCacheValue, err = g.Redis().Do("GET", cacheKey)
		if err != nil {
			return "", gerror.New("请登录")
		}
	}
	if userCacheValue == nil {
		return "", gerror.New("请登录")
	}
	userId := gconv.String(userCacheValue)
	return userId, nil
}
func RemoveCache(cacheKey string) {
	cacheType := g.Cfg().GetString("jwt.cache")
	switch cacheType {
	case "default":
		_, err := gcache.Remove(cacheKey)
		if err != nil {
			g.Log().Error(err)

		}
	case "redis":
		_, err := g.Redis().DoVar("DEL", cacheKey)
		if err != nil {
			g.Log().Error(err)

		}
	}
}
