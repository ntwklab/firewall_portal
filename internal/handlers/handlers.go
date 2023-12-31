package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ntwklab/firewall_portal/internal/config"
	"github.com/ntwklab/firewall_portal/internal/driver"
	"github.com/ntwklab/firewall_portal/internal/forms"
	"github.com/ntwklab/firewall_portal/internal/helpers"
	ciscoasa "github.com/ntwklab/firewall_portal/internal/infrastructure/cisco_asa"
	"github.com/ntwklab/firewall_portal/internal/models"
	"github.com/ntwklab/firewall_portal/internal/render"
	"github.com/ntwklab/firewall_portal/internal/repository"
	"github.com/ntwklab/firewall_portal/internal/repository/dbrepo"
)

// Repo is the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl.html", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// send the data to the template
	render.Template(w, r, "about.page.tmpl.html", &models.TemplateData{})
}

// Renders the Create Rule page
func (m *Repository) CreateRule(w http.ResponseWriter, r *http.Request) {
	var emptyCreaterule models.CreateRule
	data := make(map[string]interface{})
	data["createrule"] = emptyCreaterule

	render.Template(w, r, "create-rule.page.tmpl.html", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// Displays the form inputs to the user once the new rule ahs been added to the DB
func (m *Repository) CreateRuleSummary(w http.ResponseWriter, r *http.Request) {
	createrule, ok := m.App.Session.Get(r.Context(), "createrule").(models.CreateRule)
	if !ok {
		m.App.ErrorLog.Println("Can't get error from session")
		m.App.Session.Put(r.Context(), "error", "Can't get createrule from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "createrule")
	data := make(map[string]interface{})
	data["createrule"] = createrule

	render.Template(w, r, "create-rule-summary.page.tmpl.html", &models.TemplateData{
		Data: data,
	})
}

// PostCreateRule handles the posting of a CreateRule form to the DB and write the Terraform
func (m *Repository) PostCreateRule(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	createrule := models.CreateRule{
		SourceIP:      r.Form.Get("source_ip"),
		DestinationIP: r.Form.Get("destination_ip"),
		Port:          r.Form.Get("port"),
	}

	form := forms.New(r.PostForm)

	form.Required("source_ip", "destination_ip", "port")
	form.MinLength("source_ip", 3, r)

	if !form.Valid() {
		data := make(map[string]interface{})
		data["createrule"] = createrule

		render.Template(w, r, "create-rule.page.tmpl.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	err = m.DB.InsertRule(createrule)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "createrule", createrule)
	http.Redirect(w, r, "/create-rule-summary", http.StatusSeeOther)

	// Terraform
	// Replace dots with underscores
	formattedSourceIP := strings.ReplaceAll(createrule.SourceIP, ".", "_")
	formattedDestinationIP := strings.ReplaceAll(createrule.DestinationIP, ".", "_")

	ruleName := fmt.Sprintf("_%s_%s_%s", formattedSourceIP, formattedDestinationIP, createrule.Port)
	intf := "OUTSIDE"
	source := createrule.SourceIP
	destination := createrule.DestinationIP
	service := fmt.Sprintf("tcp/%s", createrule.Port)

	asaConfig := ciscoasa.GenerateASAConfig(ruleName, intf, source, destination, service)

	//Pull repo to go dir

	// Write to a file to main.tf
	ciscoasa.AppendASAConfigToFile(asaConfig)

	// Create new branch in GitLab and commit the New branch
	// err = gitlab.CreateBranchCommit()
	// if err != nil {
	// 	fmt.Println("Error performing Git operations:", err)
	// 	return
	// }

}

// Check for duplicates in the database without refresh of page, and display error banner if duplicate found
func (app *Repository) CheckDuplicate(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		// Handle parsing error
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve form values
	sourceIP := r.Form.Get("source_ip")
	destinationIP := r.Form.Get("destination_ip")
	port := r.Form.Get("port")

	// Create a CreateRule object with the form data
	createrule := models.CreateRule{
		SourceIP:      sourceIP,
		DestinationIP: destinationIP,
		Port:          port,
	}

	// Check for duplicate rule
	duplicateExists, err := app.DB.CheckDuplicateRule(createrule)
	if err != nil {
		// Handle database error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare response JSON
	response := map[string]bool{
		"duplicateExists": duplicateExists,
	}

	// Set response headers and write JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
