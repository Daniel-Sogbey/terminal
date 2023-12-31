package handlers

import (
	"bookings/internals/config"
	"bookings/internals/driver"
	"bookings/internals/forms"
	"bookings/internals/helpers"
	"bookings/internals/models"
	"bookings/internals/render"
	"bookings/internals/repository"
	"bookings/internals/repository/dbrepo"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// Repository type
type Respository struct {
	app *config.AppConfig
	DB  repository.DatabaseRepo
}

// Repo the repository used by the handlers
var Repo *Respository

// NewRepository creates a new repository
func NewRepository(a *config.AppConfig, db *driver.DB) *Respository {
	return &Respository{
		app: a,
		DB:  dbrepo.NewPostgresRepo(a, db.SQL),
	}
}

// NewHandler sets the repository for the handlers
func NewHandler(r *Respository) {
	Repo = r
}

// Home is the home page handler
func (m *Respository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (m *Respository) About(w http.ResponseWriter, r *http.Request) {

	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"

	remoteIP := m.app.Session.GetString(r.Context(), "remote_ip")

	fmt.Println("Remote IP 2", remoteIP)

	stringMap["remote_ip"] = remoteIP

	fmt.Println("StringMap", stringMap)

	render.Template(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Generals is the generals page handler
func (m *Respository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors is the majors page handler
func (m *Respository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Availability is the availability page handler
func (m *Respository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability is the post availability page handler
func (m *Respository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	fmt.Println("REQUEST BODY ", r.Form)
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handles request for availability and sends json response back
func (m *Respository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	jsonResponse := jsonResponse{
		OK:      true,
		Message: "Available",
	}

	js, err := json.Marshal(&jsonResponse)

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(js)
}

// Reservation is the reservation page handler
func (m *Respository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})

	data["reservation"] = emptyReservation

	form := forms.New(nil)
	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: form,
		Data: data,
	})
}

// PostReservation handles the posting of a reservation form
func (m *Respository) PostReservation(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// date layout 01-02-2006

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")
	roomId, err := strconv.Atoi(r.Form.Get("room_id"))

	if err != nil {
		helpers.ServerError(w, err)
	}

	layout := "01-02-2006"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
	}

	endDate, err := time.Parse(layout, ed)

	if err != nil {
		helpers.ServerError(w, err)
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomId,
	}

	// m.app.InfoLog.Println("POST FORM ", r.PostForm)
	// m.app.InfoLog.Println("FORM ", r.Form)

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 2)
	form.MinLength("last_name", 2)
	form.IsPhone("phone", 10)
	form.IsEmail("email")

	data := make(map[string]interface{})

	if !form.Valid() {

		data["reservation"] = reservation

		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})

		return
	}

	fmt.Println(reservation)

	data["reservation"] = reservation

	err = m.DB.InsertReservation(&reservation)

	if err != nil {
		helpers.ServerError(w, err)
	}

	m.app.Session.Put(r.Context(), "reservation", reservation)
	m.app.Session.Put(r.Context(), "flash", "Reservation submitted successfully")

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

// ReservationSummary shows the summary of the reservation
func (m *Respository) ReservationSummary(w http.ResponseWriter, r *http.Request) {

	reservation, ok := m.app.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		m.app.ErrorLog.Println("could not get reservation out of session")
		m.app.Session.Put(r.Context(), "error", "Could not get the reservation from sesssion")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	m.app.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// Contact is the contact page handler
func (m *Respository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}
