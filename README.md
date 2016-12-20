terraform-provider-arubacloud
==========================

Terraform provider for Arubacloud

## Description

With this custom terraform provider plugin you can manage your Arubacloud resources.

## Usage

Add plugin binary to your ~/.terraformrc file
```
providers {
    arubacloud = "/path/to/your/bin/terraform-provider-arubacloud"
}
```

### Provider Configuration

```
provider "arubacloud" {
    username  = "${var.arubacloud_username}"
    password  = "${var.arubacloud_password}"
    dc_number = "${var.arubacloud_dc_number}"
}
```

##### Argument Reference

The following arguments are required.

* `username` - username for accessing Arubacloud Control Panel (like AWI-12345).
* `password` - password for accessing Arubacloud Control Panel.
* `dc_number_server` - Number of the datacenter (1-6)
  * 1 -> DC1 -> Italy
  * 2 -> DC2 -> Italy
  * 3 -> DC3 -> Czech Republic
  * 4 -> DC4 -> France
  * 5 -> DC5 -> Germany
  * 6 -> DC6 -> UK

### Resource Configuration

#### `arubacloud_server_smart`

```
resource "arubacloud_server_smart" "smart-server-example" {
  smart_size       = "MEDIUM"
  name             = "smart-1"
  admin_password   = "${var.admin_password}"
  os_template_name = "Ubuntu Server 16.04 LTS 64bit"
  note             = "created with arubacloud terraform provider"
}

output "smart-server-example public ip" {
  value = "${arubacloud_server_smart.smart-server-example.public_ip}"
}
```
##### Argument Reference

The following arguments are supported.

* `smart_size` (Required) Smart server size (SMALL, MEDIUM, LARGE, EXTRALARGE)
* `name` (Required) Name of the virtual machine
* `admin_password` (Required) Admin password for accessing a VM
* `os_template_name` (Required) Operating System template name
* `note` (Optional) Free text

`public_ip` is a computed property which could be used in _output_ section

#### `arubacloud_server_pro`

```
resource "arubacloud_server_pro" "pro-server-example-1" {
  hypervisor               = "Low Cost Hyper-V"
  name                     = "srv-pro-1"
  admin_password           = "${var.admin_password}"
  os_template_name         = "CentOS 7.x 64bit"
  note                     = "my first terraformed arubacloud pro server"
  cpu_quantity             = 4
  ram_quantity             = 8
  virtual_disks            = [20, 50, 50]
}
```

##### Argument Reference

The following arguments are supported.

* `hypervisor` (Optional) VM Hypervisor (`VMWare` _(default)_, `Hyper-V`, `Low Cost Hyper-V`)
* `name` (Required) Name of the virtual machine
* `admin_password` (Required) Admin password for accessing a VM
* `os_template_name` (Required) Operating System template name
* `note` (Optional) Free text
* `cpu_quantity` (Optional) Number of Virtual CPU(s)
* `ram_quantity` (Optional) RAM in GB
* `virtual_disks` (Optional) Array of Hard Disks in GB. Default is [10]. You can add up to 4 HDDs
* `purchased_ip_resource_id` (Optional) By default **no public IP is assigned**. You can use `arubacloud_purchased_ip` resource for assigning a public IP. See example below

#### `arubacloud_purchased_ip`

```
resource "arubacloud_purchased_ip" "example" {}

output "ip" {
  value = "${arubacloud_purchased_ip.example.ip}"
}
```
##### Argument Reference

All properties are computed

* `ip`
* `subnet_mask`
* `gateway`
* `gateway_ip_v6`
* `prefix_ip_v6`
* `subnet_prefix_ip_v6`
* `start_range_ip_v6`
* `end_range_ip_v6`
* `server_id`

#### `arubacloud_vlan`

```
resource "arubacloud_vlan" "vlan-example" {
  name = "vlan-1"
}
```
##### Argument Reference

All properties are computed

* `name` (Required) Name of the VLAN

### More examples

Purchase a new IP and assign it to a PRO server:

```
resource "arubacloud_purchased_ip" "ip-srv-pro" {}

resource "arubacloud_server_pro" "cloud-srv-pro" {
  name                     = "srv-pro-terraformed"
  admin_password           = "${var.admin_password}"
  os_template_name         = "Ubuntu Server 14.04 LTS 64bit"
  virtual_disks            = [10, 80, 100]
  purchased_ip_resource_id = "${arubacloud_purchased_ip.ip-srv-pro.id}"
}

output "ip cloud-srv-pro" {
  value = "${arubacloud_server_pro.cloud-srv-pro.public_ip}"
}
```


## Contribution

Want a new feature? Do it and send a pull request.


## Licence

[MIT License](https://raw.githubusercontent.com/andrexus/terraform-provider-goarubacloud/master/LICENSE.txt)

## Author

[andrexus](https://github.com/andrexus)
