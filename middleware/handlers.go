package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Go-Lang/go-PostgreSQL/models"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func CreateConnection() *sql.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env files")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to postgres")
	return db

}

func CreateStock(w http.ResponseWriter, r *http.Request) {
	var stock models.Stock
	//we are decoding the body coming into jason and then storing it in stocks
	err := json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		log.Fatalf("Unable to decode the request body. %v", err)
	}

	insertID := insertStock(stock)

	res := response{
		ID:      insertID,
		Message: "stock created succesfully",
	}
	json.NewEncoder(w).Encode(res)

}

func GetStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int. %v", err)
	}

	stock, err := getstock(int64(id))
	if err != nil {
		log.Fatalf("unable to get stock. %v", err)
	}

	json.NewEncoder(w).Encode(stock)

}

func UpdateStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("unable to convert string to integer. %v", err)

	}

	var stock models.Stock

	err = json.NewDecoder(r.Body).Decode(stock)
	if err != nil {
		log.Fatalf("unable to ecode the request body. %v", err)
	}

	updatedRows := updateStock(int64(id), stock)

	msg := fmt.Sprintf("Stocks updated successfully. total rows/records affected %v", updatedRows)
	res := response{
		ID:      int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)

}

func GetAllStock(w http.ResponseWriter, r *http.Request) {
	stock, err := getAllStock()
	if err != nil {
		log.Fatalf("unable to get all stocks %v", err)

	}
	json.NewEncoder(w).Encode(stock)

}

func DeleteStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("unable to convert string to int. %v", err)
	}

	deletedRows := deleteStock(int64(id))
	msg := fmt.Sprintf("Stock deleted successfully. total rows/records %v", deletedRows)

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)

}

//the below functoin the toose which will work with databases

func insertStock(stock models.Stock) int64 {
	db := CreateConnection()
	defer db.Close()
	sqlStatement := `INSERT INTO stocks(name, price, company) VALUES ($1, $2, $3) RETURNING stocksid`
	var id int64

	err := db.QueryRow(sqlStatement, stock.Name, stock.Price, stock.Company).Scan(&id)
	if err != nil {
		log.Fatal("Unable to execute the query. %v", err)
	}

	return id

}

// here id is the parameter and model.stocks is the return type
func getstock(id int64) (models.Stock, error) {
	db := CreateConnection()
	defer db.Close()
	var stock models.Stock

	sqlStatement := `SELECT * FROM stocks WHERE stockid=$1`

	row := db.QueryRow(sqlStatement, id)

	err := row.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("Now rows were returned!")
		return stock, nil

	case nil:
		return stock, nil

	default:
		log.Fatal("Unable to scan the row. %v", err)

	}
	return stock, err

}

func getAllStock() ([]models.Stock, error) {
	db := CreateConnection()
	defer db.Close()

	var stocks []models.Stock
	sqlStatement := `SELECT * FROM stocks`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatalf("Unable to execute the quesry.%v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var stock models.Stock
		err = rows.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)
		if err != nil {
			log.Fatalf("unable to scan the row %v", err)
		}
		stocks = append(stocks, stock)
	}

	return stocks, err

}

func updateStock(id int64, stock models.Stock) int64 {
	db := CreateConnection()

	// close the db connection
	defer db.Close()

	// create the update sql query
	sqlStatement := `UPDATE stocks SET name=$2, price=$3, company=$4 WHERE stockid=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, id, stock.Name, stock.Price, stock.Company)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected

}

func deleteStock(id int64) int64 {
	// create the postgres db connection
	db := CreateConnection()

	// close the db connection
	defer db.Close()

	// create the delete sql query
	sqlStatement := `DELETE FROM stocks WHERE stockid=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected

}
