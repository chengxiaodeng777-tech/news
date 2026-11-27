package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type SysRoleMenu struct {
	ent.Schema
}

func (SysRoleMenu) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_role_menu"},
	}
}

func (SysRoleMenu) Mixin() []ent.Mixin {
	return []ent.Mixin{
		AuditMixin{},
	}
}

func (SysRoleMenu) Fields() []ent.Field {
	return []ent.Field{
		field.Int("role_id"),
		field.Int("menu_id"),
	}
}

func (SysRoleMenu) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("role", SysRole.Type).
			Ref("role_menus").
			Field("role_id").
			Unique().
			Required(),
		edge.From("menu", SysMenu.Type).
			Ref("role_menus").
			Field("menu_id").
			Unique().
			Required(),
	}
}

func (SysRoleMenu) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("role_id", "menu_id").
			Unique(),
	}
}
