package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type SysCasbinRule struct {
	ent.Schema
}

func (SysCasbinRule) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_casbin_rule"},
	}
}

func (SysCasbinRule) Mixin() []ent.Mixin {
	return []ent.Mixin{
		AuditMixin{},
	}
}

func (SysCasbinRule) Fields() []ent.Field {
	return []ent.Field{
		field.String("ptype").
			NotEmpty(),
		field.String("v0").
			Default(""),
		field.String("v1").
			Default(""),
		field.String("v2").
			Default(""),
		field.String("v3").
			Default(""),
		field.String("v4").
			Default(""),
		field.String("v5").
			Default(""),
	}
}

func (SysCasbinRule) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("ptype", "v0", "v1", "v2", "v3", "v4", "v5").
			Unique(),
	}
}
