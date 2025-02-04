module riccardotornesello.it/sharedtelemetry/iracing/events_models

go 1.23.2

require (
	github.com/lib/pq v1.10.9
	gorm.io/gorm v1.25.12
)

replace (
	riccardotornesello.it/sharedtelemetry/iracing/common => ../common
	riccardotornesello.it/sharedtelemetry/iracing/events => ../events
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/text v0.20.0 // indirect
)
