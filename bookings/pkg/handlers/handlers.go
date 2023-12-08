package handlers

import (
	"bookings/pkg/config"
	"bookings/pkg/models"
	"bookings/pkg/render"
	"encoding/json"
	"fmt"
	"net/http"
)

type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
}

// Repository type
type Respository struct {
	app *config.AppConfig
}

// Repo the repository used by the handlers
var Repo *Respository

// NewRepository creates a new repository
func NewRepository(a *config.AppConfig) *Respository {
	return &Respository{
		app: a,
	}
}

// NewHandler sets the repository for the handlers
func NewHandler(r *Respository) {
	Repo = r
}

// Home is the home page handler
func (m *Respository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr

	fmt.Println("Remote IP", remoteIP)

	m.app.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (m *Respository) About(w http.ResponseWriter, r *http.Request) {

	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"

	remoteIP := m.app.Session.GetString(r.Context(), "remote_ip")

	fmt.Println("Remote IP 2", remoteIP)

	stringMap["remote_ip"] = remoteIP

	fmt.Println("StringMap", stringMap)

	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Generals is the generals page handler
func (m *Respository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors is the majors page handler
func (m *Respository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Availability is the availability page handler
func (m *Respository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability is the post availability page handler
func (m *Respository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	fmt.Println("REQUEST BODY ", r.Form)
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))
}

type jsonResponse struct {
	OK           bool   `json:"ok"`
	Availability string `json:"availability"`
}

// AvailabilityJSON handles request for availability and sends json response back
func (m *Respository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	jsonResponse := jsonResponse{
		OK:           true,
		Availability: "Available",
	}

	js, _ := json.Marshal(&jsonResponse)

	w.Write(js)
}

// Reservation is the reservation page handler
func (m *Respository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{})
}

// Contact is the contact page handler
func (m *Respository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}
