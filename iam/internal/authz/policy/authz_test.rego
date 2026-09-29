package authz_test

import data.authz

client := {"uuid": "11111111-1111-1111-1111-111111111111", "role": "client"}

manager := {"uuid": "22222222-2222-2222-2222-222222222222", "role": "manager"}

own := {"owner_uuid": client.uuid}

foreign := {"owner_uuid": "33333333-3333-3333-3333-333333333333"}

test_client_can_create if {
	authz.allow with input as {"subject": client, "action": "order:create", "resource": {}}
}

test_client_manages_own_order if {
	every action in ["order:read", "order:pay", "order:cancel"] {
		authz.allow with input as {"subject": client, "action": action, "resource": own}
	}
}

test_client_denied_foreign_order if {
	every action in ["order:read", "order:pay", "order:cancel"] {
		not authz.allow with input as {"subject": client, "action": action, "resource": foreign}
	}
}

test_manager_reads_and_cancels_any_order if {
	every action in ["order:read", "order:cancel"] {
		authz.allow with input as {"subject": manager, "action": action, "resource": foreign}
	}
}

test_manager_cannot_create_or_pay if {
	not authz.allow with input as {"subject": manager, "action": "order:create", "resource": {}}
	not authz.allow with input as {"subject": manager, "action": "order:pay", "resource": foreign}
}

test_unknown_role_denied if {
	not authz.allow with input as {"subject": {"uuid": "x", "role": ""}, "action": "order:create", "resource": {}}
}
