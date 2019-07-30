package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tarm/serial"
)

// APIReturn представляет собой набор переменных, требуемых для работы homebridge-thermostat
type APIReturn struct {
	CurrentTemperature string  `json:"currentTemperature"`
	TargetTemperature  float64 `json:"targetTemperature"`
	TargetHeating      int     `json:"targetHeatingCoolingState"`
	CurrentHeating     int     `json:"currentHeatingCoolingState"`
}

var currentHeatingStatus int     // Переменная сообщает нам, происходит ли нагрев в данную минуту
var previousTempString = "00.00" // Сообщает нам последнюю полученную температуру

func sendTemperature(w http.ResponseWriter, r *http.Request) {

	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600} // Устанавливаем настройки Serial порта

	s, err := serial.OpenPort(c)
	if err != nil {
		log.Println("Can't open port")
		log.Fatal(err)
	}

	s.Flush()

	buf := make([]byte, 8) // Получаем данные из Serial порта
	n, err := s.Read(buf)
	if err != nil {
		log.Println("Can't read")
		log.Fatal(err)
	}

	s.Close()

	temperatureString := fmt.Sprintf("%s", buf[:n])                    // Здесь хранится температура в виде строки
	temperatureFloat, err := strconv.ParseFloat(temperatureString, 64) // здесь хранится температура в виде float
	if err != nil {
		log.Println(err)
	}

	if len(temperatureString) == 5 { // Если полученное сообщение содержит 5 байт, то используем её. Нужно исправить это место.

		log.Printf("%q", temperatureString)

		if temperatureFloat > 34.00 { // Меняем переменную currentHeatingStatus указывая на то, что нагрев больше не происходит
			//sendMessageToTelegram("Temperature is now - " + temperatureString + ", too hot!")
			currentHeatingStatus = 0
		}

		if temperatureFloat < 28.00 { // Нагрев происходит
			currentHeatingStatus = 1
		}

		previousTempString = temperatureString // Записываем текущее значение в переменную, на случай не полученного следующего

		response := APIReturn{CurrentTemperature: temperatureString, TargetHeating: 0, TargetTemperature: 34.00, CurrentHeating: currentHeatingStatus} // Формируем ответ используя структуру APIReturn

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	} else { // Если вернулось сообщение отличное от 5 байт, тогда используем предыдущее значение

		response := APIReturn{CurrentTemperature: previousTempString, TargetHeating: 0, TargetTemperature: 34.00, CurrentHeating: currentHeatingStatus} // Формируем ответ используя структуру APIReturn

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	}

}

//func sendMessageToTelegram(message string) {	// Отправляем сообщение в Telegram
//	var query = "https://api.telegram.org/<BOT_TOKEN>/sendMessage?chat_id=<CHATID>&text=" + message
//	http.Get(query)
//}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/status", sendTemperature).Methods("GET")
	fmt.Println(http.ListenAndServe(":8080", r))
}
