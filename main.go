package main

import (
	"github.com/twinj/uuid"

	log "github.com/sirupsen/logrus"
)

//	skolkovo_test
//	По для получения данных AP по REST, сохранение их в БД и отправки на сторонее API

// 	точка входа
func main() {
	log.SetLevel(log.InfoLevel)
	//var mySettings *SettingsStruct

	//	объявили настройки

	//	проинициализировали
	mySettings.LoadSettings()



	// Включаем роутер
	CreateWebServer()
}

// GetNewUUID - функция отдает новый UUID
func getNewUUID() string {
	return uuid.NewV4().String()
}