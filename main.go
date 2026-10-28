package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
)

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
	version string
	count   int
	str     string
}

func validate() Config {
	config := Config{}
	/*
		config_path := "happy.yaml"

		config_data, err := os.ReadFile(config_path)

		if err != nil {
			os.Exit(1)
		}
	*/
	config_data := []byte(f)

	err := yaml.Unmarshal(config_data, &config)
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
	happy := Happy{}

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		slog.Error("error connecting to db", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	for _, query := range config.Happy.Query {
		switch query.Return {
		case "string":
			err = conn.QueryRow(context.Background(), query.Statement).Scan(&happy.version)
			if err != nil {
				slog.Error("error performing query", query.Statement, err)
			}
		case "int":
			err = conn.QueryRow(context.Background(), query.Statement).Scan(&happy.count)
			if err != nil {
				slog.Error("error performing query", query.Statement, err)
			}
		default:
			slog.Error("unrecoverable error")
			os.Exit(1)
		}
	}

	slog.Info("result", "dbs", happy.count)
}
