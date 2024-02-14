package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var tpl *template.Template

type PageData struct {
	ErrorMessage string
}

func init() {
	var err error
	db, err = sql.Open("mysql", "root:./@tcp(localhost:3306)/InvoicingDB")
	if err != nil {
		log.Fatal(err)
	}

	// Ensure the database connection is established
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

// NewCustomer represents a customer object
type NewCustomer struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/newInvoice", newInvoiceHandler)
	http.HandleFunc("/newCustomer", newCustomerHandler)
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{}

	data.ErrorMessage = "No error"

	tpl.ExecuteTemplate(w, "index.gohtml", data)
}

func newInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{}

	data.ErrorMessage = "No error"

	tpl.ExecuteTemplate(w, "newInvoice.gohtml", data)
}

func newCustomerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	switch r.Method {
	case "GET":
		// Display the form
		tpl.ExecuteTemplate(w, "newCustomer.gohtml", PageData{})
	case "POST":
		// Process the form submission
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		var customer NewCustomer
		customer.FirstName = r.Form.Get("FirstName")
		customer.LastName = r.Form.Get("LastName")
		customer.Email = r.Form.Get("Email")
		customer.Phone = r.Form.Get("Phone")

		insertCustomer, err := db.Prepare("INSERT INTO Customers(FirstName, LastName, Email, Phone) VALUES(?, ?, ?, ?)")
		if err != nil {
			log.Println("Database preparation error: ", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer insertCustomer.Close()

		_, err = insertCustomer.Exec(customer.FirstName, customer.LastName, customer.Email, customer.Phone)
		if err != nil {
			log.Println("Database insertion error: ", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Redirect to a success page or redisplay the form with a success message
		tpl.ExecuteTemplate(w, "newCustomer.gohtml", PageData{})
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
