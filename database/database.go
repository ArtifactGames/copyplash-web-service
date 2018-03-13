package database

import (
	"database/sql"
	"log"
	"math/rand"

	_ "github.com/mattn/go-sqlite3" //SQLite driver
)

const (
	databaseFile = "./data.db"
)

// PrepareDatabase creates the necessary database file and tables.
func PrepareDatabase() {
	db, _ := sql.Open("sqlite3", databaseFile)

	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS lobbies (id INTEGER PRIMARY KEY AUTOINCREMENT, password TEXT)")
	statement.Exec()

	statement, _ = db.Prepare("CREATE TABLE IF NOT EXISTS lobby_player (cookie_id INTEGER PRIMARY KEY, lobby_id INTEGER, FOREIGN KEY(lobby_id) references lobbies(id))")
	statement.Exec()

	statement, _ = db.Prepare("DELETE FROM lobbies")
	statement.Exec()

	statement, _ = db.Prepare("DELETE FROM lobby_player")
	statement.Exec()
}

// CreateLobby creates a lobby for players to join
func CreateLobby() (int64, int) {
	db, _ := sql.Open("sqlite3", databaseFile)

	var randomCode int

	randomCode = rand.Intn(99999)

	rows, _ := db.Query("SELECT * FROM lobbies WHERE password = $1", randomCode)

	for rows.Next() {
		randomCode = rand.Intn(99999)
		rows, _ = db.Query("SELECT * FROM lobbies WHERE password = $1", randomCode)
	}

	statement, _ := db.Prepare("INSERT INTO lobbies(password) values(?)")
	result, _ := statement.Exec(randomCode)
	lastInsertedID, _ := result.LastInsertId()
	return lastInsertedID, randomCode
}

// DestroyLobby destroys an already created lobby
func DestroyLobby(id int) {
	db, _ := sql.Open("sqlite3", databaseFile)

	statement, _ := db.Prepare("DELETE FROM lobbies WHERE id = ?")
	statement.Exec(id)
}

// EnterLobby destroys an already created lobby
func EnterLobby(clientID string, password int) (bool) {
	db, err := sql.Open("sqlite3", databaseFile)

	if err != nil {
		log.Println("cannot open the database")
		return false
	}
	
	var lobbyID int
	err = db.QueryRow("SELECT id FROM lobbies WHERE password = $1", password).Scan(&lobbyID)
	if err != nil {
		return false
	}

	statement, err := db.Prepare("INSERT INTO lobby_player(cookie_id, lobby_id) values(?, ?)");
	if err != nil {
		log.Println("Cannot prepare the insert query")
		return false
	}

	_, err = statement.Exec(clientID, lobbyID)
	
	if err != nil {
		log.Println("Cannot insert the player")
	}

	return err == nil
}

// func borralla() {
// 	statement, _ = database.Prepare("INSERT INTO people (firstname, lastname) VALUES (?, ?)")
// 	statement.Exec("Nic", "Raboy")
// 	rows, _ := database.Query("SELECT id, firstname, lastname FROM people")
// 	var id int
// 	var firstname string
// 	var lastname string
// 	for rows.Next() {
// 		rows.Scan(&id, &firstname, &lastname)
// 		fmt.Println(strconv.Itoa(id) + ": " + firstname + " " + lastname)
// 	}
// }