package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"gopkg.in/yaml.v3"
	"os"
)

/*
var f = `
happy:
  query: 
    - name: postgres_version
      statement: "select version()"
      return: "string"
    - name: database_count
      statement: "select count (*) from pg_database" 
      return: "int"
`
*/

type Happy struct {
	Queries []string `yaml:"queries"`
}

type Resp struct {
	version string
	count   int
}

func main() {

	happy := Happy{}
	resp := Resp{}

	config_path := "happy.yaml"
	
	config, err := os.ReadFile(config_path)
	
	if err != nil {
		os.Exit(1)
	}

	err = yaml.Unmarshal(config, &happy)
	if err != nil {
		os.Exit(1)
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	for _, query := range happy.Queries {
		err = conn.QueryRow(context.Background(), query).Scan(&resp.count)
	}

	if err != nil {
		os.Exit(1)
	}

	fmt.Println("num dbs is:", resp.count)
}
