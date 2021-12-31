package schema

import (
	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/donech/pearl/pkg/ent/mixin"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Default("").NotEmpty().Annotations(
			entproto.Field(2),
		),
		field.String("account").Default("").NotEmpty().Annotations(
			entproto.Field(3),
		),
		field.String("password").Default("").NotEmpty().Annotations(
			entproto.Field(4),
		),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}

// Annotations of the User.
func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "user",
			Collation: "utf8mb4_general_ci",
		},
		entproto.Message(),
	}
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{mixin.CUDTime{}}
}
