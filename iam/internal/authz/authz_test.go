package authz_test

import (
	"context"
	"testing"

	"github.com/open-policy-agent/opa/v1/tester"
	"github.com/stretchr/testify/require"

	"github.com/PabloGolobaro/cosmic_factory/iam/internal/authz"
	"github.com/PabloGolobaro/cosmic_factory/iam/internal/model"
)

const (
	clientUUID  = "11111111-1111-1111-1111-111111111111"
	managerUUID = "22222222-2222-2222-2222-222222222222"
	foreignUUID = "33333333-3333-3333-3333-333333333333"
)

func TestEngine_Allow(t *testing.T) {
	engine, err := authz.New(context.Background())
	require.NoError(t, err)

	client := model.AuthzSubject{UUID: clientUUID, Role: model.RoleClient}
	manager := model.AuthzSubject{UUID: managerUUID, Role: model.RoleManager}
	own := model.AuthzResource{OwnerUUID: clientUUID}
	foreign := model.AuthzResource{OwnerUUID: foreignUUID}
	none := model.AuthzResource{}

	tests := []struct {
		name     string
		subject  model.AuthzSubject
		action   string
		resource model.AuthzResource
		want     bool
	}{
		{"client creates order", client, "order:create", none, true},
		{"client reads own order", client, "order:read", own, true},
		{"client pays own order", client, "order:pay", own, true},
		{"client cancels own order", client, "order:cancel", own, true},
		{"client reads foreign order", client, "order:read", foreign, false},
		{"client pays foreign order", client, "order:pay", foreign, false},
		{"client cancels foreign order", client, "order:cancel", foreign, false},
		{"client reads order without owner", client, "order:read", none, false},
		{"manager reads any order", manager, "order:read", foreign, true},
		{"manager cancels any order", manager, "order:cancel", foreign, true},
		{"manager cannot create order", manager, "order:create", none, false},
		{"manager cannot pay order", manager, "order:pay", foreign, false},
		{"unknown role denied", model.AuthzSubject{UUID: clientUUID}, "order:create", none, false},
		{"unknown action denied", client, "part:delete", own, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := engine.Allow(context.Background(), model.AuthzInput{
				Subject:  tt.subject,
				Action:   tt.action,
				Resource: tt.resource,
			})
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

// TestRegoPolicies runs policy/*_test.rego so rego tests execute without the opa CLI.
func TestRegoPolicies(t *testing.T) {
	results, err := tester.Run(context.Background(), "policy")
	require.NoError(t, err)
	require.NotEmpty(t, results)

	for _, r := range results {
		require.Truef(t, r.Pass(), "%s.%s: %v", r.Package, r.Name, r.Error)
	}
}
