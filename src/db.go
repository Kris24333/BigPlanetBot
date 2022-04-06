package src

import (
	"country_bot/conf"
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type DatabaseConnection struct {
	DB *sql.DB
}

var Connection DatabaseConnection

// Initialize connection to the database
func (connection *DatabaseConnection) InitializeDB() {
	db, err := sql.Open("mysql", conf.DB_USERNAME+":"+conf.DB_PASSWORD+"@tcp("+conf.DB_ADDRESS+")/"+conf.DB_NAME)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	connection.DB = db
}

// Get full list of countries
func (connection *DatabaseConnection) GetCountries() ([]conf.Keyboard, error) {
	scan_query := "SELECT `country_name` FROM `categorycountries` ORDER BY `country_name`"
	countries, error := connection.makeListOfCountriesBySql(scan_query, "l")
	return countries, error
}

// Get list of countries by alphabet
func (connection *DatabaseConnection) GetCountriesByAlphabet(letter string) ([]conf.Keyboard, error) {
	scan_query := "SELECT `country_name` FROM `categorycountries` WHERE `country_name` LIKE '" + letter + "%' ORDER BY `country_name`"
	countries, error := connection.makeListOfCountriesBySql(scan_query, "l")
	return countries, error
}

// Get list of countries by region
func (connection *DatabaseConnection) GetCountriesByRegion(region string) ([]conf.Keyboard, error) {
	scan_query := "SELECT `country_name` FROM `categorycountries` WHERE `region` LIKE '%" + region + "%' ORDER BY `country_name`"
	countries, error := connection.makeListOfCountriesBySql(scan_query, "l")
	return countries, error
}

// Make list of countries by sql
func (connection *DatabaseConnection) makeListOfCountriesBySql(scan_query string, command string) ([]conf.Keyboard, error) {
	var countries []conf.Keyboard
	scan, scan_error := connection.DB.Query(scan_query)
	if scan_error != nil {
		log.Println(scan_error)
		return countries, scan_error
	}
	defer scan.Close()
	for scan.Next() {
		var country conf.Keyboard
		if err := scan.Scan(&country.Text); err != nil {
			log.Println(err)
			return nil, err
		}
		country.Data = command + ":" + country.Text
		countries = append(countries, country)
	}
	return countries, nil
}

// Get alphabet list
func GetAlphabet() []conf.Keyboard {
	alphabet := []string{"А", "Б", "В", "Г", "Д", "Е", "З", "И", "Й", "К", "Л", "М", "Н", "О", "П", "Р", "С", "Т", "У", "Ф", "Х", "Ц", "Ч", "Ш", "Э", "Ю", "Я"}
	datas := getList(alphabet, "by_alphabet")
	return datas
}

// Get regions list
func GetRegions() []conf.Keyboard {
	regions := []string{"Австралия и Океания", "Азия", "Америка", "Африка", "Европа"}
	datas := getList(regions, "by_region")
	return datas
}

// Get list
func getList(list []string, command string) []conf.Keyboard {
	var datas []conf.Keyboard
	for _, v := range list {
		data := conf.Keyboard{
			Text: v,
			Data: command + ":" + v,
		}
		datas = append(datas, data)
	}
	return datas
}

// Get info about country
func (connection *DatabaseConnection) FindCountry(country string) (conf.CountryMainInfo, error) {
	var info conf.CountryMainInfo
	scan_error := connection.DB.QueryRow("SELECT country_name, capital, money, area, language, region, phone_code, latlon, capital_iata, alpha2code, alpha3code, population, population_date, inin, inside, of, country_link FROM categorycountries WHERE country_name LIKE ?", "%"+country+"%").Scan(&info.Country, &info.Capital, &info.Currency, &info.Area, &info.Lang, &info.Region, &info.PhoneCode, &info.LatLon, &info.CapitalIATA, &info.Alpha2code, &info.Alpha3code, &info.Population, &info.PopulationDate, &info.In, &info.Inside, &info.Of, &info.Link)
	switch {
	case scan_error == sql.ErrNoRows:
		log.Println(scan_error)
		return info, scan_error
	case scan_error != nil:
		log.Println(scan_error)
		return info, scan_error
	default:
		return info, nil
	}
}

// Get covidrestrictions in country
func (connection *DatabaseConnection) GetCovidrestrictions(country string) (conf.Covidrestrictions, error) {
	var infos conf.Covidrestrictions
	scan_error := connection.DB.QueryRow("SELECT categorycountries.inin, covidrestrictions.vezd_s_celyu_turizma, covidrestrictions.viza, covidrestrictions.oficialnaya_informaciya_o_vizovyh_trebovaniyah, covidrestrictions.aviasoobshchenie, covidrestrictions.obyazatelnyy_karantin_po_pribytiyu, covidrestrictions.usloviya_vezda, covidrestrictions.kategorii_grajdan_kotorym_razreshen_vezd, covidrestrictions.usloviya_po_pcr_testam_ili_vakcinam FROM categorycountries, covidrestrictions WHERE categorycountries.country_link=covidrestrictions.country AND covidrestrictions.country=?", country).Scan(&infos.In, &infos.Vezd, &infos.Viza, &infos.OfficialInfo, &infos.Avia, &infos.Karantin, &infos.Usloviya, &infos.Restrictions, &infos.PCR)
	if scan_error != nil {
		log.Println(scan_error)
		return infos, scan_error
	} else {
		scan_error2 := connection.DB.QueryRow("SELECT upd FROM admin WHERE parsed_table='covidrestrictions'").Scan(&infos.Date.String)
		if scan_error2 != nil {
			log.Println(scan_error2)
			return infos, scan_error2
		} else {
			return infos, nil
		}
	}
}

// Get russian embassies in country
func (connection *DatabaseConnection) GetEmbassiesInCountry(country string) (conf.Embassies, error) {
	var infos conf.Embassies
	scan_error := connection.DB.QueryRow("SELECT inside FROM categorycountries WHERE country_link=?", country).Scan(&infos.CountryInside)
	if scan_error != nil {
		log.Println(scan_error)
		return infos, scan_error
	}

	var embassies []conf.Embassy
	scan_query := "SELECT `embassy_name`, `adress`, `web`, `email`, `phone`, `hours`, `faks`, `added_info` FROM `embassies` WHERE `country` = '" + country + "'"
	scan, scan_error := connection.DB.Query(scan_query)
	if scan_error != nil {
		log.Println(scan_error)
		return infos, scan_error
	}
	defer scan.Close()
	for scan.Next() {
		var embassy conf.Embassy
		if err := scan.Scan(&embassy.EmbassyName, &embassy.Address, &embassy.Web, &embassy.Email, &embassy.Phone, &embassy.Hours, &embassy.Faks, &embassy.AddedInfo); err != nil {
			log.Println(err)
			return infos, err
		}
		embassies = append(embassies, embassy)
	}
	infos.EmbassyInfo = escapeCharacters(embassies)
	return infos, nil
}

// Get embassies in Russia by country
func (connection *DatabaseConnection) GetEmbassiesInRussiaByCountry(country string) (conf.EmbassiesInRu, error) {
	var infos conf.EmbassiesInRu
	scan_error := connection.DB.QueryRow("SELECT of FROM categorycountries WHERE country_link=?", country).Scan(&infos.CountryOf)
	if scan_error != nil {
		log.Println(scan_error)
		return infos, scan_error
	}

	var embassies []conf.Embassy
	scan_query := "SELECT `country_in`, `embassy_name`, `adress`, `web`, `email`, `head`, `phone`, `hours` FROM `embassiesinrussia` WHERE `country` = '" + country + "' ORDER BY `country_in`"
	scan, scan_error := connection.DB.Query(scan_query)
	if scan_error != nil {
		log.Println(scan_error)
		return infos, scan_error
	}
	defer scan.Close()
	for scan.Next() {
		var embassy conf.Embassy
		if err := scan.Scan(&embassy.CountryIn, &embassy.EmbassyName, &embassy.Address, &embassy.Web, &embassy.Email, &embassy.Head, &embassy.Phone, &embassy.Hours); err != nil {
			log.Println(err)
			return infos, err
		}
		embassies = append(embassies, embassy)
	}
	infos.EmbassyInfo = escapeCharacters(embassies)
	return infos, nil
}

// Keyboard for country
func (connection *DatabaseConnection) GetKeyboardForCountry(link string, in string, inside string, of string) []conf.Keyboard {
	var datas []conf.Keyboard
	var data conf.Keyboard
	var id int
	scan_error := connection.DB.QueryRow("SELECT id FROM covidrestrictions WHERE country=?", link).Scan(&id)
	if scan_error == nil {
		data = conf.Keyboard{
			Text: "Въезд " + in + " в условиях коронавирусных ограничений",
			Data: "covidrestrictions:" + link,
		}
		datas = append(datas, data)
	}
	scan_error1 := connection.DB.QueryRow("SELECT id FROM embassies WHERE country=?", link).Scan(&id)
	if scan_error1 == nil {
		data = conf.Keyboard{
			Text: "Консульские учреждения России " + inside,
			Data: "embassies:" + link,
		}
		datas = append(datas, data)
	}
	scan_error2 := connection.DB.QueryRow("SELECT id FROM embassiesinrussia WHERE country=?", link).Scan(&id)
	if scan_error2 == nil {
		data = conf.Keyboard{
			Text: "Консульские учреждения " + of + " в России",
			Data: "embassiesinrussia:" + link,
		}
		datas = append(datas, data)
	}
	return datas
}
