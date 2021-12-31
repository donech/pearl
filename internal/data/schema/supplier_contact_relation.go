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
type SupplierContactRelation struct {
	ent.Schema
}

// Fields of the SupplierContact.
func (SupplierContactRelation) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("contact_id").Optional().Nillable().Annotations(
			entproto.Field(2),
		),
		field.Int8("type").Default(0).Annotations(
			entproto.Field(3),
		),
		field.String("value").Default("").NotEmpty().Annotations(
			entproto.Field(4),
		),
	}
}

// Edges of the SupplierContact.
func (SupplierContactRelation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("contact", SupplierContact.Type).
			Ref("relations").
			Unique().Field("contact_id"),
	}
}

// Annotations of the SupplierContact.
func (SupplierContactRelation) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "supplier_contact_relation",
			Collation: "utf8mb4_general_ci",
		},
		entproto.Message(),
	}
}

func (SupplierContactRelation) Mixin() []ent.Mixin {
	return []ent.Mixin{mixin.CUDTime{}}
}
