package src

import (
	"country_bot/conf"
	"regexp"
	"strconv"
	"time"
)

// Format date
func formatDate(date string) string {
	var date_layout = "2006-01-02T15:04:05-07:00"
	t, _ := time.Parse(date_layout, date)
	var rus_months string
	switch t.Month().String() {
	case "January":
		rus_months = "января"
	case "February":
		rus_months = "февраля"
	case "March":
		rus_months = "марта"
	case "April":
		rus_months = "апреля"
	case "May":
		rus_months = "мая"
	case "June":
		rus_months = "июня"
	case "July":
		rus_months = "июля"
	case "August":
		rus_months = "августа"
	case "September":
		rus_months = "сентября"
	case "October":
		rus_months = "октября"
	case "December":
		rus_months = "декабря"
	case "November":
		rus_months = "ноября"
	}
	new_date := strconv.Itoa(t.Day()) + " " + rus_months + " " + strconv.Itoa(t.Year()) + " года"
	return new_date
}

// Format a number with grouped thousands
func formatNumber(string_number string) string {
	numOfDigits := len(string_number)
	numOfSpaces := (numOfDigits - 1) / 3

	out := make([]byte, len(string_number)+numOfSpaces)

	for i, j, k := len(string_number)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = string_number[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = ' '
		}
	}
}

// Escape characters '_', '*', '`', '[' for Markdown style from []conf.Embassy
func escapeCharacters(res []conf.Embassy) []conf.Embassy {
	new_res := res
	for k, v := range res {
		new_res[k].Web.String = iterateString(v.Web.String)
		new_res[k].Address.String = iterateString(v.Address.String)
		new_res[k].Email.String = iterateString(v.Email.String)
		new_res[k].Phone.String = iterateString(v.Phone.String)
		new_res[k].Faks.String = iterateString(v.Faks.String)
		new_res[k].Hours.String = iterateString(v.Hours.String)
		new_res[k].AddedInfo.String = iterateString(v.AddedInfo.String)
	}
	return new_res
}

// Escape characters '_', '*', '`', '[' for Markdown style from string
func iterateString(text string) string {
	new_string := text
	escape_characters := map[string]string{
		"_":   "\\_",
		"\\*": "\\*",
		"`":   "\\`",
		"\\[": "\\[",
	}
	var tr *regexp.Regexp
	for from, to := range escape_characters {
		tr = regexp.MustCompile(from)
		new_string = tr.ReplaceAllString(new_string, to)
	}
	return new_string
}

// Remove HTML from tring and add \n where needed
func removeHTMLFromString(text string) string {
	new_string := text
	escape_characters := map[string]string{
		"</p><p>": "\n",
		"Для въезда понадобится:": "_Для въезда понадобится:_",
		"</p>": "",
		"<p>":  "",
	}
	var tr *regexp.Regexp
	for from, to := range escape_characters {
		tr = regexp.MustCompile(from)
		new_string = tr.ReplaceAllString(new_string, to)
	}
	re, _ := regexp.Compile("<a href=\"(.+)\">")
	res := re.FindAllStringSubmatch(new_string, -1)
	if len(res) != 0 {
		return res[0][1]
	} else {
		return new_string
	}
}
