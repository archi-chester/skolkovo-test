package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// AddHeaders adds all needed headers
func AddHeaders(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Here we adding headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		// http.SetCookie(w, &http.Cookie{Name: "api_key", Value: app.GetAPIKey()})
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		inner.ServeHTTP(w, r)
		//log.Println(w.Header().Get("api_key"))
	})
}

// создание роутера
func CreateRouter(r []Route) *mux.Router {
	// создаем экземпляр роутера
	router := mux.NewRouter().StrictSlash(true)

	log.Info("createRouter")
	var handler http.Handler
	//	создаем роутер

	for _, route := range r {
		// 	ссылка на функцию
		handler = route.HandlerFunc

		handler = AddHeaders(handler)
		router.
			Methods([]string{route.Method}...).
			Path(route.Path).
			Name(route.Name).
			Handler(handler)
	}
	// log.Warnf("Router: %+v", router)
	return router
}


// Serve запуск веб-сервера
func CreateWebServer() {

	// 	создаем ключик
	//secretKey = getNewUUID()

	//	заполняем роуты
	//Routes := append(GeneralRoutes, AddressRoutes...)
	//Routes = append(Routes, AssemblyRoutes...)

	// 	создаем экземпляр роутера
	r := CreateRouter(Routes)

	// headersOk := handlers.AllowedHeaders([]string{"Content-Type", "X-Forwarded-For"})
	// originsOk := handlers.AllowedOrigins([]string{"*"})
	// methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})
	log.Infof("Запуск сервера %s:%d", "", mySettings.ListeningPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", mySettings.ListeningPort), r); err != nil {
		log.Infof("ошибка при запуске сервера: ", err)
	}

}