package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type SysMenu struct {
	ent.Schema
}

func (SysMenu) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_menu"},
	}
}

func (SysMenu) Mixin() []ent.Mixin {
	return []ent.Mixin{
		AuditMixin{},
	}
}

func (SysMenu) Fields() []ent.Field {
	return []ent.Field{
		field.Int("parent_id").
			Optional().
			Nillable(),
		field.String("name").
			NotEmpty(),
		field.String("title").
			NotEmpty(),
		field.String("path").
			Default(""),
		field.String("component").
			Default(""),
		field.String("icon").
			Default(""),
		field.String("type").
			NotEmpty(),
		field.String("permission_code").
			Default(""),
		field.String("redirect").
			Default(""),
		field.Int("sort").
			Default(0),
		field.Int8("status").
			Default(1),
		field.Bool("hidden").
			Default(false),
		field.Bool("keep_alive").
			Default(false),
		field.String("external_link").
			Default(""),
		field.String("remark").
			Default(""),
	}
}

func (SysMenu) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("children", SysMenu.Type).
			From("parent").
			Field("parent_id").
			Unique(),
		edge.To("role_menus", SysRoleMenu.Type),
	}
}

func (SysMenu) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("parent_id"),
	}
}
