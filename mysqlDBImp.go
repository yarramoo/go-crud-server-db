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

func createTable(db *sql.DB, name string) sql.Result {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second);
	defer cancelfunc();
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