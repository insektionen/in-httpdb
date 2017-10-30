package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type log_record struct {
	id        int    `json:"-"`
	Name      string `json:"name"`
	Category  string `json:"category"`
	Timestamp int64  `json:"timestamp"`
	Count     int    `json:"count"`
}

type stock_entry struct {
	Name            string `json:"name"`
	Category        string `json:"category"`
	Supplier        string `json:"supplier"`
	Barcode         string `json:"barcode"`
	Package_barcode string `json:"package_barcode"`
	Package_size    int    `json:"package_size"`
	Count           int    `json:"count"`
	Visible         int    `json:"visible"`
}

type SumEntry struct {
	Category string `json:"category"`
	Sum      int64  `json:"sum"`
}

// Stock
//   get: name, barcode, category, supplier
func http_Stock(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		log.Println("Get stock")

		name := "%"
		barcode := "%"
		category := "%"
		supplier := "%"

		if nam := r.FormValue("name"); nam != "" {
			name = nam
		}

		if bc := r.FormValue("barcode"); bc != "" {
			barcode = bc
		}

		if cat := r.FormValue("category"); cat != "" {
			category = cat
		}

		if supp := r.FormValue("supplier"); supp != "" {
			supplier = supp
		}

		rows, e := db.Query(`SELECT name, category, supplier, barcode, package_barcode, package_size, itemCount, visible FROM stock 
							WHERE name LIKE ? AND category LIKE ? AND supplier LIKE ? AND (barcode LIKE ? OR package_barcode LIKE ?)`,
			name, category, supplier, barcode, barcode)

		if e != nil {
			log.Println(e)
			return
		}
		defer rows.Close()

		res := make([]stock_entry, 0, 100)
		for rows.Next() {
			ent := stock_entry{}
			e := rows.Scan(&ent.Name, &ent.Category, &ent.Supplier, &ent.Barcode, &ent.Package_barcode, &ent.Package_size, &ent.Count, &ent.Visible)
			if e != nil {
				log.Println(e)
			}
			res = append(res, ent)
		}
		if err := rows.Err(); err != nil {
			log.Println(err)
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		encoder := json.NewEncoder(w)
		encoder.Encode(res)

	case "PUT":
		log.Println("Update stock")
	default:
		log.Println("Invalid method")
	}
}

func http_Increment(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		log.Println("POST Increment")

		name := "%"
		barcode := "%"
		amount := ""

		if nam := r.PostFormValue("name"); nam != "" {
			name = nam
		}
		if bc := r.PostFormValue("barcode"); bc != "" {
			barcode = bc
		}
		if i := r.PostFormValue("amount"); i != "" {
			amount = i
		}
		_, e := db.Exec(`UPDATE stock 
						SET itemCount=itemCount+? 
						WHERE name LIKE ? 
						AND (barcode LIKE ? OR package_barcode LIKE ?)`,
			amount, name, barcode, barcode)

		if e != nil {
			log.Println(e)
			return
		}

		rows, e := db.Query(`SELECT name, category, supplier, barcode, package_barcode, package_size, itemCount, visible FROM stock 
							WHERE name LIKE ? AND (barcode LIKE ? OR package_barcode LIKE ?)`,
			name, barcode, barcode)

		if e != nil {
			log.Println(e)
			return
		}
		defer rows.Close()

		ent := stock_entry{}
		rows.Next()
		e = rows.Scan(&ent.Name, &ent.Category, &ent.Supplier, &ent.Barcode, &ent.Package_barcode, &ent.Package_size, &ent.Count, &ent.Visible)

		if e != nil {
			log.Println(e)
		}
		if err := rows.Err(); err != nil {
			log.Println(err)
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		encoder := json.NewEncoder(w)
		encoder.Encode(ent)

	default:
		log.Println("Invalid method")
	}
}

func http_Decrement(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		log.Println("POST Decrement")

		name := "%"
		barcode := "%"
		amount := ""

		if nam := r.PostFormValue("name"); nam != "" {
			name = nam
		}
		if bc := r.PostFormValue("barcode"); bc != "" {
			barcode = bc
		}
		if i := r.PostFormValue("amount"); i != "" {
			k, _ := strconv.Atoi(i)
			if k < 0 {
				amount = strconv.Itoa(k * -1)
			} else {
				amount = i
			}
		}
		_, e := db.Exec(`UPDATE stock SET itemCount=itemCount-? WHERE name LIKE ? AND (barcode LIKE ? OR package_barcode LIKE ?)`,
			amount, name, barcode, barcode)
		if e != nil {
			log.Println(e)
			return
		}

		rows, e := db.Query(`SELECT name, category, supplier, barcode, package_barcode, package_size, itemCount, visible FROM stock 
							WHERE name LIKE ? AND (barcode LIKE ? OR package_barcode LIKE ?)`,
			name, barcode, barcode)
		if e != nil {
			log.Println(e)
			return
		}
		defer rows.Close()

		ent := stock_entry{}
		rows.Next()
		e = rows.Scan(&ent.Name, &ent.Category, &ent.Supplier, &ent.Barcode, &ent.Package_barcode, &ent.Package_size, &ent.Count, &ent.Visible)

		if e != nil {
			log.Println(e)
		}
		if err := rows.Err(); err != nil {
			log.Println(err)
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		encoder := json.NewEncoder(w)
		encoder.Encode(ent)
	default:
		log.Println("Invalid method")
	}
}

