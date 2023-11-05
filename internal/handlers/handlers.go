package handlers

import (
	"net/http"

	"github.com/ntwklab/firewall_portal/internal/config"
	"github.com/ntwklab/firewall_portal/internal/forms"
	"github.com/ntwklab/firewall_portal/internal/helpers"
	"github.com/ntwklab/firewall_portal/internal/models"
	"github.com/ntwklab/firewall_portal/internal/render"
)

// Repo is the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "home.page.tmpl.html", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// send the data to the template
	render.RenderTemplate(w, r, "about.page.tmpl.html", &models.TemplateData{})
}

// Renders the Create Rule page
func (m *Repository) CreateRule(w http.ResponseWriter, r *http.Request) {
	var emptyCreaterule models.CreateRule
	data := make(map[string]interface{})
	data["createrule"] = emptyCreaterule

	render.RenderTemplate(w, r, "create-rule.page.tmpl.html", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostCreateRule  handles the posting of a CreateRule form
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

		render.RenderTemplate(w, r, "create-rule.page.tmpl.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	m.App.Session.Put(r.Context(), "createrule", createrule)
	http.Redirect(w, r, "/create-rule-summary", http.StatusSeeOther)
}

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

	render.RenderTemplate(w, r, "create-rule-summary.page.tmpl.html", &models.TemplateData{
		Data: data,
	})
}
