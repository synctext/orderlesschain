# Variables
variable "key_name" {}
variable "public_key" {}
variable "security_group_name" {}
variable "network_name" {}

variable "image_name" {}
variable "availability_zone" {}

variable "node_vm_name" {}
variable "node_counter" {}
variable "node_flavor_name" {}

variable "orderer_vm_name" {}
variable "orderer_counter" {}
variable "orderer_flavor_name" {}

variable "client_vm_name" {}
variable "client_counter" {}
variable "client_flavor_name" {}

variable "cli_vm_name" {}
variable "cli_counter" {}
variable "cli_flavor_name" {}

# Resources
resource "openstack_compute_keypair_v2" "access_key" {
  name       = var.key_name
  public_key = var.public_key
}

resource "openstack_compute_secgroup_v2" "security_group" {
  name        = var.security_group_name
  description = "Terraform Security Group for the System"
  rule {
    from_port   = -1
    to_port     = -1
    ip_protocol = "icmp"
    cidr        = "0.0.0.0/0"
  }

  rule {
    from_port   = 22
    to_port     = 22
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }

  rule {
    from_port   = 9000
    to_port     = 9000
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }

  rule {
    from_port   = 9001
    to_port     = 9001
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }

  #  rule {
  #    from_port   = 6060
  #    to_port     = 6060
  #    ip_protocol = "tcp"
  #    cidr        = "0.0.0.0/0"
  #  }
}

resource "openstack_compute_instance_v2" "node" {
  name              = "${var.node_vm_name}-${count.index}"
  count             = var.node_counter
  image_name        = var.image_name
  availability_zone = var.availability_zone
  flavor_name       = var.node_flavor_name
  key_pair          = "${openstack_compute_keypair_v2.access_key.name}"
  security_groups   = [
    "${openstack_compute_secgroup_v2.security_group.name}"
  ]
  network {
    name = var.network_name
  }
  depends_on        = [
    openstack_compute_secgroup_v2.security_group,
    openstack_compute_keypair_v2.access_key
  ]
}

resource "openstack_compute_instance_v2" "orderer" {
  name              = "${var.orderer_vm_name}-${count.index}"
  count             = var.orderer_counter
  image_name        = var.image_name
  availability_zone = var.availability_zone
  flavor_name       = var.orderer_flavor_name
  key_pair          = "${openstack_compute_keypair_v2.access_key.name}"
  security_groups   = [
    "${openstack_compute_secgroup_v2.security_group.name}"
  ]
  network {
    name = var.network_name
  }
  depends_on        = [
    openstack_compute_secgroup_v2.security_group,
    openstack_compute_keypair_v2.access_key
  ]
}

resource "openstack_compute_instance_v2" "client" {
  name              = "${var.client_vm_name}-${count.index}"
  count             = var.client_counter
  image_name        = var.image_name
  availability_zone = var.availability_zone
  flavor_name       = var.client_flavor_name
  key_pair          = "${openstack_compute_keypair_v2.access_key.name}"
  security_groups   = [
    "${openstack_compute_secgroup_v2.security_group.name}"
  ]
  network {
    name = var.network_name
  }
  depends_on        = [
    openstack_compute_secgroup_v2.security_group,
    openstack_compute_keypair_v2.access_key
  ]
}

resource "openstack_compute_instance_v2" "cli" {
  name              = "${var.cli_vm_name}-${count.index}"
  count             = var.cli_counter
  image_name        = var.image_name
  availability_zone = var.availability_zone
  flavor_name       = var.cli_flavor_name
  key_pair          = "${openstack_compute_keypair_v2.access_key.name}"
  security_groups   = [
    "${openstack_compute_secgroup_v2.security_group.name}"
  ]
  network {
    name = var.network_name
  }
  depends_on        = [
    openstack_compute_secgroup_v2.security_group,
    openstack_compute_keypair_v2.access_key
  ]
}

#### The Ansible inventory file ############################################
resource "local_file" "ansible_inventory" {
  content  = templatefile("./templates/inventory.tmpl",
  {
    nodes-ips   = openstack_compute_instance_v2.node[*].access_ip_v4
    orderer-ips = openstack_compute_instance_v2.orderer[*].access_ip_v4
    clients-ips = openstack_compute_instance_v2.client[*].access_ip_v4
    cli-ips     = openstack_compute_instance_v2.cli[*].access_ip_v4
  })
  filename = "../ansible/ansible_remote/inventory_dir/inventory"
}

### The Nodes Endpoints
resource "local_file" "nodes_endpoints" {
  content  = templatefile("./templates/endpoints.tmpl",
  {
    nodes-ips   = openstack_compute_instance_v2.node[*].access_ip_v4
    orderer-ips = openstack_compute_instance_v2.orderer[*].access_ip_v4
    clients-ips = openstack_compute_instance_v2.client[*].access_ip_v4
    cli-ips     = openstack_compute_instance_v2.cli[*].access_ip_v4
  })
  filename = "../../configs/endpoints_remote.yml"
}

### The Nodes Endpoints for creating certificates
resource "local_file" "nodes_endpoints_certificate" {
  content  = templatefile("./templates/endpoints_certificate.tmpl",
  {
    nodes-ips   = openstack_compute_instance_v2.node[*].access_ip_v4
    orderer-ips = openstack_compute_instance_v2.orderer[*].access_ip_v4
    clients-ips = openstack_compute_instance_v2.client[*].access_ip_v4
  })
  filename = "../../certificates/endpoints_remote"
}
###############################################################################

### The Nodes Endpoints for opening log terminal
resource "local_file" "nodes_endpoints_terminal" {
  content  = templatefile("./templates/endpoints_terminals_nodes.tmpl",
  {
    nodes-ips = slice(openstack_compute_instance_v2.node[*].access_ip_v4, 0, 1)
  })
  filename = "../../scripts/open_terminals/remote_endpoints_nodes"
}

### The Orderer Endpoints for opening log terminal
resource "local_file" "orderer_endpoints_terminal" {
  content  = templatefile("./templates/endpoints_terminals_orderer.tmpl",
  {
    orderer-ips = slice(openstack_compute_instance_v2.orderer[*].access_ip_v4, 0, 1)
  })
  filename = "../../scripts/open_terminals/remote_endpoints_orderer"
}

### The Client Endpoints for opening log terminal
resource "local_file" "clients_endpoints_terminal" {
  content  = templatefile("./templates/endpoints_terminals_clients.tmpl",
  {
    clients-ips = slice(openstack_compute_instance_v2.client[*].access_ip_v4, 0, 1)
  })
  filename = "../../scripts/open_terminals/remote_endpoints_clients"
}

# Output
output "public_ip_nodes" {
  value = openstack_compute_instance_v2.node[*].access_ip_v4
}

output "public_ip_orderer" {
  value = openstack_compute_instance_v2.orderer[*].access_ip_v4
}

output "public_ip_clients" {
  value = openstack_compute_instance_v2.client[*].access_ip_v4
}

output "public_ip_cli" {
  value = openstack_compute_instance_v2.cli[*].access_ip_v4
}


