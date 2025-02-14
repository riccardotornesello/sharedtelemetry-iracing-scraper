module riccardotornesello.it/sharedtelemetry/iracing/cars_models

go 1.23.2

require (
	ariga.io/atlas-provider-gorm v0.5.0
	gorm.io/gorm v1.25.12
)

replace (
	riccardotornesello.it/sharedtelemetry/iracing/common => ../common
	riccardotornesello.it/sharedtelemetry/iracing/events => ../events
)

require (
	ariga.io/atlas-go-sdk v0.2.3 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.3.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.17 // indirect
	github.com/microsoft/go-mssqldb v1.6.0 // indirect
	golang.org/x/crypto v0.12.0 // indirect
	golang.org/x/text v0.20.0 // indirect
	gorm.io/driver/mysql v1.5.1 // indirect
	gorm.io/driver/postgres v1.5.2 // indirect
	gorm.io/driver/sqlite v1.5.2 // indirect
	gorm.io/driver/sqlserver v1.5.2 // indirect
)
