package firestore_structs

type Driver struct {
	Name     string `firestore:"name"`
	Location string `firestore:"location"`

	Stats DriverStats `firestore:"stats"`
}

type DriverStats struct {
	DirtOval   *DriverStatsDetails `firestore:"dirtOval"`
	DirtRoad   *DriverStatsDetails `firestore:"dirtRoad"`
	FormulaCar *DriverStatsDetails `firestore:"formulaCar"`
	Oval       *DriverStatsDetails `firestore:"oval"`
	Road       *DriverStatsDetails `firestore:"road"`
	SportsCar  *DriverStatsDetails `firestore:"sportsCar"`
}

type DriverStatsDetails struct {
	License string `firestore:"not null"`
	IRating int    `firestore:"not null"`
}
