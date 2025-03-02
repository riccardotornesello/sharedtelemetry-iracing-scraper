package firestore_structs

type Car struct {
	Name            string `firestore:"name"`
	NameAbbreviated string `firestore:"nameAbbreviated"`
	Brand           string `firestore:"brand"`

	Logo        string `firestore:"logo"`
	SmallImage  string `firestore:"smallImage"`
	SponsorLogo string `firestore:"sponsorLogo"`
}

type CarClass struct {
	Name      string `firestore:"name"`
	ShortName string `firestore:"shortName"`

	Cars []string `firestore:"cars"`
}
