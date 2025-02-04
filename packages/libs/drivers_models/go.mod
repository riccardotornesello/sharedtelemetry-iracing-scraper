module riccardotornesello.it/sharedtelemetry/iracing/drivers_models

go 1.23.2

require gorm.io/gorm v1.25.12

replace (
	riccardotornesello.it/sharedtelemetry/iracing/common => ../common
	riccardotornesello.it/sharedtelemetry/iracing/events => ../events
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/text v0.20.0 // indirect
)
