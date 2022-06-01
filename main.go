package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

// Тип для записи данных с XML-файла
type ValCurs struct {
	Date   string `xml:"Date,attr"`
	Valute []struct {
		CharCode string `xml:"CharCode"`
		Nominal  string `xml:"Nominal"`
		Name     string `xml:"Name"`
		Value    string `xml:"Value"`
	} `xml:"Valute"`
}

func main() {
	// Формирование запроса к сайту, декодирование ответа и запись данных в переменную currencies типа ValCurs
	response, err := http.Get("http://www.cbr.ru/scripts/XML_daily.asp")
	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()

	var currencies ValCurs

	decoder := xml.NewDecoder(response.Body)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return nil, fmt.Errorf("Unknows charset: %s", charset)
		}
	}
	err = decoder.Decode(&currencies)

	var charcodeFromCurrency string
	var charcodeToCurrency string

	var valueFromCurrency float64
	var valueToCurrency float64

	var nomFromCurrency float64
	var nomToCurrency float64

	// Перебор массива в поиске нужной валюты
	for _, currencyData := range currencies.Valute {
		// From NOK
		if currencyData.Name == "Норвежских крон" {
			charcodeFromCurrency = currencyData.CharCode
			nomFromCurrency = stringToFloat(replaceCommaToDot(currencyData.Nominal))
			valueFromCurrency = stringToFloat(replaceCommaToDot(currencyData.Value))
		}
		// To HUF
		if currencyData.Name == "Венгерских форинтов" {
			charcodeToCurrency = currencyData.CharCode
			nomToCurrency = stringToFloat(replaceCommaToDot(currencyData.Nominal))
			valueToCurrency = stringToFloat(replaceCommaToDot(currencyData.Value))
		}
	}

	// Вычисление текущего курса
	var currentExchangeRate float64
	currentExchangeRate = (valueFromCurrency / nomFromCurrency) / (valueToCurrency / nomToCurrency)

	// Формирование исходящего сообщения и вывод его в терминал
	var outMessage = "Курс 1 " + charcodeFromCurrency + " к " + charcodeToCurrency + " на " + currencies.Date + " состаляет: " + fmt.Sprintf("%v", currentExchangeRate)
	fmt.Println(outMessage) // Курс 1 {код валюты} к {код валюты} на {дата} составляет: {текущий курс}
}

// Функция для преобразования данных типа string, записанных из XML-файла, в поля типа float64 для проведения арифметических операций
func stringToFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 32)
	return f
}

// Функция замены запятой в строке на точку. Требуется для корректного преобразования данных типа string, полученных с XML-файла, в данные типа float64 для проведения арифметических операций
func replaceCommaToDot(s string) string {
	out := strings.Replace(s, ",", ".", -1)
	return out
}
