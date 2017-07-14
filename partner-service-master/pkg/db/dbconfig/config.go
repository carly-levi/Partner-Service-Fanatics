package dbconfig

import (
	"os"

	"github.com/jackc/pgx"
)

func ExtractConfig() pgx.ConnConfig {
	var config pgx.ConnConfig

	var environment = os.Getenv("environment")//true

	if environment == "" {
		config.Host = "settings.cclyw00l55b3.us-east-1.rds.amazonaws.com"
		config.User = "spam"
		config.Password =  "mapsmaps"
		config.Database = "settings"
	} else {
		config.Host = os.Getenv("partner_service_DB_HOST")
		if config.Host == "" {
			config.Host = "localhost"
		}

		config.User = os.Getenv("partner_service_DB_USER")
		if config.User == "" {
			config.User = os.Getenv("postgres")
		}

		config.Password =  os.Getenv("partner_service_DB_PASSWORD")

		config.Database = os.Getenv("partner_service_DB_DATABASE")
		if config.Database == "" {
			config.Database = "settings"
		}
	}

	return config
}
