package common

import (
	"crypto/rsa"
	"fmt"
	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"strconv"
	"time"
)

func SignGlobalJwt() string {

	var (
		signKey *rsa.PrivateKey
		err     error
	)

	exp, _ := strconv.Atoi(Config.GlobalJwtExp)
	expire := time.Unix(time.Now().Add(time.Duration(exp)*time.Second).Unix(), 0)

	now := time.Unix(time.Now().Unix(), 0)

	c := jws.Claims{}
	c.SetIssuedAt(now)
	c.SetExpiration(expire)
	c.Set("client_type", 2)
	c.SetIssuer(Config.GlobalJwtISS)

	// 私钥格式化
	privateKey := fmt.Sprintf(`-----BEGIN PRIVATE KEY-----
%s
-----END PRIVATE KEY-----`, Config.GlobalJwtPrivateKey)

	signKey, err = crypto.ParseRSAPrivateKeyFromPEM([]byte(privateKey))

	if err != nil {
		Log.Error("error to parse private key:", err)
	}

	token := jws.NewJWT(c, crypto.SigningMethodRS256)

	serializedToken, err := token.Serialize(signKey)

	if err != nil {
		Log.Error("error to sign public key:", err)
	}

	return string(serializedToken)
}
