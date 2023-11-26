package ciscoasa

import (
	"fmt"
	"os"
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

func AppendASAConfigToFile(asaConfig string) error {
	// Open the file in append mode or create if it doesn't exist
	file, err := os.OpenFile("cisco_asa_terraform.tf", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Append the ASA config to the file
	_, err = file.WriteString(asaConfig + "\n")
	if err != nil {
		return err
	}

	return nil
}
