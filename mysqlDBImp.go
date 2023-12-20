package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type MySqlServer struct {
	name string
}

// type DatabaseI interface {
// 	FetchAll() ([]Employee, error);
// 	FetchId(int) (Employee, error);
// 	UpdateId(int, Employee) error;
// 	Insert(*Employee) error;
// 	Delete(int) error;
// }

func (db *MySqlServer) FetchAll() ([]Employee, error) {
	dbconn := dbConn(db.name);
	defer dbconn.Close();
	selDB, err := dbconn.Query("SELECT * FROM Employee ORDER BY id DESC");
	if err != nil {
		return nil, err;
	}
	emp := Employee{};
	res := []Employee{};
	for selDB.Next() {
		var id int;
		var name, city string;
		err = selDB.Scan(&id, &name, &city);
		if err != nil {
			return nil, err;
		}
		emp.Id = id;
		emp.Name = name;
		emp.City = city;	
		res = append(res, emp);
	}
	return res, nil;
}

func (db *MySqlServer) FetchId(id int) (Employee, error) {
	dbconn := dbConn(db.name);
	defer dbconn.Close();
	selDB, err := dbconn.Query("SELECT * FROM Employee ORDER BY id DESC");
	emp := Employee{};
	if err != nil {
		return emp, err;
	}
	for selDB.Next() {
		var id int;
		var name, city string;
		err = selDB.Scan(&id, &name, &city);
		if err != nil {
			return emp, err;
		}
		emp.Id = id;
		emp.Name = name;
		emp.City = city;		
	}
	return emp, nil;
}

func (db *MySqlServer) UpdateId(id int, emp Employee) error {
	dbconn := dbConn(db.name);
	defer dbconn.Close();
	insForm, err := dbconn.Prepare("UPDATE Employee SET name=?, city=?, WHERE id=?");
	if err != nil {
		return err;
	}
	insForm.Exec(emp.Name, emp.City, id);
	return nil;
}

// 	Insert(*Employee) error;
func (db *MySqlServer) Insert(emp *Employee) error {
	dbconn := dbConn(db.name);
	defer dbconn.Close();
	insForm, err := dbconn.Prepare("INSERT INTO Employee(name, city) VALUES (?,?)");
	if err != nil {
		return err
	}
	insForm.Exec(emp.Name, emp.City);
	return nil;
}

// 	Delete(int) error;
func (db *MySqlServer) Delete(id int) error {
	dbconn := dbConn(db.name);
	defer dbconn.Close();
	delForm, err := dbconn.Prepare("DELETE FROM Employee WHERE id=?");
	if err != nil {
		return err;
	}
	delForm.Exec(id);
	return nil;
}

// func Delete(w http.ResponseWriter, r *http.Request) {
// 	db := dbConn();
// 	emp := r.URL.Query().Get("id");
// 	delForm, err := db.Prepare("DELETE FROM Employee WHERE id=?");
// 	if err != nil {
// 		panic(err.Error());
// 	}
// 	delForm.Exec(emp);
// 	log.Println("DELETE");
// 	defer db.Close();
// 	http.Redirect(w, r, "/", 301);
// }

// func Insert(w http.ResponseWriter, r *http.Request) {
// 	db := dbConn();
// 	if r.Method == "POST" {
// 		name := r.FormValue("name");
// 		city := r.FormValue("city");
// 		insForm, err := db.Prepare("INSERT INTO Employee(name, city) VALUES (?,?)")
// 		if err != nil {
// 			panic(err.Error());
// 		}
// 		insForm.Exec(name, city);
// 		log.Println("INSERT: Name: " + name + " | City: " + city);
// 	}
// 	defer db.Close();
// 	http.Redirect(w, r, "/", 301);
// }

// func Update(w http.ResponseWriter, r *http.Request) {
// 	db := dbConn();
// 	if r.Method == "POST" {
// 		name := r.FormValue("name");
// 		city := r.FormValue("city");
// 		id := r.FormValue("uid");
// 		insForm, err := db.Prepare("UPDATE Employee SET name=?, city=?, WHERE id=?");
// 		if err != nil {
// 			panic(err.Error());
// 		}
// 		insForm.Exec(name, city, id);
// 		log.Println("UPDATE: Name: " + name + " | City " + city);
// 	}
// 	defer db.Close();
// 	http.Redirect(w, r, "/", 301);
// }

// func Show(w http.ResponseWriter, r *http.Request) {
// 	db := dbConn();
// 	idStr := r.URL.Query().Get("id");
// 	id, err := strconv.Atoi(idStr);
// 	if err != nil {
// 		panic(err.Error());
// 	}
// 	res, err =
// 	selDB, err := db.Query("SELECT * FROM Employee WHERE id=?", id);
// 	if err != nil {
// 		panic(err.Error());
// 	}
// 	emp := Employee{}
// 	for selDB.Next() {
// 		var id int;
// 		var name, city string;
// 		err = selDB.Scan(&id, &name, &city);
// 		if err != nil {
// 			panic(err.Error());
// 		}
// 		emp.Id = id;
// 		emp.Name = name;
// 		emp.City = city;
// 	}
// 	tmpl.ExecuteTemplate(w, "Show", emp);
// 	defer db.Close();
// }
const (
	username = "root"
	password = ""
	hostname = "127.0.0.1:3306"
	dbname   = "goblog"
)

func dsn(name string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, name)
}

func dbConn(name string) *sql.DB {
	db, err := sql.Open("mysql", dsn(name));
	if err != nil {
		panic(err.Error());
	}
	return db;
}

func createDatabase(db *sql.DB, name string) sql.Result {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second);
	defer cancelfunc();
	res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+name);
	if err != nil {
		log.Panicf("Error %s when creating DB\n", err);
	}
	return res;
}

// func createTableQueryFromStructure(db *sql.DB, name string, data interface{}) string {
// 	var columns []string;
// 	t := reflect.TypeOf(data);
// 	for i := 0; i < t.NumField(); i++ {
// 		field := t.Field(i);
// 		columnName := strings.ToLower(field.Name);
// 		columnType := "";
// 		switch field.Type.Kind() {
// 		case reflect.Int:
// 			columnType = "INT";
// 		case reflect.String:
// 			columnType = "VARCHAR(255)"
// 		default:
// 			panic("createTableQueryFromStructure: Unsupported structure type");
// 		}
// 		columns = append(columns, fmt.Sprintf("%s %s", columnName, columnType));
// 	}
// 	createTableQuery := fmt.Sprintf("CREATE TABLE %s (%s);", name, strings.Join(columns, ", "));
// 	return createTableQuery;
// }

func createTable(db *sql.DB, name string) sql.Result {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second);
	defer cancelfunc();
	// res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+name);
	dropQuery := "DROP TABLE IF EXISTS`"+name+"`;";
	query := 
		"CREATE TABLE `"+name+"` (" +
			"`id` int(6) unsigned NOT NULL AUTO_INCREMENT," +
			"`name` varchar(30) NOT NULL," +
			"`city` varchar(30) NOT NULL," +
			"PRIMARY KEY (`id`)" +
		") ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;"
	res, err := db.ExecContext(ctx, dropQuery);
	res, err = db.ExecContext(ctx, query);
	if err != nil {
		log.Panicf("Error %s when creating table\n", err);
	}
	return res;
}

func setDBParams(db *sql.DB) {
	db.SetMaxOpenConns(20);
	db.SetMaxIdleConns(20);
	db.SetConnMaxLifetime(time.Minute * 5);
}