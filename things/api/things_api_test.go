package api

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/pseudoincorrect/bariot/things/mock"
	"github.com/pseudoincorrect/bariot/things/models"
)

func TestThingGetEndpoint(t *testing.T) {

	mockThing := mock.CreateMockThing()

	testData := []struct {
		name  string
		id    string
		thing *models.Thing
		want  string
		throw string
	}{
		{
			name:  "Return a validation error",
			id:    mockThing.Id[:len(mockThing.Id)-1],
			thing: nil,
			want:  "validation error\n",
			throw: "",
		},
		{
			name:  "Return a not found (thing) error",
			id:    mockThing.Id,
			thing: nil,
			want:  "thing not found\n",
			throw: "",
		},
		{
			name:  "Return an internal server error",
			id:    mockThing.Id,
			thing: &mockThing,
			want:  "SaveThingError\n",
			throw: "SaveThingError",
		},
		{
			name:  "OK, return a thing model",
			id:    mockThing.Id,
			thing: &mockThing,
			want:  mockThing.JsonString(),
			throw: "",
		},
	}

	for _, d := range testData {
		t.Run(d.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/{id}", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", d.id)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			s := mock.New()
			s.ThrowErr = d.throw
			if d.thing != nil {
				s.Thing = d.thing
			}

			thingGetEndpoint(&s)(w, r)

			got := w.Body.String()
			want := d.want
			if got != want {
				t.Errorf("got %q, want %q", got, want)
			}
		})

	}

}
