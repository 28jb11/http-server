package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

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
	http.HandleFunc("/customers/edit", editCustomerHandler)
	http.HandleFunc("/customers/", editCustomerFormHandler)
	http.HandleFunc("/customers/save", saveCustomerHandler)
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{}

	data.ErrorMessage = "This is an error message, but nothing is wrong."

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

func editCustomerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	customerID := r.URL.Query().Get("id")
	if customerID == "" {
		http.Error(w, "Bad Request: Missing customer ID", http.StatusBadRequest)
		return
	}

	// Fetch customer details from the database based on the provided customer ID
	row := db.QueryRow("SELECT CustomerID, FirstName, LastName, Email, Phone FROM Customers WHERE CustomerID = ?", customerID)
	var customer Customer
	err := row.Scan(&customer.CustomerID, &customer.FirstName, &customer.LastName, &customer.Email, &customer.Phone)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error fetching customer details:", err)
		return
	}

	// Render the edit customer form
	tpl.ExecuteTemplate(w, "edit_customer.gohtml", customer)
}

func editCustomerFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the customer ID from the URL
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Bad Request: Missing customer ID", http.StatusBadRequest)
		return
	}
	customerID := parts[2]

	// Fetch customer details from the database based on the provided customer ID
	row := db.QueryRow("SELECT CustomerID, FirstName, LastName, Email, Phone FROM Customers WHERE CustomerID = ?", customerID)
	var customer Customer
	err := row.Scan(&customer.CustomerID, &customer.FirstName, &customer.LastName, &customer.Email, &customer.Phone)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error fetching customer details:", err)
		return
	}

	// Render the edit customer form
	tpl.ExecuteTemplate(w, "edit_customer.gohtml", customer)
}

func saveCustomerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	} else {
		fmt.Println("saveCustomerHandler POST")
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println("Error parsing form data:", err)
		return
	}

	// Extract customer details from the form
	customerID := r.PostForm.Get("CustomerID")
	firstName := r.PostForm.Get("EditFirstName")
	lastName := r.PostForm.Get("EditLastName")
	email := r.PostForm.Get("EditEmail")
	phone := r.PostForm.Get("EditPhone")

	fmt.Println(customerID, firstName, lastName, email, phone)
	// Prepare the SQL query
	stmt, err := db.Prepare("UPDATE Customers SET FirstName=?, LastName=?, Email=?, Phone=? WHERE CustomerID=?")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error preparing SQL statement:", err)
		return
	}
	defer stmt.Close()

	// Execute the prepared statement with the provided parameters

	fmt.Println(customerID, firstName, lastName, email, phone)
	_, err = stmt.Exec(firstName, lastName, email, phone, customerID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error updating customer details:", err)
		return
	}

	// Redirect back to the customer list page
	http.Redirect(w, r, "/customers", http.StatusSeeOther)
}
