package src

import (
	"country_bot/conf"
	"database/sql"
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	API     *tgbotapi.BotAPI        // API телеграмма
	Updates tgbotapi.UpdatesChannel // Канал обновлений
	//ActiveContactRequests []int64                 // ID чатов, от которых мы ожидаем номер
}

var TB TelegramBot

// Initialize Telegram bot
func (telegramBot *TelegramBot) InitializeTB() {
	bot, err := tgbotapi.NewBotAPI(conf.TELEGRAM_APITOKEN)
	if err != nil {
		log.Panic(err)
	}

	telegramBot.API = bot

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)

	// Tell Telegram we should wait up to 60 seconds on each request for an update
	u.Timeout = 60

	// Start polling Telegram for updates
	updates := bot.GetUpdatesChan(u)

	telegramBot.Updates = updates
}

// Start bot
func (telegramBot *TelegramBot) Start() {
	for update := range telegramBot.Updates {
		if update.Message != nil {
			// Start analize massage
			telegramBot.analyzeUpdate(update)
		} else if update.CallbackQuery != nil {
			// Start analize CallbackQuery
			telegramBot.analyzeCallbackQuery(update)
		}
	}
}

// Analize massage
func (telegramBot *TelegramBot) analyzeUpdate(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	if update.Message.IsCommand() {
		// Analize command massage
		switch update.Message.Command() {
		case "start", "help":
			data := []conf.Keyboard{
				{
					Text: "По алфавиту",
					Data: "by_alphabet",
				}, {
					Text: "Списком",
					Data: "by_list",
				}, {
					Text: "По частям света",
					Data: "by_region",
				},
			}
			msg.Text, msg.ReplyMarkup = generateKeyboard("Привет! Я покажу тебе много полезной информации о 198 странах мира.\n\n*Введи интересующую тебя страну* ниже\n\nИли укажи, каким образом ты хочешь её выбрать:", data, nil)
			msg.ParseMode = "markdown"
		case "by_list":
			countries, err := Connection.GetCountries()
			msg.Text, msg.ReplyMarkup = generateKeyboard("Выбери интересующую тебя страну из списка ниже:", countries, err)
		case "by_alphabet":
			alphabet := GetAlphabet()
			msg.Text, msg.ReplyMarkup = generateKeyboard("Выбери интересующую тебя страну по алфавиту:", alphabet, nil)
		case "by_region":
			regions := GetRegions()
			msg.Text, msg.ReplyMarkup = generateKeyboard("Выбери интересующую тебя страну по частям света:", regions, nil)
		default:
			msg.Text = "К сожалению, я не смог понять чего ты хочешь. Попробуй ещё"
		}
	} else {
		// Analize usual massage
		msg.Text, msg.ReplyMarkup = makeMassage(update.Message.Text)
		msg.ParseMode = "markdown"
	}
	telegramBot.API.Send(msg)
}

// Analize CallbackQuery
func (telegramBot *TelegramBot) analyzeCallbackQuery(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")
	switch update.CallbackQuery.Data {
	case "by_list":
		countries, err := Connection.GetCountries()
		msg.Text, msg.ReplyMarkup = generateKeyboard("Выбери интересующую тебя страну из списка ниже:", countries, err)
	case "by_alphabet":
		alphabet := GetAlphabet()
		msg.Text, msg.ReplyMarkup = generateKeyboard("Выбери интересующую тебя страну по алфавиту:", alphabet, nil)
	case "by_region":
		regions := GetRegions()
		msg.Text, msg.ReplyMarkup = generateKeyboard("Выбери интересующую тебя страну по частям света:", regions, nil)
	default:
		// Split CallbackQuery by ":"
		splitted_command := strings.Split(update.CallbackQuery.Data, ":")
		// Analize CallbackQuery
		switch splitted_command[0] {
		case "l":
			msg.Text, msg.ReplyMarkup = makeMassage(splitted_command[1])
			msg.ParseMode = "markdown"
		case "by_alphabet":
			countries, err := Connection.GetCountriesByAlphabet(splitted_command[1])
			msg.Text, msg.ReplyMarkup = generateKeyboard("Показаны все страны на букву *"+splitted_command[1]+"*\n\nВыбери интересующую тебя страну из списка ниже:", countries, err)
			msg.ParseMode = "markdown"
		case "by_region":
			countries, err := Connection.GetCountriesByRegion(splitted_command[1])
			msg.Text, msg.ReplyMarkup = generateKeyboard("Показаны все страны относящиеся к части света *"+splitted_command[1]+"*\n\nВыбери интересующую тебя страну из списка ниже:", countries, err)
			msg.ParseMode = "markdown"
		case "embassiesinrussia":
			msg.Text = makeMassageEmbassiesInRussia(splitted_command[1])
			msg.ParseMode = "markdown"
		case "embassies":
			msg.Text = makeMassageEmbassies(splitted_command[1])
			msg.ParseMode = "markdown"
		case "covidrestrictions":
			msg.Text = makeMassageCovidrestrictions(splitted_command[1])
			msg.ParseMode = "markdown"
		default:
			msg.Text = "К сожалению, я не смог понять чего ты хочешь. Попробуй ещё"
		}
	}
	telegramBot.API.Send(msg)
}

