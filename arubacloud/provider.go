package arubacloud

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

const (
	provider_name           = "terraform-provider-arubacloud"
	provider_version string = "v0.1.0"
)

// Provider returns a schema.Provider for Arubacloud.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"dc_number": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARUBACLOUD_REGION", nil),
				Description: "The DC number for API operations.",
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARUBACLOUD_USERNAME", nil),
				Description: "Username for API operations.",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARUBACLOUD_PASSWORD", nil),
				Description: "Password for API operations.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"arubacloud_server_pro":   resourceArubacloudServerPro(),
			"arubacloud_server_smart": resourceArubacloudServerSmart(),
			"arubacloud_purchased_ip": resourceArubacloudPurchasedIP(),
			"arubacloud_vlan":         resourceArubacloudVLAN(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	log.Printf("[DEBUG] Configure %s. Version %s", provider_name, provider_version)
	config := Config{
		DatacenterRegion: d.Get("dc_number").(int),
		Username:         d.Get("username").(string),
		Password:         d.Get("password").(string),
	}

	return config.Client()
}
