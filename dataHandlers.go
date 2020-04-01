package datagraphs

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const provinceOffset = 128

type nationData struct {
	Data                        string `json:"data"`
	Stato                       string `json:"stato"`
	Ricoverati_con_sintomi      int    `json:"ricoverati_con_sintomi"`
	Terapia_intensiva           int    `json:"terapia_intensiva"`
	Totale_ospedalizzati        int    `json:"totale_ospedalizzati"`
	Isolamento_domiciliare      int    `json:"isolamento_domiciliare"`
	Totale_attualmente_positivi int    `json:"totale_attualmente_positivi"`
	Nuovi_positivi              int    `json:"nuovi_positivi"`
	Dimessi_guariti             int    `json:"dimessi_guariti"`
	Deceduti                    int    `json:"deceduti"`
	Totale_casi                 int    `json:"totale_casi"`
	Tamponi                     int    `json:"tamponi"`
	Note_it                     string `json:"note_it"`
}

type regionData struct {
	Data                        string  `json:"data"`
	Stato                       string  `json:"stato"`
	Codice_regione              int     `json:"codice_regione"`
	Denominazione_regione       string  `json:"denominazione_regione"`
	Lat                         float64 `json:"lat"`
	Long                        float64 `json:"long"`
	Ricoverati_con_sintomi      int     `json:"ricoverati_con_sintomi"`
	Terapia_intensiva           int     `json:"terapia_intensiva"`
	Totale_ospedalizzati        int     `json:"totale_ospedalizzati"`
	Isolamento_domiciliare      int     `json:"isolamento_domiciliare"`
	Totale_attualmente_positivi int     `json:"totale_attualmente_positivi"`
	Nuovi_positivi              int     `json:"nuovi_positivi"`
	Dimessi_guariti             int     `json:"dimessi_guariti"`
	Deceduti                    int     `json:"deceduti"`
	Totale_casi                 int     `json:"totale_casi"`
	Tamponi                     int     `json:"tamponi"`
	Note_it                     string  `json:"note_it"`
}

type provinceData struct {
	Data                    string  `json:"data"`
	Stato                   string  `json:"stato"`
	Codice_regione          int     `json:"codice_regione"`
	Denominazione_regione   string  `json:"denominazione_regione"`
	Codice_provincia        int     `json:"codice_provincia"`
	Denominazione_provincia string  `json:"denominazione_provincia"`
	Sigla_provincia         string  `json:"sigla_provincia"`
	Lat                     float64 `json:"lat"`
	Long                    float64 `json:"long"`
	Totale_casi             int     `json:"totale_casi"`
	Note_it                 string  `json:"note_it"`

	NuoviCasi int
}

type noteData struct {
	Codice           string `json:"codice"`
	Data             string `json:"Data"`
	Regione          string `json:"regione"`
	Provincia        string `json:"provincia"`
	Tipologia_avviso string `json:"tipologia_avviso"`
	Avviso           string `json:"avviso"`
	Note             string `json:"noteData"`
}

func getNation() (*[]nationData, error) {
	var response []nationData
	resp, err := http.Get("https://raw.githubusercontent.com/pcm-dpc/COVID-19/master/dati-json/dpc-covid19-ita-andamento-nazionale.json")
	if err != nil {
		return nil, fmt.Errorf("error receiving data: %v", err)
	} else {
		var bodyBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error in received body: %v", err)
		} else {
			err := json.Unmarshal(bodyBytes, &response)
			defer resp.Body.Close()
			if err != nil {
				return nil, fmt.Errorf("error in json unmarshal: %v", err)
			}
		}
	}

	return &response, nil
}

func getRegion() (*[]regionData, error) {
	var response []regionData
	resp, err := http.Get("https://raw.githubusercontent.com/pcm-dpc/COVID-19/master/dati-json/dpc-covid19-ita-regioni.json")
	if err != nil {
		return nil, fmt.Errorf("error receiving data: %v", err)
	} else {
		var bodyBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error in received body: %v", err)
		} else {
			err := json.Unmarshal(bodyBytes, &response)
			defer resp.Body.Close()
			if err != nil {
				return nil, fmt.Errorf("error in json unmarshal: %v", err)
			}
		}
	}
	return &response, nil
}

func getProvince() (*[]provinceData, error) {
	var response []provinceData
	resp, err := http.Get("https://raw.githubusercontent.com/pcm-dpc/COVID-19/master/dati-json/dpc-covid19-ita-province.json")
	if err != nil {
		return nil, fmt.Errorf("error receiving data: %v", err)
	} else {
		var bodyBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error in received body: %v", err)
		} else {
			err := json.Unmarshal(bodyBytes, &response)
			defer resp.Body.Close()
			if err != nil {
				return nil, fmt.Errorf("error in json unmarshal: %v", err)
			}
		}
	}

	setNuoviCasiProvince(&response)
	return &response, nil
}

