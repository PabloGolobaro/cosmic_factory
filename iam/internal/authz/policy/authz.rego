# Access policy for cosmic_factory.
#
# input:
#   subject:  {uuid, role}      — who performs the action
#   action:   "<resource>:<op>" — e.g. "order:read"
#   resource: {owner_uuid}      — empty owner_uuid for actions without a resource
package authz

default allow := false

client_owned_actions := {"order:read", "order:pay", "order:cancel"}

manager_actions := {"order:read", "order:cancel"}

is_owner if {
	input.resource.owner_uuid != ""
	input.resource.owner_uuid == input.subject.uuid
}

allow if {
	input.subject.role == "client"
	input.action == "order:create"
}

allow if {
	input.subject.role == "client"
	client_owned_actions[input.action]
	is_owner
}

allow if {
	input.subject.role == "manager"
	manager_actions[input.action]
}
