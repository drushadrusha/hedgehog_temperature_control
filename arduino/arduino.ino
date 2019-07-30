#include <OneWire.h>

int Relay = 4;  // реле
OneWire ds(7); // на пине 10 (нужен резистор 4.7 КОм)

void setup(void) {

  Serial.begin(9600);
  pinMode(Relay, OUTPUT); 

}

void loop(void) {

  byte i;

  byte present = 0;

  byte type_s;

  byte data[12];

  byte addr[8];

  float nowTemp, maxTemp, minTemp;

  if ( !ds.search(addr)) {

    ds.reset_search();

    delay(250);

    return;

  }


  if (OneWire::crc8(addr, 7) != addr[7]) {

    return;

  }

  type_s = 0;

  ds.reset();

  ds.select(addr);

  ds.write(0x44);

  delay(1000); 

  present = ds.reset();

  ds.select(addr);

  ds.write(0xBE);

  for ( i = 0; i < 9; i++) { // нам необходимо 9 байт

    data[i] = ds.read();

  }

  // конвертируем данный в фактическую температуру

  // так как результат является 16 битным целым, его надо хранить в

  // переменной с типом данных "int16_t", которая всегда равна 16 битам,

  // даже если мы проводим компиляцию на 32-х битном процессоре

  int16_t raw = (data[1] << 8) | data[0];

  if (type_s) {

    raw = raw << 3; // разрешение 9 бит по умолчанию

  if (data[7] == 0x10) {

    raw = (raw & 0xFFF0) + 12 - data[6];

  }

  } else {

    byte cfg = (data[4] & 0x60);

  // при маленьких значениях, малые биты не определены, давайте их обнулим

    if (cfg == 0x00) raw = raw & ~7; // разрешение 9 бит, 93.75 мс

    else if (cfg == 0x20) raw = raw & ~3; // разрешение 10 бит, 187.5 мс

    else if (cfg == 0x40) raw = raw & ~1; // разрешение 11 бит, 375 мс

//// разрешение по умолчанию равно 12 бит, время преобразования - 750 мс

  }

  nowTemp = (float)raw / 16.0;
  maxTemp = 34.00;
  minTemp = 28.00;

  if (nowTemp >= maxTemp){  // Если температура больше или равна maxTemp, тогда выключаем реле
    digitalWrite(Relay, HIGH); 
  }
  if (nowTemp <= minTemp){  // Если темпратура меньше minTemp, тогда включаем реле
    digitalWrite(Relay, LOW); 
  }
  
  Serial.print(nowTemp);

}
