package models

import (
	"context"
	"fmt"
)

type Metadata map[string]interface{}

type Thing struct {
	Id        string   `json:"Id"`
	CreatedAt string   `json:"CreatedAt"`
	Key       string   `json:"Key"`
	Name      string   `json:"Name"`
	UserId    string   `json:"UserId"`
	Metadata  Metadata `json:"Metadata"`
}

// String returns a string representation of the thing
func (t Thing) String() string {
	return fmt.Sprintf("Thing{\n  Id: %s,\n  CreatedAt: %s,\n  Key: %s,\n  Name: %s,\n  UserId: %s,\n  Metadata: %v\n}", t.Id, t.CreatedAt, t.Key, t.Name, t.UserId, t.Metadata)
}

// POSTGRE table things (
//   id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
//   created_at timestamp with time zone DEFAULT now(),
//   key character varying(4096) NOT NULL,
//   name character varying(255) NOT NULL,
//   user_id uuid NOT NULL,
//   metadata json NOT NULL
// );

type ThingsRepository interface {
	Save(context.Context, *Thing) (*Thing, error)
	Get(context.Context, string) (*Thing, error)
	Delete(context.Context, string) (string, error)
	Update(context.Context, *Thing) (*Thing, error)
}
