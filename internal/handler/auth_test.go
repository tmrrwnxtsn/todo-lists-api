package handler

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/model"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/service"
	mockservice "github.com/tmrrwnxtsn/todo-lists-api/internal/service/mocks"
	"net/http/httptest"
	"testing"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mockservice.MockAuthorization, user model.User)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            model.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "ok",
			inputBody: `{"name":"Pavel","username":"tmrrwnxtsn","password":"qwerty"}`,
			inputUser: model.User{
				Name:     "Pavel",
				Username: "tmrrwnxtsn",
				Password: "qwerty",
			},
			mockBehavior: func(s *mockservice.MockAuthorization, user model.User) {
				s.EXPECT().CreateUser(user).Return(uint64(1), nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:                 "empty fields",
			inputBody:            `{"username":"tmrrwnxtsn","password":"qwerty"}`,
			mockBehavior:         func(s *mockservice.MockAuthorization, user model.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "service failure",
			inputBody: `{"name":"Pavel","username":"tmrrwnxtsn","password":"qwerty"}`,
			inputUser: model.User{
				Name:     "Pavel",
				Username: "tmrrwnxtsn",
				Password: "qwerty",
			},
			mockBehavior: func(s *mockservice.MockAuthorization, user model.User) {
				s.EXPECT().CreateUser(user).Return(uint64(0), errors.New("service failure"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"service failure"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mockservice.NewMockAuthorization(c)
			tt.mockBehavior(auth, tt.inputUser)

			services := &service.Service{AuthService: auth}
			handler := NewHandler(services)

			router := gin.New()
			router.POST("/sign-up", handler.signUp)

			responseRecorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				"POST",
				"/sign-up",
				bytes.NewBufferString(tt.inputBody),
			)

			router.ServeHTTP(responseRecorder, request)

			assert.Equal(t, tt.expectedStatusCode, responseRecorder.Code)
			assert.Equal(t, tt.expectedResponseBody, responseRecorder.Body.String())
		})
	}
}

func TestHandler_signIn(t *testing.T) {
	type mockBehavior func(s *mockservice.MockAuthorization, req signInRequest)

	tests := []struct {
		name                 string
		inputBody            string
		inputCreds           signInRequest
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "ok",
			inputBody: `{"username":"tmrrwnxtsn","password":"qwerty"}`,
			inputCreds: signInRequest{
				Username: "tmrrwnxtsn",
				Password: "qwerty",
			},
			mockBehavior: func(s *mockservice.MockAuthorization, creds signInRequest) {
				s.EXPECT().GenerateToken(creds.Username, creds.Password).Return("token", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"token":"token"}`,
		},
		{
			name:                 "empty fields",
			inputBody:            `{"username":"tmrrwnxtsn"}`,
			mockBehavior:         func(s *mockservice.MockAuthorization, creds signInRequest) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "service failure",
			inputBody: `{"username":"tmrrwnxtsn","password":"qwerty"}`,
			inputCreds: signInRequest{
				Username: "tmrrwnxtsn",
				Password: "qwerty",
			},
			mockBehavior: func(s *mockservice.MockAuthorization, creds signInRequest) {
				s.EXPECT().GenerateToken(creds.Username, creds.Password).Return("", errors.New("service failure"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"service failure"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mockservice.NewMockAuthorization(c)
			tt.mockBehavior(auth, tt.inputCreds)

			services := &service.Service{AuthService: auth}
			handler := NewHandler(services)

			router := gin.New()
			router.POST("/sign-in", handler.signIn)

			responseRecorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				"POST",
				"/sign-in",
				bytes.NewBufferString(tt.inputBody),
			)

			router.ServeHTTP(responseRecorder, request)

			assert.Equal(t, tt.expectedStatusCode, responseRecorder.Code)
			assert.Equal(t, tt.expectedResponseBody, responseRecorder.Body.String())
		})
	}
}