func getNote() (*[]noteData, error) {
	var notes []noteData
	resp, err := http.Get("https://raw.githubusercontent.com/pcm-dpc/COVID-19/master/note/dpc-covid19-ita-note-it.csv")
	if err != nil {
		return nil, fmt.Errorf("error receiving data: %v", err)
	}
	//bodyBytes, err := ioutil.ReadAll(resp.Body)
	reader := csv.NewReader(bufio.NewReader(resp.Body))
	i := 0
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("error while parsing notes: %v", err)
		}

		if i == 0 {
			i++
			continue
		}
		notes = append(notes, noteData{
			Codice:           line[0],
			Data:             line[1],
			Regione:          line[5],
			Provincia:        line[7],
			Tipologia_avviso: line[9],
			Avviso:           line[10],
			Note:             line[11],
		})
	}

	return &notes, nil
}

func calculateDelta(first int, second int) (float64, string) {
	n := float64(second) - float64(first)
	delta := math.Abs(n)
	return delta, fmt.Sprintf("%+.0f", n)
}

func calculateDeltaPerDay(data *[]nationData, fieldName string, startIndex int) (*[]string, error) {
	deltas := make([]string, 0)
	var deltaStr string

	switch fieldName {
	case "Ricoverati_con_sintomi":
		for i := startIndex; i < len(*data)-1; i++ {
			_, deltaStr = calculateDelta((*data)[i].Ricoverati_con_sintomi, (*data)[i+1].Ricoverati_con_sintomi)
			deltas = append(deltas, deltaStr)
		}
		break
	case "Terapia_intensiva":
		for i := startIndex; i < len(*data)-1; i++ {
			_, deltaStr = calculateDelta((*data)[i].Terapia_intensiva, (*data)[i+1].Terapia_intensiva)
			deltas = append(deltas, deltaStr)
		}
		break
	case "Totale_ospedalizzati":
		for i := startIndex; i < len(*data)-1; i++ {
			_, deltaStr = calculateDelta((*data)[i].Totale_ospedalizzati, (*data)[i+1].Totale_ospedalizzati)
			deltas = append(deltas, deltaStr)
		}
		break
	case "Isolamento_domiciliare":
		for i := startIndex; i < len(*data)-1; i++ {
			_, deltaStr = calculateDelta((*data)[i].Isolamento_domiciliare, (*data)[i+1].Isolamento_domiciliare)
			deltas = append(deltas, deltaStr)
		}
		break
	case "Totale_attualmente_positivi":
		for i := startIndex; i < len(*data)-1; i++ {
			_, deltaStr = calculateDelta((*data)[i].Totale_attualmente_positivi, (*data)[i+1].Totale_attualmente_positivi)
			deltas = append(deltas, deltaStr)
		}
		break
	case "Nuovi_positivi":
		for i := startIndex; i < len(*data)-1; i++ {
			_, deltaStr = calculateDelta((*data)[i].Nuovi_positivi, (*data)[i+1].Nuovi_positivi)
			deltas = append(deltas, deltaStr)
		}
		break
	case "Dimessi_guariti":
		for i := startIndex; i < len(*data)-1; i++ {
			_, deltaStr = calculateDelta((*data)[i].Dimessi_guariti, (*data)[i+1].Dimessi_guariti)
			deltas = append(deltas, deltaStr)
		}
		break
	case "Deceduti":
		for i := startIndex; i < len(*data)-1; i++ {
			_, deltaStr = calculateDelta((*data)[i].Deceduti, (*data)[i+1].Deceduti)
			deltas = append(deltas, deltaStr)
		}
		break
	case "Totale_casi":
		for i := startIndex; i < len(*data)-1; i++ {
			_, deltaStr = calculateDelta((*data)[i].Totale_casi, (*data)[i+1].Totale_casi)
			deltas = append(deltas, deltaStr)
		}
		break
	case "Tamponi":
		for i := startIndex; i < len(*data)-1; i++ {
			_, deltaStr = calculateDelta((*data)[i].Tamponi, (*data)[i+1].Tamponi)
			deltas = append(deltas, deltaStr)
		}
		break
	default:
		return nil, fmt.Errorf("wrong field name passed")
	}

	return &deltas, nil
}

func setNuoviCasiProvince(data *[]provinceData) {
	for i := 0; i < len(*data)-provinceOffset; i++ {
		if i < provinceOffset {
			(*data)[i].NuoviCasi = (*data)[i].Totale_casi
		} else {
			delta, _ := calculateDelta((*data)[i].Totale_casi, (*data)[i+provinceOffset].Totale_casi)
			(*data)[i+provinceOffset].NuoviCasi = int(delta)
		}
	}
}

func intToDeltasArray(data *[]provinceData, startProvinceIndex int) *[]string {
	deltas := make([]string, 0)
	for i := startProvinceIndex; i < len(*data); i += provinceOffset {
		deltas = append(deltas, strconv.Itoa((*data)[i].NuoviCasi))
	}

	return &deltas
}

