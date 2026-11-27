package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type SysOrganization struct {
	ent.Schema
}

func (SysOrganization) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_organization"},
	}
}

func (SysOrganization) Mixin() []ent.Mixin {
	return []ent.Mixin{
		AuditMixin{},
	}
}

func (SysOrganization) Fields() []ent.Field {
	return []ent.Field{
		field.Int("parent_id").
			Optional().
			Nillable(),
		field.String("name").
			NotEmpty(),
		field.String("code").
			Unique().
			NotEmpty(),
		field.String("type").
			Default("department"),
		field.String("leader").
			Default(""),
		field.String("phone").
			Default(""),
		field.String("email").
			Default(""),
		field.Int("sort").
			Default(0),
		field.Int8("status").
			Default(1),
		field.String("remark").
			Default(""),
	}
}

func (SysOrganization) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("children", SysOrganization.Type).
			From("parent").
			Field("parent_id").
			Unique(),
		edge.To("users", SysUser.Type),
	}
}

func (SysOrganization) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("parent_id"),
	}
}
