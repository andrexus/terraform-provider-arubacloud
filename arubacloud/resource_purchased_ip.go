package arubacloud

import (
	"log"
	"strconv"

	"fmt"

	"github.com/andrexus/goarubacloud"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArubacloudPurchasedIP() *schema.Resource {
	return &schema.Resource{
		Create: resourceArubacloudPurchasedIPCreate,
		Read:   resourceArubacloudPurchasedIPRead,
		Delete: resourceArubacloudPurchasedIPDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"subnet_mask": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"gateway": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"gateway_ip_v6": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"prefix_ip_v6": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			"subnet_prefix_ip_v6": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"start_range_ip_v6": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"end_range_ip_v6": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"server_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceArubacloudPurchasedIPCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goarubacloud.Client)

	purchased_ip, _, err := client.PurchasedIPs.Purchase()

	if err != nil {
		return err
	}

	// Assign the purchased_ips id
	d.SetId(strconv.Itoa(purchased_ip.ResourceId))

	log.Printf("[INFO] Purchased IP: %s. ID: %s", purchased_ip.Value, d.Id())

	return resourceArubacloudPurchasedIPRead(d, meta)
}

func resourceArubacloudPurchasedIPRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goarubacloud.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid purchased IP resource id: %v", err)
	}

	// Retrieve the purchased IP for updating the state
	purchased_ips, _, err := client.PurchasedIPs.List()
	if err != nil {
		return fmt.Errorf("Error retrieving purchased ips: %s", err)
	}

	for _, purchased_ip := range purchased_ips {
		if purchased_ip.ResourceId == id {
			d.Set("ip", purchased_ip.Value)
			d.Set("subnet_mask", purchased_ip.SubNetMask)
			d.Set("gateway", purchased_ip.Gateway)
			d.Set("gateway_ip_v6", purchased_ip.GatewayIPv6)
			d.Set("prefix_ip_v6", purchased_ip.PrefixIPv6)
			d.Set("subnet_prefix_ip_v6", purchased_ip.SubnetPrefixIPv6)
			d.Set("start_range_ip_v6", purchased_ip.StartRangeIPv6)
			d.Set("end_range_ip_v6", purchased_ip.EndRangeIPv6)
			d.Set("server_id", purchased_ip.ServerId)

			return nil
		}
	}

	return nil
}

func resourceArubacloudPurchasedIPDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goarubacloud.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid purchased IP resource id: %v", err)
	}

	_, err = client.PurchasedIPs.Delete(id)

	if err != nil {
		return err
	}

	return nil
}
