package arubacloud

import (
	"log"
	"github.com/andrexus/goarubacloud"
)

type Config struct {
	DatacenterRegion int
	Username         string
	Password         string
}

// Client() returns a new client for accessing Arubacloud.
func (c *Config) Client() (*goarubacloud.Client, error) {
	client := goarubacloud.NewClient(goarubacloud.DataCenterRegion(c.DatacenterRegion), c.Username, c.Password)

	log.Printf("[INFO] Arubacloud Client configured for URL: %s", client.BaseURL.String())

	return client, nil
}