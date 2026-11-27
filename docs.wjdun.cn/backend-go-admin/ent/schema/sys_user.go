package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type SysUser struct {
	ent.Schema
}

func (SysUser) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_user"},
	}
}

func (SysUser) Mixin() []ent.Mixin {
	return []ent.Mixin{
		AuditMixin{},
	}
}

func (SysUser) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").
			Unique().
			NotEmpty(),
		field.String("password_hash").
			NotEmpty(),
		field.String("real_name").
			Default(""),
		field.String("nickname").
			Default(""),
		field.String("mobile").
			Default(""),
		field.String("email").
			Default(""),
		field.String("avatar").
			Default(""),
		field.Int8("status").
			Default(1),
		field.Bool("is_super_admin").
			Default(false),
		field.Int("organization_id").
			Optional().
			Nillable(),
		field.String("remark").
			Default(""),
	}
}

func (SysUser) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user_roles", SysUserRole.Type),
		edge.From("organization", SysOrganization.Type).
			Ref("users").
			Field("organization_id").
			Unique(),
	}
}

func (SysUser) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("organization_id"),
	}
}