// Generate keyboard massage
func generateKeyboard(default_text string, data []conf.Keyboard, err error) (string, interface{}) {
	var msg tgbotapi.MessageConfig
	if err != nil {
		msg.Text = "Произошла ошибка! Бот может работать некорректно"
		log.Println(err)
	} else {
		msg.Text = default_text
		keyboard := tgbotapi.InlineKeyboardMarkup{}
		for _, v := range data {
			var row []tgbotapi.InlineKeyboardButton
			btn := tgbotapi.NewInlineKeyboardButtonData(v.Text, v.Data)
			row = append(row, btn)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
		}
		msg.ReplyMarkup = keyboard
	}
	return msg.Text, msg.ReplyMarkup
}

// Make a massage about main country info
func makeMassage(country string) (string, interface{}) {
	res, err := Connection.FindCountry(country)
	switch {
	case err == sql.ErrNoRows:
		return "Я ничего не нашёл по твоему запросу, попробуй ещё раз", nil
	case err != nil:
		return "Произошла ошибка! Бот может работать неправильно!", nil
	default:
		text := fmt.Sprintf("*%s:* справочная информация\n\n", res.Country.String)
		if res.Capital.Valid {
			text = text + fmt.Sprintf("_Столица:_ %s\n", res.Capital.String)
		}
		if res.CapitalIATA.Valid {
			text = text + fmt.Sprintf("_IATA код столицы:_ %s\n", res.CapitalIATA.String)
		}
		if res.Currency.Valid {
			text = text + fmt.Sprintf("_Валюта:_ %s\n", res.Currency.String)
		}
		if res.Area.Valid {
			text = text + fmt.Sprintf("_Площадь:_ %s км2\n", res.Area.String)
		}
		if res.Lang.Valid {
			text = text + fmt.Sprintf("_Язык:_ %s\n", res.Lang.String)
		}
		if res.Region.Valid {
			text = text + fmt.Sprintf("_Часть света:_ %s\n", res.Region.String)
		}
		if res.PhoneCode.Valid {
			text = text + fmt.Sprintf("_Телефонный код:_ +%s\n", res.PhoneCode.String)
		}
		if res.Alpha2code.Valid {
			text = text + fmt.Sprintf("_Альфа-2 код:_ %s\n", res.Alpha2code.String)
		}
		if res.Alpha3code.Valid {
			text = text + fmt.Sprintf("_Альфа-3 код:_ %s\n", res.Alpha3code.String)
		}
		if res.Population.Valid {
			text = text + fmt.Sprintf("_Население:_ %s чел. (на %s)\n", formatNumber(res.Population.String), formatDate(res.PopulationDate.String))
		}
		country_keyboard := Connection.GetKeyboardForCountry(res.Link.String, res.In.String, res.Inside.String, res.Of.String)
		if len(country_keyboard) != 0 {
			_, keyboard := generateKeyboard("", country_keyboard, err)
			return text, keyboard
		} else {
			return text, nil
		}
	}
}

// Make a massage about russian embassies in country
func makeMassageEmbassies(country string) string {
	res, err := Connection.GetEmbassiesInCountry(country)
	switch {
	case err == sql.ErrNoRows:
		return "Я ничего не нашёл по твоему запросу, попробуй ещё раз"
	case err != nil:
		return "Произошла ошибка! Бот может работать неправильно!"
	default:
		text := fmt.Sprintf("Консульские учреждения России *%s*\n", res.CountryInside)
		for _, v := range res.EmbassyInfo {
			text = text + fmt.Sprintf("\n*%s*\n", v.EmbassyName.String)
			if v.AddedInfo.Valid {
				text = text + fmt.Sprintf("%s\n", v.AddedInfo.String)
			}
			if v.Address.Valid {
				text = text + fmt.Sprintf("_Адрес:_ %s\n", v.Address.String)
			}
			if v.Web.Valid {
				text = text + fmt.Sprintf("_Сайт:_ %s\n", v.Web.String)
			}
			if v.Email.Valid {
				text = text + fmt.Sprintf("_Email:_ %s\n", v.Email.String)
			}
			if v.Phone.Valid {
				text = text + fmt.Sprintf("_Телефон:_ %s\n", v.Phone.String)
			}
			if v.Faks.Valid {
				text = text + fmt.Sprintf("_Факс:_ %s\n", v.Faks.String)
			}
			if v.Hours.Valid {
				text = text + fmt.Sprintf("_Часы работы:_ %s\n", v.Hours.String)
			}
		}
		return text
	}
}

