package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const chirpy = "chirpy"

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	issuedAt := time.Now().UTC()
	expiresAt := issuedAt.Add(expiresIn)

	claims := jwt.RegisteredClaims{
		Issuer:    chirpy,
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return ss, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claimsStruct, func(t *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != chirpy {
		return uuid.Nil, err
	}

	userId, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return uuid.Nil, err
	}

	return userUUID, nil
}
