package handler

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
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
