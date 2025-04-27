package main

import (
	"ChiragKr04/go-backend/config"
	"ChiragKr04/go-backend/db"
	"log"
	"os"
	"strconv"

	mysqlCfg "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// parseInt converts a string to an int and handles errors
func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("Error parsing version '%s': %v", s, err)
	}
	return i
}

func main() {
	appConfig := mysqlCfg.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}
	db, err := db.MySQLStorage(appConfig)
	if err != nil {
		log.Fatal(err)
	}

	// init driver
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		// path for migrations for macos
		"file://cmd/migrate/migrations",
		"mysql",
		driver,
	)

	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) < 2 {
		log.Fatal("No command provided. Use 'up', 'down', or 'force <version>'")
	}

	cmd := os.Args[1]
	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	} else if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	} else if cmd == "force" {
		if len(os.Args) < 3 {
			log.Fatal("No version provided for force command. Use 'force <version>'")
		}
		// Get the version from the third argument
		version := os.Args[2]
		if err := m.Force(parseInt(version)); err != nil {
			log.Fatal(err)
		}
	} else if cmd == "status" {
		// Get current migration version
		version, dirty, err := m.Version()
		if err != nil {
			if err == migrate.ErrNilVersion {
				log.Println("No migrations have been applied yet")
				return
			}
			log.Fatal(err)
		}
		
		log.Printf("Current migration version: %d, Dirty: %t\n", version, dirty)
		if dirty {
			log.Println("WARNING: The database is in a dirty state. This means that the last migration failed.")
			log.Println("You may need to fix the database manually and then use 'force' to set the version.")
		}
	} else {
		log.Fatal("Invalid command. Use 'up', 'down', 'status', or 'force <version>'")
	}

}
