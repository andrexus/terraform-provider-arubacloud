package arubacloud

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/andrexus/goarubacloud"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArubacloudServerCreate(createRequest goarubacloud.CloudServerCreator, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goarubacloud.Client)

	log.Printf("[DEBUG] Cloud server create configuration: %#v", createRequest)

	cloud_server, _, err := client.CloudServers.Create(createRequest)

	if err != nil {
		return fmt.Errorf("Error creating cloud server: %s", err)
	}

	// Assign the cloud_servers id
	d.SetId(strconv.Itoa(cloud_server.ServerId))

	log.Printf("[INFO] Cloud server ID: %s", d.Id())

	err = goarubacloud.WaitForServerCreationDone(client, cloud_server.ServerId)
	if err != nil {
		return fmt.Errorf(
			"Error waiting for server (%s) creation is done. Error: %s", d.Id(), err)
	}

	return resourceArubacloudServerProRead(d, meta)
}

func resourceArubacloudServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goarubacloud.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid cloud server id: %v", err)
	}

	// Retrieve the cloud server details for updating the state
	cloud_server, _, err := client.CloudServers.Get(id)
	if err != nil {
		error_response := err.(*goarubacloud.ErrorResponse)
		// check if the cloud server no longer exists.
		if error_response.ResultCode == 15 { // resultCode 15 means not found
			log.Printf("[WARN] Arubacloud server (%s) not found", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving cloud server: %s", err)
	}

	d.Set("name", cloud_server.Name)
	d.Set("cpu_quantity", cloud_server.CPUQuantity)
	d.Set("ram_quantity", cloud_server.RAMQuantity)
	d.Set("dc_number", cloud_server.DatacenterId)
	d.Set("note", cloud_server.Note)
	d.Set("hypervisor", cloud_server.HypervisorType)

	setNetworkAdaptersData(d, cloud_server.NetworkAdapters)

	public_ip, err := cloud_server.GetPublicIpAddress()
	if err == nil {
		d.Set("public_ip", public_ip)
	}

	return nil
}

func resourceArubacloudServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*goarubacloud.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid cloud server id: %v", err)
	}

	client.CloudServerActions.PowerOff(id)

	err = goarubacloud.WaitForServerStatus(client, id, goarubacloud.OFF)

	if err != nil {
		return fmt.Errorf(
			"Error waiting for cloud server to be shut down for destroy (%s): %s", d.Id(), err)
	}

	log.Printf("[INFO] Deleting cloud server: %s", d.Id())

	// Destroy the cloud server
	_, err = client.CloudServers.Delete(id)

	// Handle remotely destroyed servers
	error_response := err.(*goarubacloud.ErrorResponse)
	// check if the cloud server no longer exists.
	if error_response.ResultCode == 15 { // resultCode 15 means not found
		return nil
	}

	if err != nil {
		return fmt.Errorf("Error deleting cloud server: %s", err)
	}

	// Wait for resources to be released (like assigned purchased IPs)
	time.Sleep(30 * time.Second)

	return nil
}

func setNetworkAdaptersData(d *schema.ResourceData, network_adapters []goarubacloud.NetworkAdapter) error {
	networkAdaptersData := make([]map[string]interface{}, 0, len(network_adapters))

	for _, network_adapter := range network_adapters {
		if network_adapter.Id != 0 {
			networkAdaptersData = append(networkAdaptersData, map[string]interface{}{
				"id":          network_adapter.Id,
				"mac_address": network_adapter.MacAddress,
			})
		}
	}

	return d.Set("network_adapters", networkAdaptersData)
}
