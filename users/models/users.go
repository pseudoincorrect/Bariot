package models

import (
	"context"
	"fmt"
)

type Metadata map[string]interface{}

// create table things (
//   id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
//   created_at timestamp with time zone DEFAULT now(),
//   email character varying(255) NOT NULL,
//   full_name character varying(255) NOT NULL,
//   metadata json
// );

type User struct {
	Id        string   `json:"Id"`
	CreatedAt string   `json:"CreatedAt"`
	Email     string   `json:"Email"`
	FullName  string   `json:"FullName"`
	HashPass  string   `json:"HashPass"`
	Metadata  Metadata `json:"Metadata"`
}

func (t User) String() string {
	return fmt.Sprintf("User{\n  Id: %s,\n  CreatedAt: %s,\n  Email: %s,\n  FullName: %s,\n  HashPass: %s,\n Metadata: %v\n}", t.Id, t.CreatedAt, t.Email, t.FullName, t.HashPass, t.Metadata)
}

type UsersRepository interface {
	Save(context.Context, *User) (*User, error)
	Get(context.Context, string) (*User, error)
	Delete(context.Context, string) (string, error)
	Update(context.Context, *User) (*User, error)
	GetByEmail(context.Context, string) (*User, error)
}
