package ciscoasa

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
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
	var branchName string = "FirewallRules"
	var commitMessage string = "Firewall Rule Update [ci skip]"
	var repoPath string = "/Users/stefankelly/Terraform/asa/asaterraform/"

	// Check if the branch exists
	checkCmd := exec.Command("git", "show-ref", "--verify", fmt.Sprintf("refs/heads/%s", branchName))
	checkCmd.Dir = repoPath

	err := checkCmd.Run()
	if err != nil {
		// Branch doesn't exist, create a new one
		createCmd := exec.Command("git", "checkout", "-b", branchName)
		createCmd.Dir = repoPath

		out, err := createCmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error creating branch:", err)
			fmt.Println("Output:", string(out))
			return err
		}
		fmt.Println("Branch created successfully:", branchName)
	} else {
		// Branch exists, switch to it
		switchCmd := exec.Command("git", "checkout", branchName)
		switchCmd.Dir = repoPath

		out, err := switchCmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error switching branch:", err)
			fmt.Println("Output:", string(out))
			return err
		}
		fmt.Println("Switched to branch:", branchName)
	}

	// Open the file in append mode or create if it doesn't exist
	file, err := os.OpenFile("/Users/stefankelly/Terraform/asa/asaterraform/main.tf", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Append the ASA config to the file
	fmt.Println(asaConfig)
	_, err = file.WriteString(asaConfig + "\n")
	if err != nil {
		log.Printf("Error writing to file: %v\n", err)
		return err
	}
	log.Println("File written successfully")

	// Git add, commit, and push
	cmds := [][]string{
		{"git", "add", "."},
		{"git", "commit", "-am", commitMessage},
		{"git", "push", "-u", "origin", branchName},
	}

	for _, cmdArgs := range cmds {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Dir = repoPath

		err := cmd.Run()
		if err != nil {
			log.Printf("Error with %s: %v\n", strings.Join(cmdArgs, " "), err)
			return err
		}
	}

	return nil
}
