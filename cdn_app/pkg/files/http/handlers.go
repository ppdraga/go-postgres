package http

import (
	"cdn_app/httpapi"
	fservice "cdn_app/pkg/files/service"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)




func APIHandler(filesCtrl *fservice.Controller) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name, ok := vars["name"]
		if !ok {
			httpapi.ResponseBadRequest("Couldn't parse name param", w)
			return
		}

		uri := strings.Split(r.RequestURI, "/")
		fmt.Println(uri)
		switch uri[3] {
		case "user":
			item, err := filesCtrl.GetUserFiles(name)
			if err != nil {
				errMsg := fmt.Sprintf("Error: %v", err)
				httpapi.ResponseBadRequest(errMsg, w)
				return
			}
			json.NewEncoder(w).Encode(*item)
			return
		case "server":
			item, err := filesCtrl.GetServerFiles(name)
			if err != nil {
				errMsg := fmt.Sprintf("Error: %v", err)
				httpapi.ResponseBadRequest(errMsg, w)
				return
			}
			json.NewEncoder(w).Encode(*item)
			return
		case "area":
			item, err := filesCtrl.GetAreaFiles(name)
			if err != nil {
				errMsg := fmt.Sprintf("Error: %v", err)
				httpapi.ResponseBadRequest(errMsg, w)
				return
			}
			json.NewEncoder(w).Encode(*item)
			return
		}
	}



}
