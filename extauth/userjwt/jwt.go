package userjwt

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

const (
	// jwtTTL defines how long a JWT string should be considered valid
	jwtTTL = time.Second * 30
)

var (
	// privateKey holds the RSA 4096-bit private key used to sign our JWTs
	privateKey *rsa.PrivateKey
)

// init is a load time initialization function.
func init() {

	// Load the RSA 4096-bit pem file
	keyBytes, err := ioutil.ReadFile("rsa.pem")
	if err != nil {
		log.Fatalf("failed to load RSE pem file: %v", err)
	}

	// Parse the pem bytes into an RSA key
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
	if err != nil {
		log.Fatalf("failed to parse RSE private key: %v", err)
	}
}

// CreateJWT builds and signs a crude JWT that wraps a username
func CreateJWT(username string) (string, error) {

	// The time we will use for "iat" and "nbf"
	now := time.Now().UTC()

	// When will the thing expire
	expires := now.Add(jwtTTL)

	// Define the claims that the JWy=T will contain
	claims := make(jwt.MapClaims)
	claims["sub"] = username       // The subject, i.e. user that this JWT represents
	claims["exp"] = expires.Unix() // The expiration time after which the token must be disregarded.
	claims["iat"] = now.Unix()     // The time at which the token was issued.
	claims["nbf"] = now.Unix()     // The time before which the token must be disregarded.

	// Create and sign the JWT
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}
