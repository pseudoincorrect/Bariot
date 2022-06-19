package http

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	e "github.com/pseudoincorrect/bariot/pkg/errors"
	"github.com/pseudoincorrect/bariot/tests/mocks/services"
)

func sendGetUserById(m *services.UsersMock, id string) string {
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
	m := services.NewUsersMock()
	// Case all fine
	userName := "John"
	m.On("GetUser", userId).Return(nil, userId, userName).Once()
	res := sendGetUserById(&m, userId)
	if !strings.Contains(res, userId) {
		t.Fatal("Should return a User model")
	}
	// Case user not found
	m.On("GetUser", userId).Return(nil, "", "").Once()
	res = sendGetUserById(&m, userId)
	if !strings.Contains(res, e.ErrNotFound.Error()) {
		t.Fatal("Should return an error")
	}
	// Case repository error
	m.On("GetUser", userId).Return(e.ErrDb, "", "").Once()
	res = sendGetUserById(&m, userId)
	if !strings.Contains(res, e.ErrDb.Error()) {
		t.Fatal("Should return an error")
	}
}
