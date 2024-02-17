package Controllers

import (
	"216/internal/Redis"
	"216/internal/Types"
	"216/internal/orchestrator/Database"
	"216/internal/orchestrator/Entities"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

var calculators []*Types.Calculator

func AvailableCalculators(w http.ResponseWriter, r *http.Request) {

	val, err := Redis.Client.Get("calculators").Result()
	if err != nil {
		log.Println(err)
	}

	if err := json.Unmarshal([]byte(val), &calculators); err != nil {
		log.Println(err)
	}

	for k, v := range calculators {
		if time.Unix(int64(v.HeartBeat), 0).Add(10 * time.Second).Before(time.Now()) {
			v.Status = Types.NOTAVAILABLE
			//fmt.Println(v)
			//calculators = append(calculators[:k], calculators[k+1:]...)
			//calculators = append(calculators, v)
		}
		if time.Unix(int64(v.HeartBeat), 0).Add(15 * time.Second).Before(time.Now()) {
			log.Println("удаляем")
			if k < len(calculators) {
				calculators = append(calculators[:k], calculators[k+1:]...)
			} else {
				calculators = calculators[:len(calculators)-1]
			}

		}

	}
	jsonArray, err := json.Marshal(calculators)
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonArray)
	//json.NewEncoder(w).Encode(string(jsonArray))
}

func AvailableComputingResource(w http.ResponseWriter, r *http.Request) {
	var computingResources []Entities.ComputingResource
	Database.Instance.Model(&Entities.ComputingResource{}).Where("to_timestamp(heart_beat) < Now() - interval '10 second'").Update("status", Types.NOTAVAILABLE)
	Database.Instance.Where("to_timestamp(heart_beat) < Now() - interval '20 second'").Delete(&computingResources)
	row := Database.Instance.Order("name").Find(&computingResources) //.Where("to_timestamp(heart_beat) > Now() - interval '4 second'")

	if row.RowsAffected > 0 {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(computingResources)
	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
