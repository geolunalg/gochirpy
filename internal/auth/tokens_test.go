package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	const secret = "secret"
	validToken, _ := MakeJWT(userID, secret, time.Hour)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: secret,
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: secret,
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	const jwt = "jwt_encrypted_token"

	tests := []struct {
		name    string
		headers http.Header
		token   string
		wantErr bool
	}{
		{
			name:    "valid token",
			headers: http.Header{"Authorization": {"Bearer " + jwt}},
			token:   jwt,
			wantErr: false,
		},
		{
			name:    "invalid token",
			headers: http.Header{"Authorization": {"not_a_valid_jwt"}},
			token:   "",
			wantErr: true,
		},
		{
			name:    "no auth header",
			headers: http.Header{"Content-Type": {"application/json"}},
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authToken, err := GetBearerToken(tt.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if authToken != tt.token {
				t.Errorf("GetBearerToken() gotUserID = %v, want %v", authToken, tt.token)
			}
		})
	}
}
