package models

import (
	"context"
	"fmt"
)

type Metadata map[string]interface{}

// type ThingI interface {
// 	String() string
// }

type Thing struct {
	Id        string   `json:"id"`
	CreatedAt string   `json:"createdAt"`
	Key       string   `json:"key"`
	Name      string   `json:"name"`
	UserId    string   `json:"userId"`
	Metadata  Metadata `json:"metadata"`
}

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
