package schema

import (
	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/donech/pearl/pkg/ent/mixin"
)

// Supplier holds the schema definition for the Supplier entity.
type Supplier struct {
	ent.Schema
}

// Fields of the Supplier.
func (Supplier) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Default("").NotEmpty().Annotations(
			entproto.Field(2),
		),
		field.String("address").Default("").NotEmpty().Annotations(
			entproto.Field(3),
		),
	}
}

// Edges of the Supplier.
func (Supplier) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("contacts", SupplierContact.Type),
	}
}

// Annotations of the Supplier.
func (Supplier) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "supplier",
			Collation: "utf8mb4_general_ci",
		},
		entproto.Message(),
	}
}

func (Supplier) Mixin() []ent.Mixin {
	return []ent.Mixin{mixin.CUDTime{}}
}
