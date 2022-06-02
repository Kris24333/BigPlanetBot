# BigPlanetBot

Telegram bot for getting information about 198 countries of the world.
Get it there [@BigPlanetBot](https://t.me/BigPlanetBot)

## About

This bot is written in GoLang and uses the tgbotapi library https://pkg.go.dev/github.com/go-telegram-bot-api/telegram-bot-api

## Currently available commands
|Command|Description  |
|--|--|
|/by_list|Select a country from the full list of countries.
|/by_alphabet|Select a country from the list of countries sorted alphabetically.
|/by_region|Select a country from the list of countries sorted by region.

Information about each country includes:
* The capital of the country
* IATA code of the capital
* The currency used in the country
* Country area
* Language
* Phone code
* Alpha-2 code
* Alpha-3 code
* Population

It also supports the possibility of getting:
* consular offices in the country
* conditions for entry into a particular country