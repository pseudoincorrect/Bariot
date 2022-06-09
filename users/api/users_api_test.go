package api

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	appErr "github.com/pseudoincorrect/bariot/pkg/errors"
	serviceMock "github.com/pseudoincorrect/bariot/users/mock/service"
)

func sendGetUserById(m *serviceMock.Mock, id string) string {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/{id}", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id)
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	userGetEndpoint(m)(w, r)
	return w.Body.String()

}

func TestThingGetEndpoint(t *testing.T) {
	userId := "00000000-0000-0000-0000-000000000001"
	m := new(serviceMock.Mock)
	// Case all fine
	userName := "John"
	m.On("GetUser", userId).Return(nil, userId, userName).Once()
	res := sendGetUserById(m, userId)
	if !strings.Contains(res, userId) {
		t.Fatal("Should return a User model")
	}
	// Case user not found
	m.On("GetUser", userId).Return(nil, "", "").Once()
	res = sendGetUserById(m, userId)
	if !strings.Contains(res, appErr.ErrUserNotFound.Error()) {
		t.Fatal("Should return an error")
	}
	// Case repository error
	m.On("GetUser", userId).Return(appErr.ErrDb, "", "").Once()
	res = sendGetUserById(m, userId)
	if !strings.Contains(res, appErr.ErrDb.Error()) {
		t.Fatal("Should return an error")
	}
}
