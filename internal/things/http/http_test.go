package http

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/pseudoincorrect/bariot/internal/things/models"
	"github.com/pseudoincorrect/bariot/pkg/utils/errors"
	"github.com/pseudoincorrect/bariot/tests/mocks/services"
)

func CreateMockThing() models.Thing {
	return models.Thing{
		Id:        "00000000-0000-0000-0000-000000000001",
		CreatedAt: "2022-06-01T14:35:40+03:00",
		Key:       "123456789",
		Name:      "Thing_1",
		UserId:    "00000000-0000-0000-0000-000000000002",
		Metadata:  models.Metadata{"unit": "temperature"},
	}
}

func TestThingGetEndpoint(t *testing.T) {
	mockThing := CreateMockThing()
	testData := []struct {
		name    string
		thingId string
		thing   *models.Thing
		want    string
		throw   string
	}{
		{
			name:    "Return a validation error",
			thingId: mockThing.Id[:len(mockThing.Id)-1],
			thing:   nil,
			want:    errors.ErrValidation.Error() + "\n",
			throw:   "",
		},
		{
			name:    "Return a not found (thing) error",
			thingId: mockThing.Id,
			thing:   nil,
			want:    errors.ErrNotFound.Error() + "\n",
			throw:   "",
		},
		{
			name:    "Return an internal server error",
			thingId: mockThing.Id,
			thing:   &mockThing,
			want:    "SaveThingError" + "\n",
			throw:   "SaveThingError",
		},
		{
			name:    "OK, return a thing model",
			thingId: mockThing.Id,
			thing:   &mockThing,
			want:    mockThing.JsonString(),
			throw:   "",
		},
	}
	for _, d := range testData {
		t.Run(d.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/{id}", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", d.thingId)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			s := services.NewThingsMock()
			s.ThrowErr = d.throw
			if d.thing != nil {
				s.Thing = d.thing
			}
			thingGetEndpoint(&s)(w, r)
			got := w.Body.String()
			want := d.want
			if got != want {
				t.Errorf("RES: got %q, want %q", got, want)
			}
		})
	}
}

func TestThingPostEndpoint(t *testing.T) {
	name := "SmartDispenser1"
	bodyMalformed := []byte(`{"Name":"SmartDispenser1,"Key":"1234"}`)
	bodyWrongArgs := []byte(`{"Name":"Ma","Key":"1234"}`)
	bodyCorrect := []byte(`{"Name":"` + name + `","Key":"1234"}`)
	testData := []struct {
		name   string
		userId string
		want   string
		throw  string
		body   []byte
	}{
		{
			name:   "Incorrect Body formatting (json decode)",
			userId: "123456789",
			want:   errors.ErrParsing.Error(),
			throw:  "",
			body:   bodyMalformed,
		},
		{
			name:   "Incorrect body argument",
			userId: "123456789",
			want:   errors.ErrValidation.Error(),
			throw:  "",
			body:   bodyWrongArgs,
		},
		{
			name:   "Thing service save error",
			userId: "123456789",
			want:   "SaveThingError",
			throw:  "SaveThingError",
			body:   bodyCorrect,
		},
		{
			name:   "OK, save a thing",
			userId: "123456789",
			want:   name,
			throw:  "",
			body:   bodyCorrect,
		},
	}
	for _, d := range testData {
		t.Run(d.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/things/", bytes.NewReader(d.body))
			rctx := chi.NewRouteContext()
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			r = r.WithContext(context.WithValue(r.Context(), userIdKey, d.userId))
			s := services.NewThingsMock()
			s.ThrowErr = d.throw
			thingPostEndpoint(&s)(w, r)
			got := w.Body.String()
			if !strings.Contains(got, d.want) {
				t.Errorf("RES: got %q, want %q", got, d.want)
			}
		})
	}
}

func TestThingDeleteEndpoint(t *testing.T) {
	mockThing := CreateMockThing()
	testData := []struct {
		name    string
		thingId string
		want    string
		throw   string
	}{
		{
			name:    "Return a validation error",
			thingId: mockThing.Id[:len(mockThing.Id)-1],
			want:    errors.ErrValidation.Error() + "\n",
			throw:   "",
		},
		{
			name:    "Return a save thing error",
			thingId: mockThing.Id,
			want:    "delete error" + "\n",
			throw:   "delete error",
		},
		{
			name:    "OK, return a thing ID",
			thingId: mockThing.Id,
			want:    "{\"Id\":\"" + mockThing.Id + "\"}\n",
			throw:   "",
		},
	}
	for _, d := range testData {
		t.Run(d.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/{id}", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", d.thingId)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			s := services.NewThingsMock()
			s.ThrowErr = d.throw
			thingDeleteEndpoint(&s)(w, r)
			got := w.Body.String()
			want := d.want
			if got != want {
				t.Errorf("RES: got %q, want %q", got, want)
			}
		})
	}
}

