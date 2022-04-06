package conf

import (
	"database/sql"
	"os"
)

var (
	TELEGRAM_APITOKEN = os.Getenv("TELEGRAM_APITOKEN")
	DB_USERNAME       = os.Getenv("DB_USERNAME")
	DB_PASSWORD       = os.Getenv("DB_PASSWORD")
	DB_ADDRESS        = os.Getenv("DB_ADDRESS")
	DB_NAME           = os.Getenv("DB_NAME")
)

type Keyboard struct {
	Text string
	Data string
}

type CountryMainInfo struct {
	Link           sql.NullString
	Country        sql.NullString
	Capital        sql.NullString
	Currency       sql.NullString
	Area           sql.NullString
	Lang           sql.NullString
	Region         sql.NullString
	PhoneCode      sql.NullString
	LatLon         sql.NullString
	CapitalIATA    sql.NullString
	Alpha2code     sql.NullString
	Alpha3code     sql.NullString
	Population     sql.NullString
	PopulationDate sql.NullString
	In             sql.NullString
	Inside         sql.NullString
	Of             sql.NullString
}

type Embassy struct {
	Of          sql.NullString
	CountryIn   sql.NullString
	EmbassyName sql.NullString
	Address     sql.NullString
	Web         sql.NullString
	Email       sql.NullString
	Head        sql.NullString
	Phone       sql.NullString
	Faks        sql.NullString
	Hours       sql.NullString
	AddedInfo   sql.NullString
}

type EmbassiesInRu struct {
	CountryOf   string
	EmbassyInfo []Embassy
}

type Embassies struct {
	CountryInside string
	EmbassyInfo   []Embassy
}

type Covidrestrictions struct {
	In           string
	Vezd         sql.NullString
	Viza         sql.NullString
	OfficialInfo sql.NullString
	Avia         sql.NullString
	Karantin     sql.NullString
	Usloviya     sql.NullString
	Restrictions sql.NullString
	PCR          sql.NullString
	Date         sql.NullString
}
