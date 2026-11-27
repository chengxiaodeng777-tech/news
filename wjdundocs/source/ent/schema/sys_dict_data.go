package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type SysDictData struct {
	ent.Schema
}

func (SysDictData) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_dict_data"},
	}
}

func (SysDictData) Mixin() []ent.Mixin {
	return []ent.Mixin{
		AuditMixin{},
	}
}

func (SysDictData) Fields() []ent.Field {
	return []ent.Field{
		field.Int("dict_type_id"),
		field.String("label").
			NotEmpty(),
		field.String("value").
			NotEmpty(),
		field.String("tag_type").
			Default(""),
		field.String("css_class").
			Default(""),
		field.Int("sort").
			Default(0),
		field.Int8("status").
			Default(1),
		field.Bool("is_default").
			Default(false),
		field.String("remark").
			Default(""),
	}
}

func (SysDictData) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("dict_type", SysDictType.Type).
			Ref("dict_data").
			Field("dict_type_id").
			Unique().
			Required(),
	}
}
