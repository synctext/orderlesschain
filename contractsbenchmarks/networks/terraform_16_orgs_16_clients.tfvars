# ATTENTION: DO NOT Include Sensitive Data Here!
# Cluster configurations
key_name = "terraform_key"
security_group_name = "terraform_default_group"
network_name = "cluster-net" # From project on OpenStack under Network -> Networks

# Network configurations
image_name = "kvm-ubuntu-focal"
availability_zone = "kvm-hdd"
node_vm_name = "node"
node_flavor_name = "m1.large"
node_counter = 16
orderer_vm_name = "orderer"
orderer_flavor_name = "m1.large"
orderer_counter = 1
client_vm_name = "client"
client_counter = 16
client_flavor_name = "m1.large"
cli_vm_name = "cli"
cli_counter = 1
cli_flavor_name = "m1.medium"

