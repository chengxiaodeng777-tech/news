package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type Article struct{ ent.Schema }

func (Article) Annotations() []schema.Annotation {
	return []schema.Annotation{entsql.Annotation{Table: "article"}}
}
func (Article) Mixin() []ent.Mixin { return []ent.Mixin{AuditMixin{}} }
func (Article) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").NotEmpty(),
		field.String("category").Default("").MaxLen(80),
		field.String("summary").Default(""),
		field.Text("content"),
		field.String("cover").Default(""),
		field.Int8("status").Default(1),
		field.String("author").Default(""),
		field.Int("sort").Default(0).NonNegative(),
	}
}
