package arubacloud

import (
	"fmt"

	"github.com/andrexus/goarubacloud"
	"strings"
)

func validateHypervisorName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	switch strings.ToUpper(value) {
	case "VMWARE":
	case "HYPER-V":
	case "LOW COST HYPER-V":
	default:
		errors = append(errors, fmt.Errorf("hypervisor '%s' is wrong. Supported values are: 'VMWare' (default), 'Hyper-V', 'Low Cost Hyper-V'", value))
	}

	return
}

func validateSmartServerSize(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	_, err := goarubacloud.GetServerSmartSize(value)

	if err != nil {
		errors = append(errors, fmt.Errorf("Validation failed for field '%s'. Error: %s", k, err))
	}

	return
}
