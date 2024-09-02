package log

/*
import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

type Log struct {
	mu      sync.Mutex
	records []Plant
}

var almanac Log

type Plant struct {
	PlantName string `json:"plantName"`
	Cost      uint64 `json:"cost"`
	Recharge  string `json:"recharge"`
	offset    uint64 `"json:offset"`
}

func main() {

	peashooter := Plant{
		PlantName: "peashooter",
		Cost:      100,
		Recharge:  "normal",
		offset:    1,
	}
	almanac.records = append(almanac.records, peashooter)
	handler := mux.NewRouter()
	subHandler := handler.PathPrefix("/almanac").Subrouter()
	subHandler.HandleFunc("/{plant}", getPlant).Methods(http.MethodGet)
	subHandler.HandleFunc("/new", postPlant).Methods(http.MethodPost)

	http.ListenAndServe(":80", handler)
}

func getPlant(w http.ResponseWriter, r *http.Request) {

	almanac.mu.Lock()
	defer almanac.mu.Unlock()

	//fmt.Printf("%+v\n", almanac.records)

	vars := mux.Vars(r)
	plantName := vars["plant"]
	plantExists := false
	offset := 0
	//fmt.Printf("%s", plantName)
	for _, plant := range almanac.records {

		if plant.PlantName == plantName {
			plantExists = true
			break
		}
		offset++
	}
	if plantExists {
		err := json.NewEncoder(w).Encode(almanac.records[offset])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println("Error")

		} else {
			fmt.Print("Success")
		}
	}
	if !plantExists {
		w.WriteHeader(http.StatusBadRequest)
	}

}
func postPlant(w http.ResponseWriter, r *http.Request) {
	almanac.mu.Lock()
	defer almanac.mu.Unlock()

	var plant Plant
	err := json.NewDecoder(r.Body).Decode(&plant)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Error")
	} else {
		almanac.records = append(almanac.records, plant)

	}
}
*/