func http_UpdateItem(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		log.Println("POST UpdateItem")

		name := ""
		barcode := ""
		package_barcode := ""
		package_size := ""
		supplier := ""
		visible := ""

		if nam := r.PostFormValue("name"); nam != "" {
			name = nam
		}
		if bc := r.PostFormValue("barcode"); bc != "" {
			barcode = bc
		}
		if pbc := r.PostFormValue("packageBarcode"); pbc != "" {
			package_barcode = pbc
		}
		if psize := r.PostFormValue("packageSize"); psize != "" {
			package_size = psize
		}
		if sup := r.PostFormValue("supplier"); sup != "" {
			supplier = sup
		}
		if vis := r.PostFormValue("visible"); vis != "" {
			visible = vis
		}

		_, e := db.Exec(`UPDATE stock SET barcode=?, package_barcode=?, package_size=?, supplier=?, visible=? WHERE name=?`,
			barcode, package_barcode, package_size, supplier, visible, name)

		if e != nil {
			log.Println(e)
			return
		}
	default:
		log.Println("Invalid method")
	}
}

// Log
//   get: log date_start,date_end,category,name
func http_Log(w http.ResponseWriter, r *http.Request) {
	date_end := time.Now()
	start_date := time.Unix(0, 0)
	category := "%"
	name := "%"

	if cat := r.FormValue("category"); cat != "" {
		category = cat
	}

	if nam := r.FormValue("name"); nam != "" {
		name = nam
	}

	if stamp := r.FormValue("date_start"); stamp != "" {
		tmp, _ := strconv.ParseInt(stamp, 10, 64)
		start_date = time.Unix(tmp, 0)
	}
	if stamp := r.FormValue("date_end"); stamp != "" {
		tmp, _ := strconv.ParseInt(stamp, 10, 64)
		date_end = time.Unix(tmp, 0)
	}

	log.Println("Fetching log between", start_date.Unix(), "to", date_end.Unix())
	rows, e := db.Query(`SELECT id, name, itemCount, UNIX_TIMESTAMP(timestamp), category FROM log WHERE
						name like ? AND category like ? AND timestamp >= FROM_UNIXTIME(?) AND timestamp <= FROM_UNIXTIME(?)`,
		name, category, start_date.Unix(), date_end.Unix())
	if e != nil {
		log.Println(e)
		w.Write([]byte(e.Error()))
		return
	}
	defer rows.Close()

	res := make([]log_record, 0, 100)

	for rows.Next() {
		rec := log_record{}
		e := rows.Scan(&rec.id, &rec.Name, &rec.Count, &rec.Timestamp, &rec.Category)
		if e != nil {
			log.Println(e)
		}
		res = append(res, rec)

	}
	if err := rows.Err(); err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	encoder := json.NewEncoder(w)
	encoder.Encode(res)
}

