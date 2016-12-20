package arubacloud

import (
	"fmt"

	"strings"

	"errors"

	"github.com/andrexus/goarubacloud"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArubacloudServerPro() *schema.Resource {
	return &schema.Resource{
		Create: resourceArubacloudServerProCreate,
		Read:   resourceArubacloudServerProRead,
		Update: resourceArubacloudServerProUpdate,
		Delete: resourceArubacloudServerProDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"hypervisor": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "VMWare",
				ForceNew:     true,
				ValidateFunc: validateHypervisorName,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"admin_password": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},

			"os_template_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"purchased_ip_resource_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"note": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"cpu_quantity": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"ram_quantity": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"virtual_disks": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},

			"dc_number": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"public_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"network_adapters": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mac_address": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceArubacloudServerProCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goarubacloud.Client)

	hypervisorName := d.Get("hypervisor").(string)
	hypervisorType := goarubacloud.VMWare_Cloud_Pro

	switch strings.ToUpper(hypervisorName) {
	case "VMWARE":
		hypervisorType = goarubacloud.VMWare_Cloud_Pro
	case "HYPER-V":
		hypervisorType = goarubacloud.Microsoft_Hyper_V
	case "LOW COST HYPER-V":
		hypervisorType = goarubacloud.Microsoft_Hyper_V_Low_Cost
	}

	os_template_name := d.Get("os_template_name").(string)
	template, err := client.Hypervisors.FindOsTemplate(hypervisorType, os_template_name)

	if err != nil {
		return err
	}

	// Build up our creation options
	createRequest := goarubacloud.NewCloudServerProCreateRequest(d.Get("name").(string), d.Get("admin_password").(string), template.Id)

	if attr, ok := d.GetOk("cpu_quantity"); ok {
		createRequest.SetCPUQuantity(attr.(int))
	}

	if attr, ok := d.GetOk("ram_quantity"); ok {
		createRequest.SetRAMQuantity(attr.(int))
	}

	if attr, ok := d.GetOk("note"); ok {
		createRequest.SetNote(attr.(string))
	}

	if attr, ok := d.GetOk("virtual_disks"); ok {
		for _, size := range attr.([]interface{}) {
			err = createRequest.AddVirtualDisk(size.(int))
			if err != nil {
				return err
			}
		}
	}

	if attr, ok := d.GetOk("purchased_ip_resource_id"); ok {
		// Add public ip
		err = createRequest.AddPublicIp(attr.(int))
		if err != nil {
			return fmt.Errorf("Unable to add a public IP address for the server %d. Error: %s", d.Id(), err)
		}
	}

	err = resourceArubacloudServerCreate(createRequest, d, meta)
	if err != nil {
		return err
	}

	return resourceArubacloudServerProRead(d, meta)
}

func resourceArubacloudServerProUpdate(d *schema.ResourceData, meta interface{}) error {
	return errors.New("Not implemented yet")
}

func resourceArubacloudServerProRead(d *schema.ResourceData, meta interface{}) error {
	return resourceArubacloudServerRead(d, meta)
}

func resourceArubacloudServerProDelete(d *schema.ResourceData, meta interface{}) error {
	return resourceArubacloudServerDelete(d, meta)
}
