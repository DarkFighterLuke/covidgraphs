package covidgraphs

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
	"sort"
	"strconv"
	"strings"
	"time"
)

// National data struct containing fields from the parsed JSON
type NationData struct {
	Data                   string `json:"data"`
	Stato                  string `json:"stato"`
	Ricoverati_con_sintomi int    `json:"ricoverati_con_sintomi"`
	Terapia_intensiva      int    `json:"terapia_intensiva"`
	Totale_ospedalizzati   int    `json:"totale_ospedalizzati"`
	Isolamento_domiciliare int    `json:"isolamento_domiciliare"`
	Totale_positivi        int    `json:"totale_positivi"`
	Nuovi_positivi         int    `json:"nuovi_positivi"`
	Dimessi_guariti        int    `json:"dimessi_guariti"`
	Deceduti               int    `json:"deceduti"`
	Totale_casi            int    `json:"totale_casi"`
	Tamponi                int    `json:"tamponi"`
	Note_it                string `json:"note_it"`
}

// Regional data struct containing fields from the parsed JSON
type RegionData struct {
	Data                   string  `json:"data"`
	Stato                  string  `json:"stato"`
	Codice_regione         int     `json:"codice_regione"`
	Denominazione_regione  string  `json:"denominazione_regione"`
	Lat                    float64 `json:"lat"`
	Long                   float64 `json:"long"`
	Ricoverati_con_sintomi int     `json:"ricoverati_con_sintomi"`
	Terapia_intensiva      int     `json:"terapia_intensiva"`
	Totale_ospedalizzati   int     `json:"totale_ospedalizzati"`
	Isolamento_domiciliare int     `json:"isolamento_domiciliare"`
	Totale_positivi        int     `json:"totale_positivi"`
	Nuovi_positivi         int     `json:"nuovi_positivi"`
	Dimessi_guariti        int     `json:"dimessi_guariti"`
	Deceduti               int     `json:"deceduti"`
	Totale_casi            int     `json:"totale_casi"`
	Tamponi                int     `json:"tamponi"`
	Note_it                string  `json:"note_it"`
}

// Provincial data struct containing fields from the parsed JSON
type ProvinceData struct {
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

// Notes data struct containing fields from the parsed CSV
type NoteData struct {
	Codice           string `json:"codice"`
	Data             string `json:"Data"`
	Regione          string `json:"regione"`
	Provincia        string `json:"provincia"`
	Tipologia_avviso string `json:"tipologia_avviso"`
	Avviso           string `json:"avviso"`
	Note             string `json:"NoteData"`
}

var lastUpdateNation time.Time
var lastUpdateRegions time.Time
var lastUpdateProvinces time.Time
var lastUpdateNotes time.Time


// Retrieves and parses nation data from the pcm repo
func GetNation() (*[]NationData, error) {
	var response []NationData
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

	lastUpdateNation = time.Now()
	return &response, nil
}

// Retrieves and parses regions data from the pcm repo
func GetRegions() (*[]RegionData, error) {
	var response []RegionData
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

	lastUpdateRegions = time.Now()
	return &response, nil
}

// Retrieves and parses provinces data from the pcm repo
func GetProvinces() (*[]ProvinceData, error) {
	var response []ProvinceData
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

	//eraseTrashProvinces(&response)
	setNuoviCasiProvince(&response)
	lastUpdateProvinces = time.Now()
	return &response, nil
}

// Retrieves and parses notes data from the pcm repo
func GetNotes() (*[]NoteData, error) {
	var notes []NoteData
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
		notes = append(notes, NoteData{
			Codice:           line[0],
			Data:             line[1],
			Regione:          line[5],
			Provincia:        line[7],
			Tipologia_avviso: line[9],
			Avviso:           line[10],
			Note:             line[11],
		})
	}

	lastUpdateNotes = time.Now()
	return &notes, nil
}

// Calculates delta between two integer quantities
func CalculateDelta(first int, second int) (float64, string) {
	n := float64(second) - float64(first)
	delta := math.Abs(n)
	return delta, fmt.Sprintf("%+.0f", n)
}

