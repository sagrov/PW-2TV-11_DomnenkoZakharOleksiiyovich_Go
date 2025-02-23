package main

import (
	"encoding/json"
	"fmt"
	_ "math"
	"net/http"
	"strconv"
)

type Input struct {
	Values []float64 `json:"values"`
}

var fuelStorage = map[string]map[string]float64{
	"coal": {
		"qri":  20.47,
		"avin": 0.8,
		"ar":   25.20,
		"gvin": 1.5,
	},
	"fuel_oil": {
		"qri":  40.40,
		"avin": 1.0,
		"ar":   0.15,
		"gvin": 0.0,
	},
	"natural_gas": {
		"qri":  33.08,
		"avin": 1.25,
		"ar":   0.0,
		"gvin": 0.0,
	},
}

func calculateEmission(fuelId string, value float64) ([]float64, error) {
	fuelData, exists := fuelStorage[fuelId]
	if !exists {
		return nil, fmt.Errorf("Fuel ID not found")
	}

	qri, qriExists := fuelData["qri"]
	avin, avinExists := fuelData["avin"]
	ar, arExists := fuelData["ar"]
	gvin, gvinExists := fuelData["gvin"]
	if !qriExists || !avinExists || !arExists || !gvinExists {
		return nil, fmt.Errorf("Missing required fuel data")
	}
	ktv := (1e6 / qri) * avin * (ar / (100 - gvin)) * (1 - 0.985)
	etv := 1e-6 * ktv * qri * value
	return []float64{round(ktv), round(etv)}, nil
}

func round(val float64) float64 {
	res, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", val), 64)
	return res
}

func calculateTask1(coal, fuelOil, naturalGas float64) string {
	var output string
	if ktvEtv, err := calculateEmission("coal", coal); err == nil {
		output += fmt.Sprintf("Показник емісії твердих частинок при спалюванні вугілля становитиме: %.2f г/ГДж\n", ktvEtv[0])
		output += fmt.Sprintf("Валовий викид при спалюванні вугілля становитиме: %.2f т\n", ktvEtv[1])
	}
	if ktvEtv, err := calculateEmission("fuel_oil", fuelOil); err == nil {
		output += fmt.Sprintf("Показник емісії твердих частинок при спалюванні мазуту становитиме: %.2f г/ГДж\n", ktvEtv[0])
		output += fmt.Sprintf("Валовий викид при спалюванні мазуту становитиме: %.2f т\n", ktvEtv[1])
	}
	if ktvEtv, err := calculateEmission("natural_gas", naturalGas); err == nil {
		output += fmt.Sprintf("Показник емісії твердих частинок при спалюванні природнього газу становитиме: %.2f г/ГДж\n", ktvEtv[0])
		output += fmt.Sprintf("Валовий викид при спалюванні природнього газу становитиме: %.2f т\n", ktvEtv[1])
	}
	return output
}

func calculator1Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var input Input
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if len(input.Values) != 3 {
		http.Error(w, "Invalid number of inputs", http.StatusBadRequest)
		return
	}
	result := calculateTask1(input.Values[0], input.Values[1], input.Values[2])

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"result": result})
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/api/calculator1", calculator1Handler)

	fmt.Println("Server running at http://localhost:8082")
	http.ListenAndServe(":8082", nil)
}
