package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"gopkg.in/yaml.v3"
	"log/slog"
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

type Config struct {
	Happy struct {
		Query []struct {
			Name      string `yaml:"name"`
			Statement string `yaml:"statement"`
			Return    string `yaml:"return"`
		} `yaml:"query"`
	} `yaml:"happy"`
}

type Happy struct {
	Name   string `json:"name"`
	Result any    `json:"result"`
}

func validate() Config {
	config := Config{}
	configPath := "happy.yaml"

	//configData := []byte(f)
	configData, err := os.ReadFile(configPath)

	if err != nil {
		slog.Error("couldn't read file", configPath, err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		slog.Error("error unmarshalling yaml", err)
		os.Exit(1)
	}

	for _, query := range config.Happy.Query {
		switch query.Return {
		case "int", "string", "bool":
			slog.Info("config valid")
		default:
			slog.Error("invalid type", query.Return, "for query", query.Name)
		}
	}
	return config
}

func main() {
	config := validate()
	happy := make([]Happy, len(config.Happy.Query))

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		slog.Error("error connecting to db", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	for i := 0; i < len(config.Happy.Query); i++ {
		query := config.Happy.Query[i]
		happy[i].Name = query.Name
		err = conn.QueryRow(context.Background(), query.Statement).Scan(&happy[i].Result)
		if err != nil {
			slog.Error("error performing query", query.Statement, err)
		}
	}

	jsonData, err := json.MarshalIndent(happy, "", " ")
	if err != nil {
		slog.Error("error marshing json", "error", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonData))
	slog.Info("health", "results", jsonData)
	watch()
}
