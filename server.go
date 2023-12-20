package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

var tmpl = template.Must(template.ParseGlob("form/*"));

func IndexWrapper(db DatabaseI) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := db.FetchAll();
		fmt.Println(res);
		if err != nil {
			panic(err.Error());
		}
		tmpl.ExecuteTemplate(w, "Index", res);
	}
}

func ShowWrapper(db DatabaseI) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id");
		id, err := strconv.Atoi(idStr);
		if err != nil {
			panic(err.Error());
		}
		res, err := db.FetchId(id);
		if err != nil {
			panic(err.Error());
		}
		tmpl.ExecuteTemplate(w, "Show", res);
	}
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil);
}

func EditWrapper(db DatabaseI) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id");
		id, err := strconv.Atoi(idStr);
		if err != nil {
			panic(err.Error());
		}
		res, err := db.FetchId(id);
		if err != nil {
			panic(err.Error());
		}
		tmpl.ExecuteTemplate(w, "Edit", res);
	}	
}

// func Edit(w http.ResponseWriter, r *http.Request) {
// 	db := dbConn();
// 	nId := r.URL.Query().Get("id");
// 	selDB, err := db.Query("SELECT * FROM Employee WHERE id=?", nId);
// 	if err != nil {
// 		panic(err.Error());
// 	}
// 	emp := Employee{};
// 	for selDB.Next() {
// 		var id int;
// 		var name, city string;
// 		err = selDB.Scan(&id, &name, &city);
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		emp.Id = id
//         emp.Name = name
//         emp.City = city
//     }
//     tmpl.ExecuteTemplate(w, "Edit", emp)
//     defer db.Close()
// }

func InsertWrapper(db DatabaseI) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			name := r.FormValue("name");
			city := r.FormValue("city");
			emp := Employee {
				Name: name,
				City: city,
			};
			err := db.Insert(&emp);
			if err != nil {
				panic(err.Error());
			}
			log.Println("INSERT: Name: " + name + " | City: " + city);
		}
		http.Redirect(w, r, "/", 301);

	}	
}

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

func UpdateWrapper(db DatabaseI) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			name := r.FormValue("name");
			city := r.FormValue("city");
			emp := Employee {
				Name: name,
				City: city,
			};
			id, err := strconv.Atoi(r.FormValue("uid"));
			if err != nil {
				panic(err.Error());
			}
			err = db.UpdateId(id, emp);
			if err != nil {
				panic(err.Error());
			}
			log.Println("UPDATE: Name: " + name + " | City " + city);
		}
		http.Redirect(w, r, "/", 301);

	}	
}

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

func DeleteWrapper(db DatabaseI) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"));
		if err != nil {
			panic(err.Error());
		}
		err = db.Delete(id);
		if err != nil {
			panic(err.Error());
		}
		log.Println("DELETE");
		http.Redirect(w, r, "/", 301);
	}
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