func TestThingGetTokenEndpoint(t *testing.T) {
	mockThing := CreateMockThing()
	mockToken := "123.123.123"
	userId := mockThing.UserId
	testData := []struct {
		name    string
		thingId string
		want    string
		throw   string
		token   string
		userId  string
	}{
		{
			name:    "Return a Wrong JWT error",
			thingId: mockThing.Id,
			want:    "error jwt" + "\n",
			throw:   "error jwt",
			token:   mockToken,
			userId:  userId,
		},
		{
			name:    "OK, return a JWT",
			thingId: mockThing.Id,
			want:    "{\"Token\":\"" + mockToken + "\"}\n",
			throw:   "",
			token:   mockToken,
			userId:  userId,
		},
	}
	for _, d := range testData {
		t.Run(d.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/{id}", nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", d.thingId)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			r = r.WithContext(context.WithValue(r.Context(), userIdKey, d.userId))

			s := services.NewThingsMock()
			s.ThrowErr = d.throw
			s.Token = d.token

			thingGetTokenEndpoint(&s)(w, r)

			got := w.Body.String()
			want := d.want
			if got != want {
				t.Errorf("RES: got %q, want %q", got, want)
			}
		})
	}
}

func TestUserOnly(t *testing.T) {
	mockThing := CreateMockThing()
	mockToken := "123.123.123"
	mockThingId := mockThing.Id
	mockUserId := mockThing.UserId
	testData := []struct {
		name    string
		thingId string
		userId  string
		resWant string
		throw   string
	}{
		{
			name:    "Return a Wrong JWT error",
			thingId: mockThingId,
			resWant: "wrong jwt" + "\n",
			throw:   "wrong jwt",
			userId:  "",
		},
		{
			name:    "OK, pass the correct userId to the next handler",
			thingId: mockThingId,
			resWant: "",
			throw:   "",
			userId:  mockUserId,
		},
	}
	for _, d := range testData {
		t.Run(d.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/{id}", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", d.thingId)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			r.Header.Set("Authorization", mockToken)
			s := services.NewThingsMock()
			s.ThrowErr = d.throw
			s.UserId = d.userId
			var userId string
			next := func(res http.ResponseWriter, req *http.Request) {
				userId = req.Context().Value(userIdKey).(string)
			}
			userOnly(&s, http.HandlerFunc(next))(w, r)
			got := w.Body.String()
			resWant := d.resWant
			if userId != d.userId {
				t.Errorf("UserId: got %q, resWant %q", userId, d.userId)
			}
			if got != resWant {
				t.Errorf("RES: got %q, want %q", got, resWant)
			}
		})
	}
}

func TestUserOfThingOnly(t *testing.T) {
	mockThing := CreateMockThing()
	mockToken := "123.123.123"
	mockThingId := mockThing.Id
	mockThingIdMalformed := mockThing.Id[:len(mockThingId)-1]
	mockUserId := mockThing.UserId
	testData := []struct {
		name    string
		thingId string
		userId  string
		resWant string
		throw   string
	}{
		{
			name:    "Thing Id validation error",
			thingId: mockThingIdMalformed,
			userId:  "",
			resWant: errors.ErrValidation.Error() + "\n",
			throw:   "",
		},
		{
			name:    "Wrong token error",
			thingId: mockThingId,
			userId:  "",
			resWant: "wrong token" + "\n",
			throw:   "wrong token",
		},
		{
			name:    "OK, pass the user id to the next handler",
			thingId: mockThingId,
			userId:  mockUserId,
			resWant: "",
			throw:   "",
		},
	}
	for _, d := range testData {
		t.Run(d.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/{id}", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", d.thingId)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			r.Header.Set("Authorization", mockToken)
			s := services.NewThingsMock()
			s.ThrowErr = d.throw
			s.UserId = d.userId
			var userId string
			next := func(res http.ResponseWriter, req *http.Request) {
				userId = req.Context().Value(userIdKey).(string)
			}

			userOfThingOnly(&s, http.HandlerFunc(next))(w, r)

			got := w.Body.String()
			resWant := d.resWant
			if userId != d.userId {
				t.Errorf("UserId: got %q, resWant %q", userId, d.userId)
			}
			if got != resWant {
				t.Errorf("RES: got %q, want %q", got, resWant)
			}
		})
	}
}

func TestUserOfThingOrAdmin(t *testing.T) {
	mockThing := CreateMockThing()
	mockToken := "123.123.123"
	mockThingId := mockThing.Id
	mockThingIdMalformed := mockThing.Id[:len(mockThingId)-1]
	mockUserId := mockThing.UserId
	testData := []struct {
		name    string
		thingId string
		userId  string
		resWant string
		throw   string
	}{
		{
			name:    "Thing Id validation error",
			thingId: mockThingIdMalformed,
			userId:  "",
			resWant: errors.ErrValidation.Error() + "\n",
			throw:   "",
		},
		{
			name:    "Wrong token error",
			thingId: mockThingId,
			userId:  "",
			resWant: "wrong token" + "\n",
			throw:   "wrong token",
		},
		{
			name:    "OK, pass the user id to the next handler",
			thingId: mockThingId,
			userId:  mockUserId,
			resWant: "",
			throw:   "",
		},
	}
	for _, d := range testData {
		t.Run(d.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/{id}", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", d.thingId)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			r.Header.Set("Authorization", mockToken)
			s := services.NewThingsMock()
			s.ThrowErr = d.throw
			s.UserId = d.userId
			var userId string
			next := func(res http.ResponseWriter, req *http.Request) {
				userId = req.Context().Value(userIdKey).(string)
			}

			userOfThingOrAdmin(&s, http.HandlerFunc(next))(w, r)

			got := w.Body.String()
			resWant := d.resWant
			if userId != d.userId {
				t.Errorf("UserId: got %q, resWant %q", userId, d.userId)
			}
			if got != resWant {
				t.Errorf("RES: got %q, want %q", got, resWant)
			}
		})
	}
}
