package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// 你现有的函数：保持签名不变，向后兼容
// @secretKey: JWT 加解密密钥
// @iat: 时间戳
// @seconds: 过期时间，单位秒
// @UserId: 数据载体（用户ID）
func GetJwtToken(secretKey string, iat, seconds int64, UserId int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["UserId"] = UserId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

// 新增：带额外 claims（typ/jti），内部仍用相同签名方式
// typ: "access" | "refresh"; jti: 仅 refresh 需要，access 传空
func GetJwtTokenWithClaims(secretKey string, iat, seconds int64, userId int64, typ, jti string) (string, error) {
	claims := jwt.MapClaims{
		"exp":    iat + seconds,
		"iat":    iat,
		"UserId": userId,
		"typ":    typ, // 标识 token 类型
	}
	if typ == "refresh" && jti != "" {
		claims["jti"] = jti
	}
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

// 解析并校验（签名/exp），返回 MapClaims
func ParseJwt(tokenStr, secret string) (jwt.MapClaims, error) {
	t, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// 只接受 HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	// 额外校验 exp
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, errors.New("token expired")
		}
	}
	return claims, nil
}

// 便捷取值
func ClaimInt64(m jwt.MapClaims, key string) int64 {
	if v, ok := m[key].(float64); ok {
		return int64(v)
	}
	return 0
}
func ClaimString(m jwt.MapClaims, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}
