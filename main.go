package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var err error

// Mahasiswa struct (Model) ...
type Mahasiswa struct {
	MahasiswaID   string `json:"MahasiswaID"`
	NamaMahasiswa string `json:"NamaMahasiswa"`
	NomorBp       string `json:"NomorBp"`
	Kelas         string `json:"Kelas"`
	Alamat        string `json:"Alamat"`
}

func getMahasiswas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var mahasiswas []Mahasiswa

	sql := `SELECT
				MahasiswaID,
				IFNULL(NamaMahasiswa,''),
				IFNULL(NomorBP,'') NomorBP,
				IFNULL(Kelas,'') Kelas,
				IFNULL(Alamat,'') Alamat
			FROM mahasiswa`

	result, err := db.Query(sql)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		var mahasiswa Mahasiswa
		err := result.Scan(&mahasiswa.MahasiswaID, &mahasiswa.NamaMahasiswa, &mahasiswa.NomorBp, &mahasiswa.Kelas, &mahasiswa.Alamat)

		if err != nil {
			panic(err.Error())
		}
		mahasiswas = append(mahasiswas, mahasiswa)
	}

	json.NewEncoder(w).Encode(mahasiswas)

}

func createMahasiswa(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		mahasiswaID := r.FormValue("MahasiswaID")
		mahasiswaNama := r.FormValue("NamaMahasiswa")
		nomorBP := r.FormValue("NomorBP")
		kelas := r.FormValue("Kelas")
		alamat := r.FormValue("Alamat")

		stmt, err := db.Prepare("INSERT INTO mahasiswa (MahasiswaID,NamaMahasiswa,NomorBP,Kelas,Alamat) VALUES (?,?,?,?,?)")

		_, err = stmt.Exec(mahasiswaID, mahasiswaNama, nomorBP, kelas, alamat)
		if err != nil {
			fmt.Fprint(w, "Data Duplicate")
		} else {
			fmt.Fprint(w, "Data Created")
		}
	}

}
func updateMahasiswa(w http.ResponseWriter, r *http.Request) {

	if r.Method == "PUT" {

		params := mux.Vars(r)

		newNamaMahasiswa := r.FormValue("NamaMahasiswa")
		newNomorBp := r.FormValue("NomorBp")
		newKelas := r.FormValue("Kelas")
		newAlamat := r.FormValue("Alamat")

		stmt, err := db.Prepare("UPDATE mahasiswa SET NamaMahasiswa = ?, NomorBp = ?, Kelas = ?, Alamat = ? WHERE MahasiswaID= ?")

		_, err = stmt.Exec(newNamaMahasiswa, newNomorBp, newKelas, newAlamat, params["id"])

		if err != nil {
			fmt.Fprint(w, "Data not found, or Request error")
		} else {

			fmt.Fprintf(w, "Mahasiswa with MahasiswaId = %s was update ", params["id"])
		}
	}
}

func deleteMahasiswa(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM mahasiswa WHERE MahasiswaId = ?")

	_, err = stmt.Exec(params["id"])
	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "Customer with ID = %s was deleted", params["id"])
}

func delMahasiswa(w http.ResponseWriter, r *http.Request) {

	MahasiswaID := r.FormValue("Mahasiswa")
	NamaMahasiswa := r.FormValue("CompanyName")

	stmt, err := db.Prepare("DELETE FROM mahasiswa WHERE MahasiswaID = ? AND CompanyName = ?")

	_, err = stmt.Exec(MahasiswaID, NamaMahasiswa)

	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "Customer with ID = %s was deleted", MahasiswaID)
}
func getMahasiswa(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var mahasiswas []Mahasiswa
	params := mux.Vars(r)

	sql := `SELECT
			MahasiswaID,
			IFNULL(NamaMahasiswa,'') NamaMahasiswa,
			IFNULL(NomorBP,'') NomorBP,
			IFNULL(Kelas,'') Kelas,
			IFNULL(Alamat,'') Alamat

			FROM mahasiswa WHERE MahasiswaID = ?`

	result, err := db.Query(sql, params["id"])

	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var mahasiswa Mahasiswa

	for result.Next() {

		err := result.Scan(&mahasiswa.MahasiswaID, &mahasiswa.NamaMahasiswa, &mahasiswa.NomorBp, &mahasiswa.Kelas, &mahasiswa.Alamat)

		if err != nil {
			panic(err.Error())
		}

		mahasiswas = append(mahasiswas, mahasiswa)
	}

	json.NewEncoder(w).Encode(mahasiswas)
}

func getPost(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var mahasiswas []Mahasiswa

	MahasiswaID := r.FormValue("MahasiswaID")
	NamaMahasiswa := r.FormValue("NamaMahasiswa")

	sql := `SELECT
	MahasiswaID,
	IFNULL(NamaMahasiswa,''),
	IFNULL(NomorBP,'') NomorBP,
	IFNULL(Kelas,'') Kelas,
	IFNULL(Alamat,'') Alamat

	FROM mahasiswa WHERE MahasiswaID = ? AND NamaMahasiswa = ?`

	result, err := db.Query(sql, MahasiswaID, NamaMahasiswa)

	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var mahasiswa Mahasiswa

	for result.Next() {

		err := result.Scan(&mahasiswa.MahasiswaID, &mahasiswa.NamaMahasiswa, &mahasiswa.NomorBp, &mahasiswa.Kelas, &mahasiswa.Alamat)

		if err != nil {
			panic(err.Error())
		}

		mahasiswas = append(mahasiswas, mahasiswa)
	}

	json.NewEncoder(w).Encode(mahasiswas)

}

func main() {

	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/challangesdb")
	if err != nil {
		panic(err.Error)
	}

	defer db.Close()

	//init router
	r := mux.NewRouter()

	//Route Handles & endpoints

	r.HandleFunc("/mahasiswa", getMahasiswas).Methods("GET")
	r.HandleFunc("/mahasiswa/{id}", getMahasiswa).Methods("GET")
	r.HandleFunc("/mahasiswa", createMahasiswa).Methods("POST")
	r.HandleFunc("/mahasiswa/{id}", deleteMahasiswa).Methods("DELETE")
	r.HandleFunc("/mahasiswa/{id}", updateMahasiswa).Methods("PUT")

	//New
	r.HandleFunc("/getMahasiswa", getPost).Methods("POST")

	//start server

	log.Fatal(http.ListenAndServe(":8282", r))
}
