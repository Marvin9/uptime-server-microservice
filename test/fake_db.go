package test

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
)

const (
	CREATE = "CREATE"
	DROP   = "DROP"
)

func FakeDB(operation string) {
	godotenv.Load("../../.env")
	dbname := os.Getenv("DATABASE_NAME")
	var message string
	var command string
	if operation == CREATE {
		message = fmt.Sprintf("Error creating database %v: ", dbname)
		command = "createdb"
	} else {
		message = fmt.Sprintf("Error dropping database %v: ", dbname)
		command = "dropdb"
	}

	cmd := exec.Command(command, "-h", "localhost", "-U", os.Getenv("PSQL_USER"), "-e", dbname)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%v", os.Getenv("PSQL_PASSWORD")))
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		if operation == CREATE {
			log.Printf("Database is there while creating..\n\n")
			FakeDB(DROP)
		}
		log.Printf("%v\n", message)
		log.Println(err)
	}

	// log.Printf("Output: %v", out.String())
}
