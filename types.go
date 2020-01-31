package main

import "net/http"

//	структура для настроек
type SettingsStruct struct {
	ListeningPort int    //	порт сервиса
	ServerDB      string //	IP БД
	PortDB        int    //	порт БД
	IndexDB       string //	название индекса
	TypeDB        string //	название типа
	ProtoDB       string //	название типа
	PartnerURL    string //	URL для отправки сообщения
}

// 	структура для маршрута
type Route struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Method      string           `json:"method"`
	Path        string           `json:"pattern"`
	HandlerFunc http.HandlerFunc `json:"-"`
}

//	структура для сообщения
type Message struct {
	UID           string         `json:"uid"`
	ApID          string         `json:"ap_id"`
	ProbeRequests []ProbeRequest `json:"probe_requests"`
}

type ProbeRequest struct {
	MAC       string `json:"mac"`
	BSSID     string `json:"bssid"`
	TimeStamp string `json:"timestamp"`
	SSID      string `json:"ssid"`
}