// Calculates delta per day of the specified nation data field
func calculateDeltaPerDay(data *[]NationData, fieldName string, startIndex int) (*[]string, error) {
	deltas := make([]string, 0)
	var deltaStr string

	switch fieldName {
	case "Ricoverati_con_sintomi":
		for i := startIndex; i < len(*data)-1; i++ {
			_, deltaStr = CalculateDelta((*data)[i].Ricoverati_con_sintomi, (*data)[i+1].Ricoverati_con_sintomi)
			deltas = append(deltas, deltaStr)
		}
		break
	case "Terapia_intensiva":
		for i := startIndex; i < len(*data)-1; i++ {
			_, deltaStr = CalculateDelta((*data)[i].Terapia_intensiva, (*data)[i+1].Terapia_intensiva)
			deltas = append(deltas, deltaStr)
		}
		break
	case "Totale_ospedalizzati":
		for i := startIndex; i < len(*data)-1; i++ {
			_, deltaStr = CalculateDelta((*data)[i].Totale_ospedalizzati, (*data)[i+1].Totale_ospedalizzati)
			deltas = append(deltas, deltaStr)
		}
		break
	case "Isolamento_domiciliare":
		for i := startIndex; i < len(*data)-1; i++ {
			_, deltaStr = CalculateDelta((*data)[i].Isolamento_domiciliare, (*data)[i+1].Isolamento_domiciliare)
			deltas = append(deltas, deltaStr)
		}
		break
	case "attualmente_positivi":
		for i := startIndex; i < len(*data)-1; i++ {
			_, deltaStr = CalculateDelta((*data)[i].Totale_positivi, (*data)[i+1].Totale_positivi)
			deltas = append(deltas, deltaStr)
		}
		break
	case "Nuovi_positivi":
		for i := startIndex; i < len(*data)-1; i++ {
			_, deltaStr = CalculateDelta((*data)[i].Nuovi_positivi, (*data)[i+1].Nuovi_positivi)
			deltas = append(deltas, deltaStr)
		}
		break
	case "Dimessi_guariti":
		for i := startIndex; i < len(*data)-1; i++ {
			_, deltaStr = CalculateDelta((*data)[i].Dimessi_guariti, (*data)[i+1].Dimessi_guariti)
			deltas = append(deltas, deltaStr)
		}
		break
	case "Deceduti":
		for i := startIndex; i < len(*data)-1; i++ {
			_, deltaStr = CalculateDelta((*data)[i].Deceduti, (*data)[i+1].Deceduti)
			deltas = append(deltas, deltaStr)
		}
		break
	case "Totale_casi":
		for i := startIndex; i < len(*data)-1; i++ {
			_, deltaStr = CalculateDelta((*data)[i].Totale_casi, (*data)[i+1].Totale_casi)
			deltas = append(deltas, deltaStr)
		}
		break
	case "Tamponi":
		for i := startIndex; i < len(*data)-1; i++ {
			_, deltaStr = CalculateDelta((*data)[i].Tamponi, (*data)[i+1].Tamponi)
			deltas = append(deltas, deltaStr)
		}
		break
	default:
		return nil, fmt.Errorf("wrong field name passed")
	}

	return &deltas, nil
}

// Creates an array of delta for creating annotations
func intToDeltasArray(data *[]ProvinceData, provinceIndexes *[]int) *[]string {
	deltas := make([]string, 0)
	for i := range *provinceIndexes {
		deltas = append(deltas, strconv.Itoa((*data)[i].NuoviCasi))
	}

	return &deltas
}

// Sets the artificial NuoviCasi field for provinces
func setNuoviCasiProvince(data *[]ProvinceData) {
	//set -1 to all indicating as not already filled
	for i, _ := range (*data) {
		(*data)[i].NuoviCasi = -1
	}

	for i := len(*data) - 1; i > 0; i-- {
		if (*data)[i].NuoviCasi == -1 &&
			(*data)[i].Denominazione_provincia != "Fuori Regione / Provincia Autonoma" &&
			(*data)[i].Denominazione_provincia != "In fase di definizione/aggiornamento" {
			var lastOccurrenceIndex = i
			for j := i; j > 0; j-- {
				if (*data)[j].Denominazione_provincia == (*data)[lastOccurrenceIndex].Denominazione_provincia {
					delta, _ := CalculateDelta((*data)[j].Totale_casi, (*data)[lastOccurrenceIndex].Totale_casi)
					(*data)[lastOccurrenceIndex].NuoviCasi = int(delta)
					lastOccurrenceIndex = j
				}
			}
			firstDayIndex, err := FindFirstOccurrenceProvince(data, "denominazione_provincia", (*data)[i].Denominazione_provincia)
			if err != nil {
				fmt.Println("Error setting NuoviCasi for provinces: ", err)
			}
			(*data)[firstDayIndex].NuoviCasi = (*data)[firstDayIndex].Totale_casi
		}
	}
}

