package main

import (
	//"strings"
	"log"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/pkg/errors"
	"github.com/thinkerou/favicon"

	"bruharmy/models"
)

var (
	tomlConf = &models.Config{}
	configPath = "config.conf"
)

func main() {
	// setup database
	// if _, err := os.Stat("./database.db"); errors.Is(err, os.ErrNotExist) {
	// 	create_database()
	// }
	models.ReadConfig(tomlConf, configPath)
	err := models.CheckConfig(tomlConf)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "illegal config"))
	}

	// setup router
	router := gin.Default()
	router.Use(favicon.New("./assets/favicon.ico"))
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))
	router.Static("/assets", "./assets")
	// router.LoadHTMLGlob("templates/*.html")
	router.MaxMultipartMemory = 8 << 20 // 8Mib

	// router.Use(sessions.Sessions("session", cookie.NewStore(globals.Secret)))

	public := router.Group("/")
	PublicRoutes(public)

	private := router.Group("/")
	private.Use(JwtAuthRequired)
	PrivateRoutes(private)

	LoadPortGroups()

	if tomlConf.Https {
		log.Fatal(router.RunTLS(":" + fmt.Sprint(tomlConf.Port), tomlConf.Cert, tomlConf.Key))
	} else {
		log.Fatal(router.Run(":" + fmt.Sprint(tomlConf.Port)))
	}
}

// func create_database() {
// 	os.Create("./database.db")
// 	stmt := ``
// 	db, err := sql.Open("sqlite3", "./database.db")
// 	if err != nil {
// 		log.Println("Failed to open database!")
// 		os.Exit(1)
// 	}
// 	// sql
// 	stmt = `
// 	PRAGMA foreign_keys = ON;
// 	CREATE TABLE users (user_id INTEGER PRIMARY KEY AUTOINCREMENT, username VARCHAR(64) NOT NULL, password VARCHAR(64) NOT NULL, color VARCHAR(7) DEFAULT "#777777");
// 	INSERT INTO users (user_id, username, password) VALUES (0, "No Assignee", "");
// 	CREATE TABLE boxes (ip VARCHAR(64) PRIMARY KEY, hostname VARCHAR(64) NULL, codename VARCHAR(64) NULL, assignee INTEGER DEFAULT 0, usershells INTEGER DEFAULT 0, rootshells INTEGER DEFAULT 0, FOREIGN KEY (assignee) REFERENCES users(user_id));
// 	CREATE TABLE ports (port_id INTEGER PRIMARY KEY AUTOINCREMENT, port_number INTEGER NOT NULL, protocol VARCHAR(64) NOT NULL, service_name VARCHAR(64) NULL, service_details VARCHAR(8192), box_ip VARCHAR(64) NOT NULL, FOREIGN KEY (box_ip) REFERENCES boxes(ip));
// 	CREATE TABLE credentials (credential_id INTEGER PRIMARY KEY AUTOINCREMENT, ip VARCHAR(64) NULL, hostname VARCHAR(64) NULL, port INTEGER NOT NULL, service VARCHAR(64) NOT NULL, username VARCHAR(64) NOT NULL, password VARCHAR(64) NOT NULL);
// 	`
// 	_, err = db.Exec(stmt)
// 	if err != nil {
// 		log.Println("Failed to create tables!", err)
// 		os.Exit(1)
// 	}
// 	defer db.Close()
	
// }