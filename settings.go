package main

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
)

//	константы инициализации
const (
	//	значение порта по умолчанию
	LISTENING_PORT = 9003        //	порт сервиса
	DB_SERVER      = "elasticsearch" //	IP БД
	DB_PORT        = 9200        //	порт БД
	DB_INDEX       = "skolkovo"  //	название индекса
	DB_TYPE        = "type"      //	название типа
	DB_PROTO       = "http"        //	порт БД
	SETTINGS_FILE_NAME = "skolkovo_test.conf"
	PARTNER_URL = "http://127.0.0.1:9003"
)

var mySettings SettingsStruct	//	настройки

//	загрузки настроек
func (mySettings *SettingsStruct) LoadSettings() {

	//	Читаем содержимое файла настроек
	myFile, err := os.Open(SETTINGS_FILE_NAME)
	if err != nil {
		//	не вышло - пробуем в корень
		log.Error("Не нашел конфиг по относительнмоу пути - пробую абсолютный.")
		myFile, err = os.Open("/opt/skolkovo-test/" + SETTINGS_FILE_NAME)
		if err != nil {

			//	Файла нет - создаем
			log.Warn("Файла нет - создаем.")
			myFile, err := os.OpenFile("./"+SETTINGS_FILE_NAME, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				log.Error("Не получилось создать файл: ", err)
				return
			}
			defer myFile.Close()

			//	Объявляем для примера один параметр
			mySettings.ListeningPort = LISTENING_PORT
			mySettings.ServerDB = DB_SERVER
			mySettings.PortDB = DB_PORT
			mySettings.IndexDB = DB_INDEX
			mySettings.TypeDB = DB_TYPE
			mySettings.ProtoDB = DB_PROTO
			mySettings.PartnerURL = PARTNER_URL
			//	Маршализируем
			buf, _ := json.Marshal(*mySettings)

			//log.Infof("%v", *mySettings)
			//	Копируем структурку в файлик
			myFile.Write(buf)

			return

		}
		defer myFile.Close()
	}
	//	Отложенно закрываем
	defer myFile.Close()

	// Получить размер файла
	stat, err := myFile.Stat()
	if err != nil {
		return
	}

	// Чтение файла
	buf := make([]byte, stat.Size())
	_, err = myFile.Read(buf)
	if err != nil {
		return
	}

	//	Маршалим прочитанное в структуру
	json.Unmarshal(buf, &mySettings)

	//	Выводим сообщение
	log.Infof("Процесс инициализации завершен\nТекущая конфигурация: %v", *mySettings)
}