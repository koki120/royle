package cmd

import (
	"database/sql"
	"log"
	"os"

	"github.com/digeon-inc/royle/filter/consumer"
	"github.com/digeon-inc/royle/filter/producer"
	"github.com/digeon-inc/royle/filter/transformer"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
)

var (
	title      string
	dbUser     string
	dbPassword string
	dbHost     string
	dbPort     string
	database   string
)

var rootCmd = &cobra.Command{
	Use:   "royle",
	Short: "Generates documentation for the MySQL tables.",
	Long:  "This is a command-line application written in Go that connects to a MySQL database, extracts table information, and generates a file documenting the database tables.",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := sql.Open("mysql", INFORMATION_SCHEMA_DSN())
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		source, err := producer.FetchColumnMetadata(db, DatabaseName())
		if err != nil {
			log.Fatal(err)
		}

		tables := transformer.ConvertColumnMetadataToTableMetaData(source)

		if err = consumer.ExportToMarkdown(os.Stdout, Title(), tables); err != nil {
			log.Fatal(err)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&title, "title", "t", "ROYLE", "document title")

	rootCmd.Flags().StringVarP(&dbUser, "user", "u", "", "mysql user")
	if err := rootCmd.MarkFlagRequired("user"); err != nil {
		log.Fatal(err)
	}

	rootCmd.Flags().StringVarP(&dbPassword, "password", "p", "", "mysql password")
	if err := rootCmd.MarkFlagRequired("password"); err != nil {
		log.Fatal(err)
	}

	rootCmd.Flags().StringVarP(&dbHost, "host", "s", "", "mysql host")
	if err := rootCmd.MarkFlagRequired("host"); err != nil {
		log.Fatal(err)
	}

	rootCmd.Flags().StringVarP(&dbPort, "port", "r", "", "mysql port")
	if err := rootCmd.MarkFlagRequired("port"); err != nil {
		log.Fatal(err)
	}

	rootCmd.Flags().StringVarP(&database, "database", "d", "", "mysql database name")
	if err := rootCmd.MarkFlagRequired("database"); err != nil {
		log.Fatal(err)
	}
}
