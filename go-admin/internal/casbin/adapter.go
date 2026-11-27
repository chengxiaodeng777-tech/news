package casbin

import (
	"context"

	"base/ent"

	"github.com/casbin/casbin/v2/model"
)

type Adapter struct {
	client *ent.Client
}

func NewAdapter(client *ent.Client) *Adapter {
	return &Adapter{client: client}
}

func (a *Adapter) LoadPolicy(m model.Model) error {
	return loadPolicy(context.Background(), a.client, m)
}

func (a *Adapter) SavePolicy(m model.Model) error {
	return savePolicy(context.Background(), a.client, m)
}

func (a *Adapter) AddPolicy(string, string, []string) error {
	return nil
}

func (a *Adapter) RemovePolicy(string, string, []string) error {
	return nil
}

func (a *Adapter) RemoveFilteredPolicy(string, string, int, ...string) error {
	return nil
}

func (a *Adapter) IsFiltered() bool {
	return false
}
