module riccardotornesello.it/sharedtelemetry/iracing/drivers_downloader

go 1.23.2

require (
	github.com/joho/godotenv v1.5.1
	gorm.io/gorm v1.25.12
	riccardotornesello.it/sharedtelemetry/iracing/cloudrun_utils v0.0.0-00010101000000-000000000000
	riccardotornesello.it/sharedtelemetry/iracing/drivers_models v0.0.0-00010101000000-000000000000
	riccardotornesello.it/sharedtelemetry/iracing/gorm_utils v0.0.0-00010101000000-000000000000
	riccardotornesello.it/sharedtelemetry/iracing/irapi v0.0.0-00010101000000-000000000000
)

replace (
	riccardotornesello.it/sharedtelemetry/iracing/cloudrun_utils => ../../libs/cloudrun_utils
	riccardotornesello.it/sharedtelemetry/iracing/drivers_models => ../../libs/drivers_models
	riccardotornesello.it/sharedtelemetry/iracing/gorm_utils => ../../libs/gorm_utils
	riccardotornesello.it/sharedtelemetry/iracing/irapi => ../../libs/irapi
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d80cb9 // indirect
	github.com/jackc/pgx/v5 v5.5.5 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/crypto v0.29.0 // indirect
	golang.org/x/sync v0.9.0 // indirect
	golang.org/x/text v0.20.0 // indirect
	gorm.io/driver/postgres v1.5.11 // indirect
)
