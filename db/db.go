package db

import (

	"log"
	"database/sql"

	"github.com/jmoiron/sqlx"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	  
	"github.com/RainrainWu/gopherbot/config"
)

var schema = `
DROP TABLE IF EXISTS resource_to_team;
DROP TABLE IF EXISTS resource;
DROP TABLE IF EXISTS team;

CREATE TABLE resource (
	id INT GENERATED ALWAYS AS IDENTITY,
	name VARCHAR(255) NOT NULL,
	url VARCHAR(1023) NOT NULL,
	PRIMARY KEY(id)
);

CREATE TABLE team (
	id INT GENERATED ALWAYS AS IDENTITY,
	name VARCHAR(255) NOT NULL,
	PRIMARY KEY(id)
);

CREATE TABLE resource_to_team (
	id INT GENERATED ALWAYS AS IDENTITY,
	resource_id INT,
	team_id INT,
	CONSTRAINT fk_resource
		FOREIGN KEY(resource_id)
			REFERENCES resource(id)
			ON DELETE CASCADE,
	CONSTRAINT fk_team
		FOREIGN KEY(team_id)
			REFERENCES team(id)
			ON DELETE CASCADE
);
`

type Resource struct {

	Id		int 	`db:"id"`
	Name	string 	`db:"name"`
	Url 	string 	`db:"url"`
}

type ResourceToTeam struct {

	Id			int `db:"id"`
	ResourceId	int	`db:"resource_id"`
	TeamId		int	`db:"team_id"`
}

type Team struct {

	Id		int		`db:"id"`
	Name	string	`db:"name"`
}

var (

	database *sqlx.DB
	err interface{}
)

func init() {

	ConnectDatabase()
}

func ConnectDatabase()  {

	connStr := "host=" + config.UsingConfig.Host
	connStr += " port=" + config.UsingConfig.Port
	connStr += " user=" + config.UsingConfig.Username
	connStr += " dbname=" + config.UsingConfig.DBname
	connStr += " password=" + config.UsingConfig.Password
	connStr += " sslmode=" + config.UsingConfig.SSLmode

	database, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	database.MustExec(schema)
}

func DisconnectDatabase()  {

	database.Close()
}

func CreateResource(name, url string) {

	command := "INSERT INTO resource (name, url) VALUES ($1, $2)"
	tx := database.MustBegin()
	tx.MustExec(command, name, url)
	tx.Commit()
}

func FindResource(name string) (Resource) {

	command := "SELECT * FROM resource WHERE name=$1"
	resource := Resource{}
	err := database.Get(&resource, command, name)
	if err != nil {
		log.Fatal(err)
	}
	return resource
}

func QueryResources(team string) {

	var teams []ResourceToTeam
	team_id := FindTeam(team).Id
	team_sql := "SELECT * FROM resource_to_team WHERE team_id = $1"
	database.Select(&teams, team_sql, team_id)
	log.Print(teams)
}

func RegisterResource(resource string, team string) {

	exists := ResourceToTeam{}
	resource_id := FindResource(resource).Id
	team_id := FindTeam(team).Id
	check_cmd := "SELECT 1 FROM resource_to_team WHERE resource_id = $1 AND team_id = $2"
	register_cmd := "INSERT INTO resource_to_team (resource_id, team_id) VALUES ($1, $2)"
	
	tx := database.MustBegin()
	row := tx.QueryRow(check_cmd, resource_id, team_id)
	err := row.Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			tx.MustExec(register_cmd, resource_id, team_id)
		} else {
			log.Fatal(err)
		}
	}
	tx.Commit()
}

func CreateTeam(name string) {
	
	command := "INSERT INTO team (name) VALUES ($1)"
	tx := database.MustBegin()
	tx.MustExec(command, name)
	tx.Commit()
}

func FindTeam(name string) (Team) {

	command := "SELECT * FROM team WHERE name=$1"
	team := Team{}
	err := database.Get(&team, command, name)
	if err != nil {
		log.Fatal(err)
	}
	return team
}
