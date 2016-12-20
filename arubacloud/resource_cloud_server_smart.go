package arubacloud

import (
	"github.com/andrexus/goarubacloud"
	"github.com/hashicorp/terraform/helper/schema"
	"errors"
)

func resourceArubacloudServerSmart() *schema.Resource {
	return &schema.Resource{
		Create: resourceArubacloudServerSmartCreate,
		Read:   resourceArubacloudServerSmartRead,
		Update: resourceArubacloudServerSmartUpdate,
		Delete: resourceArubacloudServerSmartDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"smart_size": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSmartServerSize,
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

			"note": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			"dc_number": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"public_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArubacloudServerSmartCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goarubacloud.Client)

	hypervisorType := goarubacloud.VMWare_Cloud_Smart

	os_template_name := d.Get("os_template_name").(string)
	template, err := client.Hypervisors.FindOsTemplate(hypervisorType, os_template_name)

	if err != nil {
		return err
	}

	smartServerSize, _ := goarubacloud.GetServerSmartSize(d.Get("smart_size").(string))

	// Build up our creation options
	createRequest := goarubacloud.NewCloudServerSmartCreateRequest(smartServerSize, d.Get("name").(string), d.Get("admin_password").(string), template.Id)

	if attr, ok := d.GetOk("note"); ok {
		createRequest.SetNote(attr.(string))
	}

	err = resourceArubacloudServerCreate(createRequest, d, meta)
	if err != nil {
		return err
	}

	return resourceArubacloudServerSmartRead(d, meta)
}

func resourceArubacloudServerSmartUpdate(d *schema.ResourceData, meta interface{}) error {
	return errors.New("Not implemented yet")
}

func resourceArubacloudServerSmartRead(d *schema.ResourceData, meta interface{}) error {
	return resourceArubacloudServerRead(d, meta)
}

func resourceArubacloudServerSmartDelete(d *schema.ResourceData, meta interface{}) error {
	return resourceArubacloudServerDelete(d, meta)
}
