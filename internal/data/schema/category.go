package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/donech/pearl/pkg/ent/mixin"
)

// Category holds the schema definition for the User entity.
type Category struct {
	ent.Schema
}

// Fields of the User.
func (Category) Fields() []ent.Field {
	return []ent.Field{
		field.Int("parent_id").Positive().Comment("父ID"),
		field.Int("creator_id").Positive().Comment("创建者ID"),
		field.String("name").Default("类目名称"),
		field.String("type").Default("类目类型"),
	}
}

// Edges of the User.
func (Category) Edges() []ent.Edge {
	return nil
}

func (Category) Mixin() []ent.Mixin {
	return []ent.Mixin{mixin.CUDTime{}}
}
