package http_test

import (
	authhttp "github.com/Lucas-Onofre/financial-chat/chat-service/internal/auth/http"
	"github.com/Lucas-Onofre/financial-chat/chat-service/internal/auth/jwt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_AuthMiddleware(t *testing.T) {
	jwtService := jwt.NewJWTService("secret", time.Minute*5)
	validToken, err := jwtService.GenerateToken("user123")
	if err != nil {
		t.Fatalf("failed to generate valid token: %v", err)
	}

	tests := []struct {
		name               string
		authHeader         string
		expectedStatusCode int
	}{
		{
			name:               "Given no Authorization header, When request is made, Then respond with 401",
			authHeader:         "",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "Given malformed Authorization header, When request is made, Then respond with 401",
			authHeader:         "InvalidHeader",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "Given invalid token, When request is made, Then respond with 401",
			authHeader:         "Bearer invalid.token.here",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "Given valid token, When request is made, Then respond with 200",
			authHeader:         "Bearer " + validToken,
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authMw := authhttp.AuthMiddleware(jwtService)

			finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("success"))
			})

			handlerToTest := authMw(finalHandler)

			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			rr := httptest.NewRecorder()
			handlerToTest.ServeHTTP(rr, req)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}