// Make a massage about embassies in Russia by country
func makeMassageEmbassiesInRussia(country string) string {
	res, err := Connection.GetEmbassiesInRussiaByCountry(country)
	switch {
	case err == sql.ErrNoRows:
		return "Я ничего не нашёл по твоему запросу, попробуй ещё раз"
	case err != nil:
		return "Произошла ошибка! Бот может работать неправильно!"
	default:
		text := fmt.Sprintf("Консульские учреждения *%s* в России\n", res.CountryOf)
		for _, v := range res.EmbassyInfo {
			var country_in string
			switch v.CountryIn.String {
			case "ekaterinburg":
				country_in = "в Екатеринбурге"
			case "irkutsk":
				country_in = "в Иркутске"
			case "kaliningrad":
				country_in = "в Калининграде"
			case "kazan":
				country_in = "в Казани"
			case "krasnoyarsk":
				country_in = "в Красноярске"
			case "moskva":
				country_in = "в Москве"
			case "nijniy-novgorod":
				country_in = "в Нижнем Новгороде"
			case "novosibirsk":
				country_in = "в Новосибирске"
			case "rostov-na-donu":
				country_in = "в Ростове-на-Дону"
			case "sankt-peterburg":
				country_in = "в Санкт-Петербурге"
			}
			text = text + fmt.Sprintf("\n*%s %s*\n", v.EmbassyName.String, country_in)
			if v.Address.Valid {
				text = text + fmt.Sprintf("_Адрес:_ %s\n", v.Address.String)
			}
			if v.Web.Valid {
				text = text + fmt.Sprintf("_Сайт:_ %s\n", v.Web.String)
			}
			if v.Email.Valid {
				text = text + fmt.Sprintf("_Email:_ %s\n", v.Email.String)
			}
			if v.Head.Valid {
				text = text + fmt.Sprintf("_Руководитель:_ %s\n", v.Head.String)
			}
			if v.Phone.Valid {
				text = text + fmt.Sprintf("_Телефон:_ %s\n", v.Phone.String)
			}
			if v.Hours.Valid {
				text = text + fmt.Sprintf("_Часы работы:_ %s\n", v.Hours.String)
			}
		}
		return text
	}
}

// Make a massage about covidrestrictions in country
func makeMassageCovidrestrictions(country string) string {
	res, err := Connection.GetCovidrestrictions(country)
	switch {
	case err == sql.ErrNoRows:
		return "Я ничего не нашёл по твоему запросу, попробуй ещё раз"
	case err != nil:
		return "Произошла ошибка! Бот может работать неправильно!"
	default:
		text := fmt.Sprintf("Условия для въезда *%s*\n\n", res.In)
		if res.Vezd.Valid {
			text = text + fmt.Sprintf("_Въезд с целью туризма:_ %s\n", removeHTMLFromString(res.Vezd.String))
		}
		if res.Viza.Valid {
			text = text + fmt.Sprintf("_Виза:_ %s\n", removeHTMLFromString(res.Viza.String))
		}
		if res.OfficialInfo.Valid {
			text = text + fmt.Sprintf("_Официальная информация о визовых требованиях:_ %s\n", iterateString(removeHTMLFromString(res.OfficialInfo.String)))
		}
		if res.Avia.Valid {
			text = text + fmt.Sprintf("_Авиасообщение:_ %s\n", removeHTMLFromString(res.Avia.String))
		}
		if res.Karantin.Valid {
			text = text + fmt.Sprintf("_Обязательный карантин по прибытию:_ %s\n", removeHTMLFromString(res.Karantin.String))
		}
		if res.Usloviya.Valid {
			text = text + fmt.Sprintf("%s\n", removeHTMLFromString(res.Usloviya.String))
		}
		if res.Restrictions.Valid {
			text = text + fmt.Sprintf("_Категории граждан, которым разрешен въезд:_ %s\n", removeHTMLFromString(res.Restrictions.String))
		}
		if res.PCR.Valid {
			text = text + fmt.Sprintf("_Условия по ПЦР-тестам и вакцинам:_ %s\n", removeHTMLFromString(res.PCR.String))
		}
		text = text + fmt.Sprintf("\nПоследнее обновление: %s\n", formatDate(res.Date.String))
		text = text + "\nИсточник: [Федеральное агенство по туризму](https://tourism.gov.ru/contents/covid-19/deystvuyushchie-ogranicheniya-po-vezdu-v-inostrannye-gosudarstva/)\n"
		return text
	}
}
