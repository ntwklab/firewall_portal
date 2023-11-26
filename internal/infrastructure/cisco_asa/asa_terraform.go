package ciscoasa

import (
	"fmt"
)

// GenerateASAConfig generates Terraform configuration for Cisco ASA firewall rule
func GenerateASAConfig(ruleName, intf, source, destination, service string) string {
	config := fmt.Sprintf(`
resource "ciscoasa_access_in_rules" "%s" {
  interface = "%s"
  rule {
    source              = "%s"
    destination         = "%s"
    destination_service = "%s"
  }
}
`, ruleName, intf, source, destination, service)

	return config
}
