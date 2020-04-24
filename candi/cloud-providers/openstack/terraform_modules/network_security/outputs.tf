output "security_group_names" {
  value = var.enabled ? [
    data.openstack_networking_secgroup_v2.default[0].name,
    openstack_networking_secgroup_v2.kube[0].name
  ] : []
}
