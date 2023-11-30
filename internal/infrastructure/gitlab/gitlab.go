package gitlab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Your function that calls CreateBranchCommit
func PerformGitOperations() error {

	// Create new branch in GitLab and commit the new branch
	err := CreateBranchCommit()
	if err != nil {
		fmt.Println("Failed to create branch and commit:", err)
		return err
	}

	return nil
}

func CreateBranchCommit() error {
	// export GITLAB_ACCESS_TOKEN=""
	accessToken := os.Getenv("GITLAB_ACCESS_TOKEN")
	if accessToken == "" {
		// Handle case where environment variable is not set
		return fmt.Errorf("GITLAB_ACCESS_TOKEN environment variable is not set")
	}
	gitlabURL := "https://gitlab.com"
	projectID := "50849403"
	branchName := "new-branch"
	sourceBranch := "main" // Branch you're branching off from

	baseURL := fmt.Sprintf("%s/api/v4/projects/%s/repository/branches", gitlabURL, projectID)

	branchData := map[string]string{
		"access_token": accessToken,
	}

	jsonData, err := json.Marshal(branchData)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return err
	}

	url := fmt.Sprintf("%s?branch=%s&ref=%s", baseURL, branchName, sourceBranch)
	fmt.Println(url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Println("Failed to create branch:", resp.Status)
		return err
	}

	fmt.Println("Branch created successfully")

	return nil
}