func HttpCategories(w http.ResponseWriter, r *http.Request) {
	res, err := getCategories()

	if err != nil {
		log.Println("AAA", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func HttpSumCategories(w http.ResponseWriter, r *http.Request) {
	dateEnd := time.Now()
	dateStart := time.Unix(0, 0)

	if stamp := r.FormValue("date_start"); stamp != "" {
		tmp, _ := strconv.ParseInt(stamp, 10, 64)
		dateStart = time.Unix(tmp, 0)
	}
	if stamp := r.FormValue("date_end"); stamp != "" {
		tmp, _ := strconv.ParseInt(stamp, 10, 64)
		dateEnd = time.Unix(tmp, 0)
	}

	res, err := sumAllCategories(dateStart, dateEnd)

	if err != nil {
		log.Println("LogSum:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println("SumCategories:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func HttpLatestAdded(w http.ResponseWriter, r *http.Request) {
	category := "%"
	count := 5

	if cat := r.FormValue("category"); cat != "" {
		category = cat
	}

	if cou := r.FormValue("count"); cou != "" {
		tmp, _ := strconv.Atoi(cou)
		count = tmp
	}

	res, err := getLatestAdded(category, count)

	if err != nil {
		log.Println("LatestAdded:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println("LatestAdded:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func HttpTop(w http.ResponseWriter, r *http.Request) {
	category := "%"
	dateEnd := time.Now()
	dateStart := time.Unix(0, 0)

	if cat := r.FormValue("category"); cat != "" {
		category = cat
	}

	if stamp := r.FormValue("date_start"); stamp != "" {
		tmp, _ := strconv.ParseInt(stamp, 10, 64)
		dateStart = time.Unix(tmp, 0)
	}

	if stamp := r.FormValue("date_end"); stamp != "" {
		tmp, _ := strconv.ParseInt(stamp, 10, 64)
		dateEnd = time.Unix(tmp, 0)
	}

	res, err := getTop(category, dateStart, dateEnd)

	if err != nil {
		log.Println("LogSum:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println("LogSum:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func HttpLogSum(w http.ResponseWriter, r *http.Request) {
	dateEnd := time.Now()
	dateStart := time.Unix(0, 0)

	if stamp := r.FormValue("date_start"); stamp != "" {
		tmp, _ := strconv.ParseInt(stamp, 10, 64)
		log.Println("Date from:", tmp)
		dateStart = time.Unix(tmp, 0)
	}
	if stamp := r.FormValue("date_end"); stamp != "" {
		tmp, _ := strconv.ParseInt(stamp, 10, 64)
		log.Println("Date to:", tmp)
		dateEnd = time.Unix(tmp, 0)
	}

	category := "%"
	name := "%"

	if nam := r.FormValue("name"); nam != "" {
		name = nam
	}
	if cat := r.FormValue("category"); cat != "" {
		category = cat
	}

	res, err := logSum(category, name, dateStart, dateEnd)

	if err != nil {
		log.Println("LogSum:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println("LogSum:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func getCategories() ([]string, error) {
	rows, err := db.Query(`SELECT DISTINCT category FROM stock`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]string, 0, 100)

	for rows.Next() {
		var str string
		err := rows.Scan(&str)
		if err != nil {
			return nil, err
		}
		res = append(res, str)

	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func sumAllCategories(dateStart time.Time, dateEnd time.Time) ([]SumEntry, error) {
	rows, err := db.Query(`SELECT category, IFNULL(SUM(itemCount), 0) as "sum" FROM log WHERE timestamp > ? GROUP BY category`, dateStart)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]SumEntry, 0, 100)
	for rows.Next() {
		ent := SumEntry{}
		err = rows.Scan(&ent.Category, &ent.Sum)

		res = append(res, ent)
	}

	return res, nil
}

func getTop(category string, dateStart time.Time, dateEnd time.Time) ([]string, error) {
	rows, err := db.Query(`SELECT name, IFNULL(SUM(itemCount), 0) as "sum" FROM log WHERE category LIKE ? AND timestamp > ? GROUP BY name ORDER BY sum DESC LIMIT 10`, category, dateStart)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]string, 0, 100)
	for rows.Next() {
		var name string
		var sum string
		err = rows.Scan(&name, &sum)

		res = append(res, name)
	}

	return res, nil
}

func logSum(category string, name string, dateStart time.Time, dateEnd time.Time) (*SumEntry, error) {
	var row *sql.Row
	if name == "%" {
		row = db.QueryRow(`SELECT IFNULL(SUM(itemCount), 0) as sum FROM log WHERE category LIKE ? AND timestamp > ?`,
			category, dateStart)
	} else {
		row = db.QueryRow(`SELECT IFNULL(SUM(itemCount), 0) as sum FROM log WHERE name LIKE ? AND timestamp > ?`,
			"%"+name+"%", dateStart)
	}
	ent := SumEntry{category, 0}

	err := row.Scan(&ent.Sum)
	if err == sql.ErrNoRows {
		return &ent, nil
	} else if err != nil {
		return nil, err
	}
	return &ent, nil
}

func getLatestAdded(category string, count int) ([]string, error) {
	rows, err := db.Query(`SELECT name FROM stock WHERE category LIKE ? ORDER BY added DESC LIMIT ?`, category, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]string, 0, 100)

	for rows.Next() {
		var str string
		err := rows.Scan(&str)
		if err != nil {
			return nil, err
		}
		res = append(res, str)

	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func main() {
	username := os.Getenv("MYSQL_USERNAME")
	if username == "" {
		log.Fatal("Set MYSQL_USERNAME env")
	}

	password := os.Getenv("MYSQL_PASSWORD")
	if password == "" {
		log.Fatal("Set MYSQL_PASSWORD env")
	}
	ip := os.Getenv("MYSQL_IP")
	if ip == "" {
		log.Fatal("Set MYSQL_IP env")
	}
	port := os.Getenv("MYSQL_PORT")
	if port == "" {
		log.Fatal("Set MYSQL_PORT env")
	}
	databaseName := os.Getenv("MYSQL_DATABASE")
	if databaseName == "" {
		log.Fatal("Set MYSQL_DATABASE env")
	}
	bind := os.Getenv("BIND")
	if bind == "" {
		log.Fatal("Set BIND env")
	}

	var err error
	path := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, ip, port, databaseName)
	db, err = sql.Open("mysql", path)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Starting...")

	http.HandleFunc("/log", http_Log)
	http.HandleFunc("/stock", http_Stock)
	http.HandleFunc("/stock/increment", http_Increment)
	http.HandleFunc("/stock/decrement", http_Decrement)
	http.HandleFunc("/stock/updateItem", http_UpdateItem)
	http.HandleFunc("/sum", HttpLogSum)
	http.HandleFunc("/sum/categories", HttpSumCategories)
	http.HandleFunc("/categories", HttpCategories)
	http.HandleFunc("/top", HttpTop)
	http.HandleFunc("/latest/added", HttpLatestAdded)

	log.Println(http.ListenAndServe(bind, nil))
}
