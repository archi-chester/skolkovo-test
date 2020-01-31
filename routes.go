package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// 	массив роутинга
var Routes = []Route{
	//	сами функции
	Route{
		"getMessage",
		"Получаем сообщение",
		"POST",
		"/messages/get",
		getMessage,
	},
	Route{
		"interceptMessage",
		"Перехватываем сообщение для теста",
		"POST",
		"/messages/intercept",
		interceptMessage,
	},
	//	STATUS OK
	Route{
		"statusOK",
		"Создаем объект",
		"OPTIONS",
		"/messages/intercept",
		statusOK,
	},
	Route{
		"statusOK",
		"Создаем объект",
		"OPTIONS",
		"/messages/get",
		statusOK,
	},
}

//	обработчик получения данных
func getMessage(w http.ResponseWriter, r *http.Request) {
	log.Infof("getMessage, mySettings: %+v", mySettings)
	//	переменные
	var message Message

	//	декодирование переданной структуры из запроса
	//	проверяем, что не вернули ошибку
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		//	вернули ошибку - возвращаем 422
		log.Errorf("Ошибка при распаковке объекта - неверный формат запроса: %+v", err)
		http.Error(w, fmt.Sprintf("{\"status\":\"error\",\"error\":\"ООшибка при распаковке объекта - неверный формат запроса: : %+v\"}", err), http.StatusUnprocessableEntity)
		return
	} else {
		//	формируем uuid для сообщения
		message.UID = getNewUUID()

		//	пробуем сохранить сообщение
		if err := saveMessage(message); err != nil {
			log.Errorf("не удалось сохранить сообщение в базу: %+v", err)
		}
	}

	//	если массив не пустой
	if len(message.ProbeRequests) > 0 {
		for index, probeRequest := range message.ProbeRequests {

			//	если BSSID пустой - заменяем на дефолт
			if probeRequest.BSSID == "" {
				message.ProbeRequests[index].BSSID = "FF-FF-FF-FF-FF-FF"
			}

			//	если SSID пустой - заменяем на дефолт
			if probeRequest.SSID == "" {
				message.ProbeRequests[index].SSID = "Unknown"
			}
		}

		//	отправляем сообщение
		if err := sendMessage(message); err != nil {
			log.Errorf("не удалось отправить сообщение партнерам: %+v", err)
			http.Error(w, fmt.Sprintf("{\"status\":\"error\",\"error\":\"не удалось отправить сообщение партнерам : %+v\"}", err), http.StatusNotFound)
			return
		} else {
			//  все нормально - возвращаем 200
			w.WriteHeader(200)
			_ = json.NewEncoder(w).Encode(message)
			return
		}
	}
}

//	сохранение сообщение
func interceptMessage(w http.ResponseWriter, r *http.Request) {
	//	переменные
	var message Message

	//	декодирование переданной структуры из запроса
	//	проверяем, что не вернули ошибку
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		//	вернули ошибку - возвращаем 422
		log.Errorf("Ошибка при распаковке объекта - неверный формат запроса: %+v", err)
		http.Error(w, fmt.Sprintf("{\"status\":\"error\",\"error\":\"ООшибка при распаковке объекта - неверный формат запроса: : %+v\"}", err), http.StatusUnprocessableEntity)
		return
	} else {
		log.Infof("Получено перехватчиком: %+v", message)
		//  все нормально - возвращаем 200
		w.WriteHeader(200)
		//	возвращаем полученные объекты
		_ = json.NewEncoder(w).Encode(message)
	}
}

//	сохранение сообщение
func saveMessage(message Message) error {
	//	переменные
	var reqBodyBytes []byte
	reqBodyBytes, _ = json.Marshal(&message)

	//	на всякий случай отрубаем проверку сертифката, если там https
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	//	объявляем клиента для запроса
	client := &http.Client{Transport: tr}

	//	дергаем запрос
	req, _ := http.NewRequest("POST", mySettings.ProtoDB + "://" + mySettings.ServerDB + ":" +
		fmt.Sprintf("%d", mySettings.PortDB) + "/" + mySettings.IndexDB + "/" + mySettings.TypeDB + "/" + message.UID,
		bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")


	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Ошибка в сохранении в БД. resp: %+v err: %+v", resp, err)
		return err
	} else {
		log.Warnf(mySettings.ProtoDB + "://" + mySettings.ServerDB + "/" +
			mySettings.IndexDB + "/" + mySettings.TypeDB + "/" + message.UID)
		log.Warnf("Ошибка в сохранении в БД. resp: %+v", resp)
		log.Infof("saveMessage: %+v", message)
		return nil
	}
}

//	отправляем сообщение
func sendMessage(message Message) error {
	//	переменные
	var reqBodyBytes []byte
	reqBodyBytes, _ = json.Marshal(&message)

	//	на всякий случай отрубаем проверку сертифката, если там https
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	//	объявляем клиента для запроса
	client := &http.Client{Transport: tr}

	//	дергаем запрос
	log.Infof("PartnerURL: %+v", mySettings.PartnerURL)
	req, _ := http.NewRequest("POST", mySettings.PartnerURL,
		bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")


	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Ошибка при отправке сообщения партнерам. %+v err: %+v", resp, err)
		return err
	} else {
		log.Infof("sendMessage: %+v", message)
		return nil
	}
}

// 	statusOK - статус ОК
func statusOK(w http.ResponseWriter, r *http.Request) {
	// Here we adding headers
	// w.Header().Set("Access-Control-Allow-Headers", "*")
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
	// // http.SetCookie(w, &http.Cookie{Name: "api_key", Value: app.GetAPIKey()})
	// w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// inner.ServeHTTP(w, r)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode("{'status':'OK'}")
	return
}