package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type SysRole struct {
	ent.Schema
}

func (SysRole) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_role"},
	}
}

func (SysRole) Mixin() []ent.Mixin {
	return []ent.Mixin{
		AuditMixin{},
	}
}

func (SysRole) Fields() []ent.Field {
	return []ent.Field{
		field.String("code").
			Unique().
			NotEmpty(),
		field.String("name").
			NotEmpty(),
		field.Int8("status").
			Default(1),
		field.Int("sort").
			Default(0),
		field.String("remark").
			Default(""),
	}
}

func (SysRole) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user_roles", SysUserRole.Type),
		edge.To("role_menus", SysRoleMenu.Type),
	}
}
