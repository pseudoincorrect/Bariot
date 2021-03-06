package service

import (
	"context"
	"testing"

	"github.com/pseudoincorrect/bariot/internal/things/models"
	e "github.com/pseudoincorrect/bariot/pkg/utils/errors"
	"github.com/pseudoincorrect/bariot/tests/mocks/caches"
	"github.com/pseudoincorrect/bariot/tests/mocks/clients"
	"github.com/pseudoincorrect/bariot/tests/mocks/repos"
	assert "github.com/stretchr/testify/require"
)

func createService() (Things, *repos.ThingsRepoMock, *clients.AuthClientMock, *caches.CacheMock) {
	repo := repos.NewRepoMock()
	auth := clients.NewAuthClientMock()
	cach := caches.NewCacheMock()
	srvc := New(&repo, &auth, &cach)
	return srvc, &repo, &auth, &cach
}

func CreateMockThing(userId string) models.Thing {
	return models.Thing{
		Id:        "00000000-0000-0000-0000-000000000001",
		CreatedAt: "2022-06-01T14:35:40+03:00",
		Key:       "123456789",
		Name:      "Thing_1",
		UserId:    userId,
		Metadata:  models.Metadata{"unit": "temperature"},
	}
}

func TestSaveThing(t *testing.T) {
	thing := CreateMockThing("123456789")
	ts, repo, _, _ := createService()

	t.Run("test repo failure", func(t *testing.T) {
		repo.ThrowErr = "repository failure"
		ctx := context.Background()
		err := ts.SaveThing(ctx, &thing)
		assert.NotNil(t, err, "It should throw an error")
		got := err.Error()
		want := repo.ThrowErr
		assert.Equal(t, got, want, "It should throw ")
	})

	t.Run("test repo failure", func(t *testing.T) {
		repo.ThrowErr = ""
		ctx := context.Background()
		err := ts.SaveThing(ctx, &thing)
		assert.Nil(t, err, "It should not throw an error")
		assert.NotNil(t, thing, "It should return a Thing")
		assert.Equal(t, thing, thing, "should be the same thing")
	})
}

func TestUserOfThingOrAdminRefactor(t *testing.T) {
	userId := "123456"
	wrongUserThing := CreateMockThing(userId[:len(userId)-1])
	thing := CreateMockThing(userId)

	testData := []struct {
		name      string
		repoThrow string
		authThrow string
		userId    string
		thingId   string
		role      string
		thing     *models.Thing
		expUserId string
		expError  error
	}{
		{
			name:      "Authentication failure",
			repoThrow: "",
			authThrow: "auth failure",
			userId:    "",
			thingId:   "",
			role:      "",
			thing:     nil,
			expUserId: "",
			expError:  e.ErrAuthn,
		},
		{
			name:      "OK test role admin",
			repoThrow: "",
			authThrow: "",
			userId:    userId,
			thingId:   "",
			role:      clients.AdminRole,
			thing:     nil,
			expUserId: userId,
			expError:  nil,
		},
		{
			name:      "Wrong role ",
			repoThrow: "",
			authThrow: "",
			userId:    userId,
			thingId:   "",
			role:      clients.ThingRole,
			thing:     nil,
			expUserId: "",
			expError:  e.ErrAuthz,
		},
		{
			name:      "Repository failure",
			repoThrow: "repository failure",
			authThrow: "",
			userId:    userId,
			thingId:   "",
			role:      clients.UserRole,
			thing:     nil,
			expUserId: "",
			expError:  e.ErrDb,
		},
		{
			name:      "Thing not found",
			repoThrow: "",
			authThrow: "",
			userId:    userId,
			thingId:   "",
			role:      clients.UserRole,
			thing:     nil,
			expUserId: "",
			expError:  e.ErrNotFound,
		},
		{
			name:      "Wrong user id",
			repoThrow: "",
			authThrow: "",
			userId:    userId,
			thingId:   "",
			role:      clients.UserRole,
			thing:     &wrongUserThing,
			expUserId: "",
			expError:  e.ErrAuthz,
		},
		{
			name:      "All OK",
			repoThrow: "",
			authThrow: "",
			userId:    userId,
			thingId:   "",
			role:      clients.UserRole,
			thing:     &thing,
			expUserId: userId,
			expError:  nil,
		},
	}
	for _, d := range testData {
		t.Run(d.name, func(t *testing.T) {
			ts, repo, auth, _ := createService()
			auth.ThrowErr = d.authThrow
			auth.UserId = d.userId
			auth.UserRole = d.role
			repo.ThrowErr = d.repoThrow
			repo.Thing = d.thing
			token := "it does not matter"
			thingId := "it does not matter too"

			ctx := context.Background()
			userId, err := ts.UserOfThingOrAdmin(ctx, token, thingId)

			assert.Equal(t, d.expError, err, "It should throw an error")
			if err == nil {
				assert.Equal(t, d.expUserId, userId, "It should be the same User")
			}
		})
	}
}

func TestGetUserOfThing(t *testing.T) {
	thingId := "000.000.001"
	userId := "000.000.002"
	thing := CreateMockThing(userId)
	ts, repo, _, _ := createService()
	repo.Thing = &thing

	t.Run("test repo failure", func(t *testing.T) {
		repo.ThrowErr = "repository failure"
		ctx := context.Background()
		userId, err := ts.GetUserOfThing(ctx, thingId)
		assert.NotNil(t, err, "It should throw an error")
		assert.Equal(t, "", userId, "user id should be empty")
	})

	t.Run("all pass", func(t *testing.T) {
		repo.ThrowErr = ""
		ctx := context.Background()
		res, err := ts.GetUserOfThing(ctx, thingId)
		assert.Nil(t, err, "It should not throw an error")
		assert.Equal(t, userId, res, "should be the same thing")
	})
}
