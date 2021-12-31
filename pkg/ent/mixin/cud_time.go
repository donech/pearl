package mixin

import (
	"time"

	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// CUDTime composes create/update/delete time mixin.
type CUDTime struct{ mixin.Schema }

// Fields of the time mixin.
func (CUDTime) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_time").Default(time.Now).Immutable().SchemaType(map[string]string{
			dialect.MySQL: "datetime DEFAULT CURRENT_TIMESTAMP",
		}).Annotations(
			entproto.Field(101),
		),
		field.Time("updated_time").Default(time.Now).UpdateDefault(time.Now).SchemaType(map[string]string{
			dialect.MySQL: "datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP",
		}).Immutable().Annotations(
			entproto.Field(102),
		),
		field.Time("deleted_time").Optional().Nillable().Annotations(
			entproto.Field(103),
		),
	}
}

// time mixin must implement `Mixin` interface.
var _ ent.Mixin = (*CUDTime)(nil)
