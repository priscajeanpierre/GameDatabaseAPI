package main

import (
	"database/sql"
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"net/http"
	"strconv"
)

type GameData struct {
	Name            string
	ReleaseDate     int
	ReccCount       int
	SteamSpyOwners  int
	SteamSpyPlayers int
	Windows         bool
	Linux           bool
	Mac             bool
}

func main() {
	//do I actually need this

	mygameDatabase := OpenDataBase("games-features.xlsx")
	defer mygameDatabase.Close()
	createTables(mygameDatabase)
	addData(mygameDatabase)
	addGameData(mygameDatabase)

}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	gamedata := GameData{}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
}

func rows() string {
	gameExcelFile, err := excelize.OpenFile("games-features.xlsx")
	if err != nil {
		log.Fatalln(err)
	}
	allRows, err := gameExcelFile.GetRows("games-features")
	if err != nil {
		log.Fatalln(err)
	}
	for number, row := range allRows {
		if number <= 0 {
			continue
		}
		if len(row) <= 1 {
			continue
		}

		if row[0] == "10" {
			fmt.Println("Game Name: ", row[2])
			fmt.Println("Release Date: ", row[4])
			fmt.Println("Recommendation Count: ", row[12])
			fmt.Println("Steam Spy Owners: ", row[15])
			fmt.Println("Steam Spy Players: ", row[17])
			fmt.Println("Windows: ", row[26])
			fmt.Println("Linux: ", row[27])
			fmt.Println("Mac: ", row[28])

		}
	}
	return rows()
}

func OpenDataBase(dbfile string) *sql.DB {
	gdatabase, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		log.Fatal(err)
	}
	return gdatabase
}

func createTables(gdatabase *sql.DB) {
	gameCreateStatement := "CREATE TABLE IF NOT EXISTS gameInfo(    " +
		"Name TEXT PRIMARY KEY," +
		"release_Date INTEGER DEFAULT 0" +
		"required_Age INTEGER DEFAULT 0, " +
		"rec_Count INTEGER DEFAULT 0, " +
		"steamSpy_owners INTEGER DEFAULT 0, " +
		"steamSpy_players INTEGER DEFAULT 0, " +
		"platforms TEXT NOT NULL);"
	gdatabase.Exec(gameCreateStatement)
}

func addData(gdatabase *sql.DB) {

	//use all rows
}

func addGameData(gdatabase *sql.DB) {
	//LOOP through database
	insertStatment := "INSERT INTO GAMEINFO (Name, release_date, required_age, rec_Count,  steamSpy_owners, steamSpy_players, platforms) " +
		"VALUES (?,?,?,?,?,?,?);"
	for gameData, desc := range rows() {
		Name := gameData[:4]
		release_date := gameData[4:]
		intGameData, err := strconv.Atoi(gameData)

		prepped_statement, err := gdatabase.Prepare(insertStatment)
		if err != nil {
			log.Fatalln(err)
		}
		prepped_statement.Exec(Name, release_date, intGameData, desc)

	}

}
