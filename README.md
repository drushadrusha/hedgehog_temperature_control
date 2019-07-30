# Система контроля температуры у ежа
Программа для Arduino, включающая обогреватель стоящий рядом с клеткой моего ежа, когда температура падает ниже 28 градусов и выключающая когда она достигает 34. Arduino отправляет температуру на Raspberry Pi через USB в Homebridge.

Использовались примеры кода с сайта [arduino-diy.com](http://arduino-diy.com/arduino-tsifrovoy-datchik-temperatury-DS18B20)

# Требования
- Go
- [gorilla/mux](https://github.com/gorilla/mux)
- [tarm/serial](https://github.com/tarm/serial)
- Arduino Duemilanove
- [OneWire](https://github.com/PaulStoffregen/OneWire)
- Реле
- DS18B20
- Server (Raspberry Pi)
- Homebridge
- [homebridge-thermostat](https://github.com/PJCzx/homebridge-thermostat)

# Подключение
1. Arduino прошиваем программой из папки `/arduino/`. Меняем переменные в начале кода на нужные пины, если требуется. В программе реле висит на 4 пине, датчик на 7.
2. Arduino подключаем к Raspberry Pi или любому другому компьютеру.
3. В программе `main.go` корректируем адрес Serial-порта (сейчас там стоит `/dev/ttyUSB0`) и требуемый HTTP порт (8080).
4. Устанавливаем все зависимости программы на Go. Компилируем и запускаем программу:
```
go get github.com/gorilla/mux
go get github.com/tarm/serial
go build main.go
./main
```
5. На Raspberry устанавливаем Homebridge и homebridge-thermostat:
 ```
sudo npm install -g --unsafe-perm homebridge
sudo npm install -g homebridge-thermostat
 ```
6. В `config.json` добавляем:
```   
{
    "accessory": "Thermostat",
    "name": "Ёж",
    "apiroute": "http://localhost:8080"
}
```
# TODO
- Управление температурой непосредственно через Homebridge
- Уведомления в Telegram, если температура слишком большая.
- Сделать нормальный парсинг температуры к Serial-порта.
- Можно будет переписать всё на Go через [brutella/hc](https://github.com/brutella/hc) без использования Homebridge.