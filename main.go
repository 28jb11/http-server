package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var tpl *template.Template

type Customer struct {
	CustomerID int
	FirstName  string
	LastName   string
	Email      string
	Phone      string
}

type PageData struct {
	ErrorMessage string
	Customers    []Customer
}

func init() {
	var err error
	db, err = sql.Open("mysql", "root:./@tcp(localhost:3306)/InvoicingDB")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/invoices", newInvoiceHandler)
	http.HandleFunc("/customers", customerHandler)
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{}

	data.ErrorMessage = "This is an error message. Nothing is wrong."

	tpl.ExecuteTemplate(w, "index.gohtml", data)
}

func newInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{}

	data.ErrorMessage = "No invoice page yet."

	tpl.ExecuteTemplate(w, "index.gohtml", data)
}

func customerHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	switch r.Method {
	case "GET":
		fmt.Println("customer GET")

		// Pull customer table from database
		rows, err := db.Query("SELECT CustomerID, FirstName, LastName, Email, Phone FROM Customers")
		if err != nil {
			fmt.Println("db.Query error")
			http.Error(w, "Database query error", http.StatusInternalServerError)
			data.ErrorMessage = "Database query error"
			tpl.ExecuteTemplate(w, "customers.gohtml", data)
			return
		}
		defer rows.Close()

		var customers []Customer
		for rows.Next() {
			var customer Customer
			err := rows.Scan(&customer.CustomerID, &customer.FirstName, &customer.LastName, &customer.Email, &customer.Phone)
			if err != nil {
				http.Error(w, "Error scanning customer from database", http.StatusInternalServerError)
				data.ErrorMessage = "Error scanning customer from database"
				tpl.ExecuteTemplate(w, "customers.gohtml", data)
				return
			}
			customers = append(customers, customer)
		}

		// Pass the customers to the template
		tpl.ExecuteTemplate(w, "customers.gohtml", PageData{Customers: customers})

	case "POST":
		fmt.Println("customer POST")

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		var customer Customer
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

		// Redirect to the same endpoint to refresh the page
		http.Redirect(w, r, "/customers", http.StatusSeeOther)
	default:
		fmt.Printf("Received unexpected method: %s\n", r.Method)
		fmt.Println("default case")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
