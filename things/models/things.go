package models

import "context"

type Metadata map[string]interface{}

type Thing struct {
	Id        string
	CreatedAt string
	Key       string
	Name      string
	UserId    string
	// Metadata  Metadata
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
	Delete(context.Context, string) (*Thing, error)
	Update(context.Context, string, *Thing) error
}
