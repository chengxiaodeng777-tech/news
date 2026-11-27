package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type SysDictType struct {
	ent.Schema
}

func (SysDictType) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_dict_type"},
	}
}

func (SysDictType) Mixin() []ent.Mixin {
	return []ent.Mixin{
		AuditMixin{},
	}
}

func (SysDictType) Fields() []ent.Field {
	return []ent.Field{
		field.String("code").
			Unique().
			NotEmpty(),
		field.String("name").
			NotEmpty(),
		field.Int8("status").
			Default(1),
		field.String("remark").
			Default(""),
	}
}

func (SysDictType) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("dict_data", SysDictData.Type),
	}
}
