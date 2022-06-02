package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/service"
	mockservice "github.com/tmrrwnxtsn/todo-lists-api/internal/service/mocks"
	"net/http/httptest"
	"testing"
)

func TestHandler_identifyUser(t *testing.T) {
	type mockBehavior func(s *mockservice.MockAuthorization, token string)

	tests := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "ok",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *mockservice.MockAuthorization, token string) {
				r.EXPECT().ParseToken(token).Return(uint64(1), nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "1",
		},
		{
			name:                 "invalid header name",
			headerName:           "",
			headerValue:          "Bearer token",
			token:                "token",
			mockBehavior:         func(r *mockservice.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"empty auth header"}`,
		},
		{
			name:                 "invalid header parts length",
			headerName:           "Authorization",
			headerValue:          "Bearr token token",
			token:                "token",
			mockBehavior:         func(r *mockservice.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"invalid auth header"}`,
		},
		{
			name:                 "invalid header value",
			headerName:           "Authorization",
			headerValue:          "Bearr token",
			token:                "token",
			mockBehavior:         func(r *mockservice.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"invalid auth header"}`,
		},
		{
			name:                 "empty token",
			headerName:           "Authorization",
			headerValue:          "Bearer ",
			token:                "token",
			mockBehavior:         func(r *mockservice.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"token is empty"}`,
		},
		{
			name:        "parse error",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *mockservice.MockAuthorization, token string) {
				r.EXPECT().ParseToken(token).Return(uint64(0), errors.New("invalid token"))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"invalid token"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mockservice.NewMockAuthorization(c)
			tt.mockBehavior(auth, tt.token)

			services := &service.Service{AuthService: auth}
			handler := NewHandler(services)

			router := gin.New()
			router.POST("/protected", handler.identifyUser, func(context *gin.Context) {
				id, _ := context.Get(userCtx)
				context.String(200, fmt.Sprintf("%d", id.(uint64)))
			})

			responseRecorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				"POST",
				"/protected",
				nil,
			)
			request.Header.Set(tt.headerName, tt.headerValue)

			router.ServeHTTP(responseRecorder, request)

			assert.Equal(t, tt.expectedStatusCode, responseRecorder.Code)
			assert.Equal(t, tt.expectedResponseBody, responseRecorder.Body.String())
		})
	}
}
