package cmn

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWT结构体
type JWT struct {
	secret []byte
}

// 创建JWT对象
func NewJWT(secret string) *JWT {
	return &JWT{secret: StringToBytes(secret)}
}

// 创建令牌（默认HS256算法）
func (j *JWT) CreateToken(mapKv MapString, exp time.Duration) (string, error) {
	claims := make(jwt.MapClaims)
	for k, v := range mapKv {
		claims[k] = v
	}
	claims["exp"] = time.Now().Add(exp).Unix() // 设定超时时间
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.secret)
}

// 续签令牌（复制原令牌后创建新令牌，原令牌可以是过期令牌）
func (j *JWT) ExpandToken(token string, exp time.Duration) (string, error) {
	m, err := j.Parse(token)
	if err != nil {
		return "", err
	}
	return j.CreateToken(m, exp)
}

// 判断令牌是否已过期（过期令牌不会返回error，令牌无效时将返回error）
func (j *JWT) IsExpired(token string) (bool, error) {
	tk, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的令牌: " + token)
		}
		return j.secret, nil
	}, jwt.WithoutClaimsValidation())

	if err != nil {
		return false, err
	}

	if !tk.Valid {
		return false, errors.New("无效的令牌: " + token)
	}

	mc := tk.Claims.(jwt.MapClaims)
	return Float64ToInt64(mc["exp"].(float64)) <= time.Now().Unix(), nil
}

// 解析令牌（过期令牌不会产生错误，返回值不包含"exp"属性）
func (j *JWT) Parse(token string) (MapString, error) {
	tk, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的令牌: " + token)
		}
		return j.secret, nil
	}, jwt.WithoutClaimsValidation())

	if err != nil {
		return nil, err
	}

	if !tk.Valid {
		return nil, errors.New("无效的令牌: " + token)
	}

	rs := NewMapString()
	for k, v := range tk.Claims.(jwt.MapClaims) {
		if k != "exp" {
			rs.Put(k, v.(string))
		}
	}

	return rs, nil
}

// 校验令牌（过期令牌会返回error，返回值不包含"exp"属性）
func (j *JWT) Validate(token string) (MapString, error) {
	tk, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的令牌: " + token)
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if !tk.Valid {
		return nil, errors.New("无效的令牌: " + token)
	}

	rs := NewMapString()
	for k, v := range tk.Claims.(jwt.MapClaims) {
		if k != "exp" {
			rs.Put(k, v.(string))
		}
	}

	return rs, nil
}
