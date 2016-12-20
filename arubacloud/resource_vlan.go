package arubacloud

import (
	"log"
	"strconv"

	"fmt"

	"github.com/andrexus/goarubacloud"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArubacloudVLAN() *schema.Resource {
	return &schema.Resource{
		Create: resourceArubacloudVLANCreate,
		Read:   resourceArubacloudVLANRead,
		Delete: resourceArubacloudVLANDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vlan_code": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"server_ids": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func resourceArubacloudVLANCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goarubacloud.Client)

	vlan, _, err := client.VLANs.Purchase(d.Get("name").(string))

	if err != nil {
		return err
	}

	// Assign the VLANs id
	d.SetId(strconv.Itoa(vlan.ResourceId))

	log.Printf("[INFO] VLAN ID: %s", d.Id())

	return resourceArubacloudVLANRead(d, meta)
}

func resourceArubacloudVLANRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goarubacloud.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid VLAN resource id: %v", err)
	}

	// Retrieve the VLAN for updating the state
	all_vlans, _, err := client.VLANs.List()
	if err != nil {
		return fmt.Errorf("Error retrieving VLANs: %s", err)
	}

	for _, vlan := range all_vlans {
		if vlan.ResourceId == id {
			d.Set("name", vlan.Name)
			d.Set("vlan_code", vlan.VlanCode)
			d.Set("server_ids", vlan.ServerIds)
			return nil
		}
	}

	return nil
}

func resourceArubacloudVLANDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goarubacloud.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid VLAN resource id: %v", err)
	}

	_, err = client.VLANs.Delete(id)

	if err != nil {
		return err
	}

	return nil
}
