package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"gopkg.in/yaml.v3"
	"os"
)

var f = `
queries: 
  - "select count (*) from pg_database" 
`

type Config struct {
	Queries []string `yaml:"queries"`
}

type Resp struct {
	version string
	count   int
}

func main() {

	config := Config{}
	resp := Resp{}

	config_path := "config.yaml"
	


	err := yaml.Unmarshal([]byte(f), &config)
	if err != nil {
		os.Exit(2)
	}


	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	fmt.Println(config)
	for _, query := range config.Queries {
		fmt.Println(query)
		err = conn.QueryRow(context.Background(), query).Scan(&resp.count)
	}

	if err != nil {
		os.Exit(1)
	}

	fmt.Println("num dbs is:", resp.count)
}