func findFirstOccurrenceNation(data *[]nationData, fieldName string, toFind interface{}) (int, error) {
	var find interface{}
	switch toFind.(type) {
	case string:
		find = toFind.(string)
		break
	case int:
		find = toFind.(int)
	default:
		return -1, fmt.Errorf("wrong tofind type")
		break
	}

	for i, v := range *data {
		switch fieldName {
		case "Ricoverati_con_sintomi":
			if v.Ricoverati_con_sintomi == find {
				return i, nil
			}
			break
		case "Terapia_intensiva":
			if v.Terapia_intensiva == find {
				return i, nil
			}
			break
		case "Totale_ospedalizzati":
			if v.Totale_ospedalizzati == find {
				return i, nil
			}
			break
		case "Isolamento_domiciliare":
			if v.Isolamento_domiciliare == find {
				return i, nil
			}
			break
		case "Totale_attualmente_positivi":
			if v.Totale_attualmente_positivi == find {
				return i, nil
			}
			break
		case "Nuovi_positivi":
			if v.Nuovi_positivi == find {
				return i, nil
			}
			break
		case "Dimessi_guariti":
			if v.Dimessi_guariti == find {
				return i, nil
			}
			break
		case "Deceduti":
			if v.Deceduti == find {
				return i, nil
			}
			break
		case "Totale_casi":
			if v.Totale_casi == find {
				return i, nil
			}
			break
		case "Tamponi":
			if v.Tamponi == find {
				return i, nil
			}
			break
		default:
			return -1, fmt.Errorf("wrong field name passed")
		}
	}
	return -1, fmt.Errorf("element not found")
}

func findFirstOccurrenceRegion(data *[]regionData, fieldName string, toFind interface{}) (int, error) {
	var find interface{}
	switch toFind.(type) {
	case string:
		find = toFind.(string)
		break
	case int:
		find = toFind.(int)
	default:
		return -1, fmt.Errorf("wrong tofind type")
		break
	}

	for i, v := range *data {
		switch strings.ToLower(fieldName) {
		case "codice_regione":
			if v.Codice_regione == find {
				return i, nil
			}
			break
		case "denominazione_regione":
			if v.Denominazione_regione == find {
				return i, nil
			}
			break
		case "ricoverati_con_sintomi":
			if v.Ricoverati_con_sintomi == find {
				return i, nil
			}
			break
		case "terapia_intensiva":
			if v.Terapia_intensiva == find {
				return i, nil
			}
			break
		case "totale_ospedalizzati":
			if v.Totale_ospedalizzati == find {
				return i, nil
			}
			break
		case "isolamento_domiciliare":
			if v.Isolamento_domiciliare == find {
				return i, nil
			}
			break
		case "totale_attualmente_positivi":
			if v.Totale_attualmente_positivi == find {
				return i, nil
			}
			break
		case "nuovi_positivi":
			if v.Nuovi_positivi == find {
				return i, nil
			}
			break
		case "dimessi_guariti":
			if v.Dimessi_guariti == find {
				return i, nil
			}
			break
		case "deceduti":
			if v.Deceduti == find {
				return i, nil
			}
			break
		case "totale_casi":
			if v.Totale_casi == find {
				return i, nil
			}
			break
		case "tamponi":
			if v.Tamponi == find {
				return i, nil
			}
			break
		default:
			return -1, fmt.Errorf("wrong field name passed")
		}
	}
	return -1, fmt.Errorf("element not found")
}

func findFirstOccurrenceProvince(data *[]provinceData, fieldName string, toFind interface{}) (int, error) {
	var find interface{}
	switch toFind.(type) {
	case string:
		find = toFind.(string)
		break
	case int:
		find = toFind.(int)
		break
	default:
		return -1, fmt.Errorf("wrong tofind type")
		break
	}

	for i, v := range *data {
		switch strings.ToLower(fieldName) {
		case "codice_regione":
			if v.Codice_regione == find {
				return i, nil
			}
		case "denominazione_provincia":
			if v.Denominazione_provincia == find {
				return i, nil
			}
		case "sigla_provincia":
			if v.Sigla_provincia == find {
				return i, nil
			}
		case "totale_casi":
			if v.Totale_casi == find {
				return i, nil
			}
		default:
			break
		}
	}

	return -1, fmt.Errorf("element not found")
}

func findFirstOccurrenceNote(data *[]noteData, fieldName string, toFind interface{}) (int, error) {
	var find interface{}
	switch toFind.(type) {
	case string:
		find = toFind.(string)
		break
	case int:
		find = toFind.(int)
		break
	default:
		return -1, fmt.Errorf("wrong tofind type")
		break
	}

	for i, v := range *data {
		switch strings.ToLower(fieldName) {
		case "codice":
			if v.Codice == find {
				return i, nil
			}
		case "data":
			if v.Data == find {
				return i, nil
			}
		case "regione":
			if v.Regione == find {
				return i, nil
			}
		case "provincia":
			if v.Provincia == find {
				return i, nil
			}
		case "tipologia_avviso":
			if v.Tipologia_avviso == find {
				return i, nil
			}
		case "avviso":
			if v.Avviso == find {
				return i, nil
			}
		case "note":
			if v.Note == find {
				return i, nil
			}
		default:
			break
		}
	}

	return -1, fmt.Errorf("element not found")
}

func deleteFile(filename string) error {
	err := os.Remove(filename)
	if err != nil {
		return fmt.Errorf("error deleting file: %v", err)
	}

	return nil
}
