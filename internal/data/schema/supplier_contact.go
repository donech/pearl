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

// SupplierContact holds the schema definition for the SupplierContact entity.
type SupplierContact struct {
	ent.Schema
}

// Fields of the SupplierContact.
func (SupplierContact) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("supplier_id").Optional().Annotations(
			entproto.Field(2),
		),
		field.String("firstname").Default("").NotEmpty().Annotations(
			entproto.Field(3),
		),
		field.String("lastname").Default("").NotEmpty().Annotations(
			entproto.Field(4),
		),
		field.String("title").Default("").NotEmpty().Annotations(
			entproto.Field(5),
		),
	}
}

// Edges of the SupplierContact.
func (SupplierContact) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("supplier", Supplier.Type).
			Ref("contacts").
			Unique().Field("supplier_id"),
		edge.To("relations", SupplierContactRelation.Type),
	}
}

// Annotations of the SupplierContact.
func (SupplierContact) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "supplier_contact",
			Collation: "utf8mb4_general_ci",
		},
		entproto.Message(),
	}
}

func (SupplierContact) Mixin() []ent.Mixin {
	return []ent.Mixin{mixin.CUDTime{}}
}
