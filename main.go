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
		Db    string `yaml:"db"`
		Aws   string `yaml:"aws"`
		HealthChecks []struct {
			Name      string `yaml:"name"`
			Statement string `yaml:"statement"`
			Return    string `yaml:"return"`
		} `yaml:"healthchecks"`
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

	for _, hc := range config.Happy.HealthChecks {
		switch hc.Return {
		case "int", "string", "bool":
			slog.Info("config valid")
		default:
			slog.Error("invalid type", hc.Return, "for query", hc.Name)
		}
	}
	return config
}

func main() {
	config := validate()
	happy := make([]Happy, len(config.Happy.HealthChecks))

	db_connection_string := os.Getenv("DB_CONNECTION_STRING")
	if db_connection_string == "" {
		db_connection_string = config.Happy.Db
	}

	conn, err := pgx.Connect(context.Background(), db_connection_string)
	if err != nil {
		slog.Error("error connecting to db", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	for i := 0; i < len(config.Happy.HealthChecks); i++ {
		hc := config.Happy.HealthChecks[i]
		happy[i].Name = hc.Name
		err = conn.QueryRow(context.Background(), hc.Statement).Scan(&happy[i].Result)
		if err != nil {
			slog.Error("error performing query", hc.Statement, err)
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