// Finds the first occurence in the nation data array for the specified field
func FindFirstOccurrenceNation(data *[]NationData, fieldName string, toFind interface{}) (int, error) {
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
		case "attualmente_positivi":
			if v.Totale_positivi == find {
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

// Finds the first occurence in the regions data array for the specified field
func FindFirstOccurrenceRegion(data *[]RegionData, fieldName string, toFind interface{}) (int, error) {
	var find interface{}
	switch toFind.(type) {
	case string:
		find = strings.ToLower(toFind.(string))
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
			break
		case "denominazione_regione":
			if strings.Replace(strings.ToLower(v.Denominazione_regione), "-", " ", -1) == strings.Replace(strings.ToLower(find.(string)), "-", " ", -1) {
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
		case "attualmente_positivi":
			if v.Totale_positivi == find {
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

// Finds the first occurence in the provinces data array for the specified field
func FindFirstOccurrenceProvince(data *[]ProvinceData, fieldName string, toFind interface{}) (int, error) {
	var find interface{}
	switch toFind.(type) {
	case string:
		find = strings.ToLower(toFind.(string))
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
			if strings.Replace(strings.ToLower(v.Denominazione_provincia), "-", " ", -1) == strings.Replace(strings.ToLower(find.(string)), "-", " ", -1) {
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

// Finds the first occurence in the notes data array for the specified field
func FindFirstOccurrenceNote(data *[]NoteData, fieldName string, toFind interface{}) (int, error) {
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

// Returns regions names list
func GetRegionsNamesList(data *[]RegionData) []string {
	regionsNames := make([]string, 0)
	for _, v := range *data {
		regionsNames = append(regionsNames, v.Denominazione_regione)
	}

	return regionsNames
}

// Returns northern regions names list
func GetNordRegionsNamesList() []string {
	nordRegions := []string{"Piemonte", "Valle d'Aosta", "Liguria", "Lombardia", "P.A. Trento", "P.A. Bolzano", "Veneto", "Friuli Venezia Giulia", "Emilia Romagna"}
	return nordRegions
}

// Returns central regions names list
func GetCentroRegionsNamesList() []string {
	centroRegions := []string{"Toscana", "Umbria", "Marche", "Lazio"}
	return centroRegions
}

// Returns southern regions names list
func GetSudRegionsNamesList() []string {
	sudRegions := []string{"Abruzzo", "Molise", "Campania", "Puglia", "Basilicata", "Calabria", "Sicilia", "Sardegna"}
	return sudRegions
}

// Returns top regions according to field totale_contagi
func GetTopTenRegionsTotaleContagi(data *[]RegionData) *[]RegionData {
	latestData := make([]RegionData, 21)
	copy(latestData, (*data)[len(*data)-21:len(*data)])

	sort.Slice(latestData, func(i, j int) bool {
		return latestData[i].Totale_casi > latestData[j].Totale_casi
	})

	return &latestData
}

// Returns top provinces according to field totale_casi
func GetTopTenProvincesTotaleContagi(data *[]ProvinceData) *[]ProvinceData {
	latestData:=make([]ProvinceData, 0)
	latestDate, _:=time.Parse("2006-01-02T15:04:05", (*data)[len(*data)-1].Data)
	for i:=len(*data)-1; i>0; i--{
		date, _:=time.Parse("2006-01-02T15:04:05", (*data)[i].Data)
		if latestDate.Day()-date.Day()>0{
			break
		}else{
			latestData=append(latestData, (*data)[i])
		}
	}

	sort.Slice(latestData, func(i, j int) bool {
		return latestData[i].Totale_casi > latestData[j].Totale_casi
	})

	return &latestData
}

// Finds the last occurence in the regions data array for the specified field
func FindLastOccurrenceRegion(data *[]RegionData, fieldName string, toFind interface{}) (int, error) {
	latestData := (*data)[len(*data)-21 : len(*data)]

	var find interface{}
	switch toFind.(type) {
	case string:
		find = strings.ToLower(toFind.(string))
		break
	case int:
		find = toFind.(int)
	default:
		return -1, fmt.Errorf("wrong tofind type")
		break
	}

	for i, v := range latestData {
		switch strings.ToLower(fieldName) {
		case "codice_regione":
			if v.Codice_regione == find {
				return i + len(*data) - 21, nil
			}
			break
		case "denominazione_regione":
			if strings.Replace(strings.ToLower(v.Denominazione_regione), "-", " ", -1) == strings.Replace(strings.ToLower(find.(string)), "-", " ", -1) {
				return i + len(*data) - 21, nil
			}
			break
		case "ricoverati_con_sintomi":
			if v.Ricoverati_con_sintomi == find {
				return i + len(*data) - 21, nil
			}
			break
		case "terapia_intensiva":
			if v.Terapia_intensiva == find {
				return i + len(*data) - 21, nil
			}
			break
		case "totale_ospedalizzati":
			if v.Totale_ospedalizzati == find {
				return i + len(*data) - 21, nil
			}
			break
		case "isolamento_domiciliare":
			if v.Isolamento_domiciliare == find {
				return i + len(*data) - 21, nil
			}
			break
		case "attualmente_positivi":
			if v.Totale_positivi == find {
				return i + len(*data) - 21, nil
			}
			break
		case "nuovi_positivi":
			if v.Nuovi_positivi == find {
				return i + len(*data) - 21, nil
			}
			break
		case "dimessi_guariti":
			if v.Dimessi_guariti == find {
				return i + len(*data) - 21, nil
			}
			break
		case "deceduti":
			if v.Deceduti == find {
				return i + len(*data) - 21, nil
			}
			break
		case "totale_casi":
			if v.Totale_casi == find {
				return i + len(*data) - 21, nil
			}
			break
		case "tamponi":
			if v.Tamponi == find {
				return i + len(*data) - 21, nil
			}
			break
		default:
			return -1, fmt.Errorf("wrong field name passed")
		}
	}
	return -1, fmt.Errorf("element not found")
}

// Finds the last occurence in the provinces data array for the specified field
func FindLastOccurrenceProvince(data *[]ProvinceData, fieldName string, toFind interface{}) (int, error) {
	var find interface{}
	switch toFind.(type) {
	case string:
		find = strings.ToLower(toFind.(string))
		break
	case int:
		find = toFind.(int)
		break
	default:
		return -1, fmt.Errorf("wrong toFind type")
		break
	}

	for i := len(*data)-1; i > 0; i-- {
		switch strings.ToLower(fieldName) {
		case "codice_regione":
			if (*data)[i].Codice_regione == find {
				return i, nil
			}
		case "denominazione_provincia":
			if strings.Replace(strings.ToLower((*data)[i].Denominazione_provincia), "-", " ", -1) == strings.Replace(strings.ToLower(find.(string)), "-", " ", -1) {
				return i, nil
			}
		case "sigla_provincia":
			if (*data)[i].Sigla_provincia == find {
				return i, nil
			}
		case "totale_casi":
			if (*data)[i].Totale_casi == find {
				return i, nil
			}
		default:
			break
		}
	}

	return -1, fmt.Errorf("element not found")
}

// Returns the last provinces data according to the given region name
func GetLastProvincesByRegionName(data *[]ProvinceData, regionName string) *[]ProvinceData {
	provinces:=make([]ProvinceData, 0)

	latestDate, _:=time.Parse("2006-01-02T15:04:05", (*data)[len(*data)-1].Data)
	var firstLatestIndex int
	for i:=len(*data)-1; i>0; i--{
		date, _:=time.Parse("2006-01-02T15:04:05", (*data)[i].Data)
		if latestDate.Day()-date.Day()>0{
			firstLatestIndex=i+1
			break
		}
	}

	for i := firstLatestIndex; i<len(*data); i++ {
		if strings.ToLower((*data)[i].Denominazione_regione) == strings.ToLower(regionName) &&
			strings.ToLower((*data)[i].Denominazione_provincia) != "in fase di definizione/aggiornamento" &&
			strings.ToLower((*data)[i].Denominazione_provincia) != "fuori regione / provincia autonoma" {
			provinces = append(provinces, (*data)[i])
		}
	}

	return &provinces
}

// Returns a slice with all the given province data
func GetProvinceIndexesByName(data *[]ProvinceData, provinceName string) *[]int {
	provinceIndexes := make([]int, 0)
	for i, pData := range *data {
		if strings.ToLower(pData.Denominazione_provincia) == strings.ToLower(provinceName) {
			provinceIndexes = append(provinceIndexes, i)
		}
	}

	return &provinceIndexes
}

// Deletes the specified file
func DeleteFile(filename string) error {
	err := os.Remove(filename)
	if err != nil {
		return fmt.Errorf("error deleting file: %v", err)
	}

	return nil
}

// Deletes plots folder and recreates it
func DeleteAllPlots(folder string) {
	path, _ := os.Getwd()
	path = path + folder
	os.RemoveAll(path)
	os.Mkdir(path, 0755)
}
