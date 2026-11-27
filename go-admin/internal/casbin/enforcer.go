package casbin

import (
	"context"
	_ "embed"

	"base/ent"

	casbinlib "github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	casbinpersist "github.com/casbin/casbin/v2/persist"
)

//go:embed model.conf
var modelText string

func NewEnforcer(client *ent.Client) (*casbinlib.Enforcer, error) {
	m, err := model.NewModelFromString(modelText)
	if err != nil {
		return nil, err
	}

	adapter := NewAdapter(client)
	enforcer, err := casbinlib.NewEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}
	return enforcer, nil
}

func loadPolicy(ctx context.Context, client *ent.Client, m model.Model) error {
	rules, err := client.SysCasbinRule.Query().All(ctx)
	if err != nil {
		return err
	}
	for _, rule := range rules {
		line := []string{rule.V0, rule.V1, rule.V2, rule.V3, rule.V4, rule.V5}
		if err := casbinpersist.LoadPolicyLine(rule.Ptype+", "+joinPolicy(line), m); err != nil {
			return err
		}
	}
	return nil
}

func savePolicy(ctx context.Context, client *ent.Client, m model.Model) error {
	if _, err := client.SysCasbinRule.Delete().Exec(ctx); err != nil {
		return err
	}

	for ptype, ast := range m["p"] {
		for _, rule := range ast.Policy {
			if err := createRule(ctx, client, ptype, rule); err != nil {
				return err
			}
		}
	}
	for ptype, ast := range m["g"] {
		for _, rule := range ast.Policy {
			if err := createRule(ctx, client, ptype, rule); err != nil {
				return err
			}
		}
	}
	return nil
}

func createRule(ctx context.Context, client *ent.Client, ptype string, rule []string) error {
	values := make([]string, 6)
	copy(values, rule)

	return client.SysCasbinRule.Create().
		SetPtype(ptype).
		SetV0(values[0]).
		SetV1(values[1]).
		SetV2(values[2]).
		SetV3(values[3]).
		SetV4(values[4]).
		SetV5(values[5]).
		Exec(ctx)
}

func joinPolicy(values []string) string {
	result := ""
	for _, value := range values {
		if value == "" {
			continue
		}
		if result != "" {
			result += ", "
		}
		result += value
	}
	return result
}

func Reload(ctx context.Context, client *ent.Client, enforcer *casbinlib.Enforcer) error {
	return enforcer.LoadPolicy()
}
