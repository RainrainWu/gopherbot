package db

import (

	"log"

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
	resource_id INT NOT NULL,
	team_id INT NOT NULL,
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

	createSQL := `
	INSERT INTO resource (name, url)
	SELECT $1::VARCHAR, $2::VARCHAR
	WHERE NOT EXISTS (
		SELECT * FROM resource
		WHERE resource.name = $1 AND resource.url = $2
	);
	`
	tx := database.MustBegin()
	tx.MustExec(createSQL, name, url)
	tx.Commit()
}

func GetResource(name string) (Resource) {

	getSQL := "SELECT * FROM resource WHERE name=$1"
	resource := Resource{}
	err := database.Get(&resource, getSQL, name)
	if err != nil {
		log.Fatal(err)
	}
	return resource
}

func ListResources() ([]Resource) {

	listSQL := "SELECT * FROM resource"
	var resources []Resource
	err := database.Select(&resources, listSQL)
	if err != nil {
		log.Fatal(err)
	}
	return resources
}

func QueryResources(team string) {

	resourcesSQL := `
	SELECT r.id, r.name, r.url
	FROM   resource r
	JOIN resource_to_team r_t ON r.id = r_t.resource_id
	JOIN team t ON t.id = r_t.team_id
    WHERE t.name = $1
	`
	var resources []Resource
	err := database.Select(&resources, resourcesSQL, team)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(resources)
}

func CreateTeam(name string) {
	
	createSQL := `
	INSERT INTO team (name)
	SELECT $1::VARCHAR
	WHERE NOT EXISTS (
		SELECT * FROM team
		WHERE team.name = $1
	);
	`
	tx := database.MustBegin()
	tx.MustExec(createSQL, name)
	tx.Commit()
}

func GetTeam(name string) (Team) {

	getSQL := "SELECT * FROM team WHERE name=$1"
	team := Team{}
	err := database.Get(&team, getSQL, name)
	if err != nil {
		log.Fatal(err)
	}
	return team
}

func ListTeams() ([]Team) {

	listSQL := "SELECT * FROM team"
	var teams []Team
	err := database.Select(&teams, listSQL)
	if err != nil {
		log.Fatal(err)
	}
	return teams
}

func RegisterResource(resource string, team string) {

	registerSQL := `
	INSERT INTO resource_to_team (resource_id, team_id)
	SELECT $1::INT, $2::INT
	WHERE NOT EXISTS(
		SELECT * FROM resource_to_team
		WHERE resource_id = $1 AND team_id = $2
	);
	`
	resourceId := GetResource(resource).Id
	teamId := GetTeam(team).Id
	tx := database.MustBegin()
	tx.MustExec(registerSQL, resourceId, teamId)
	tx.Commit()
}

func DeregisterResource(resource string, team string) {

	deregisterSQL := `
	DELETE FROM resource_to_team
	WHERE resource_id = $1 AND team_id = $2;
	`
	resourceId := GetResource(resource).Id
	teamId := GetTeam(team).Id
	tx := database.MustBegin()
	tx.MustExec(deregisterSQL, resourceId, teamId)
	tx.Commit()
}
