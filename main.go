// https://golangbot.com/connect-create-db-mysql/
// ^ Used to connect to sql server and create database

// https://www.golangprograms.com/example-of-golang-crud-using-mysql-from-scratch.html
// ^ Used to structure a CRUD app
package main 

import (
	"net/http"
)

// What can a database do?
// Make a new one
// Add an entry
// Update an entry
// Remove an entry
// Query
// Let's hardcode the entry type for now
type DatabaseI interface {
	FetchAll() ([]Employee, error);
	FetchId(int) (Employee, error);
	UpdateId(int, Employee) error;
	Insert(*Employee) error;
	Delete(int) error;
}

type Employee struct {
	Id   int
	Name string
	City string
}

func initDB() {
	db := dbConn("");
	_ = createDatabase(db, "goblog");
	db.Close();
	db = dbConn("goblog");
	_ = createTable(db, "employee");
}


func main() {
	initDB();
	sqlDB := MySqlServer{
		name: "goblog",
	}

	http.HandleFunc("/", IndexWrapper(&sqlDB));
	http.HandleFunc("/show", ShowWrapper(&sqlDB));
	http.HandleFunc("/new", New);
	http.HandleFunc("/edit", EditWrapper(&sqlDB));
	http.HandleFunc("/insert", InsertWrapper(&sqlDB));
	http.HandleFunc("/update", UpdateWrapper(&sqlDB));
	http.HandleFunc("/delete", DeleteWrapper(&sqlDB));
	http.ListenAndServe(":8080", nil);
}