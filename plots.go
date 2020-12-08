package covidgraphs

import (
	"fmt"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
	"os"
	"strings"
	"time"
)

const (
	daySwitch   = "06:00:00"
	nightSwitch = "19:00:00"
)

// Calculates annotations containing the difference to the previous point on the plot
func deltaAnnotations(deltas *[]string, xValues *[]time.Time, yValues *[]float64) chart.AnnotationSeries {
	value2 := make([]chart.Value2, 0)
	for i := 0; i < len(*deltas)-1; i++ {
		value2 = append(value2, chart.Value2{
			Style:  chart.Style{},
			Label:  (*deltas)[i],
			XValue: chart.TimeToFloat64((*xValues)[i+1]),
			YValue: (*yValues)[i+1],
		})
	}
	return chart.AnnotationSeries{Annotations: value2}
}

// Creates a plot with the given series
func timeseriesChart(charts *[]chart.TimeSeries, gridLines *[]chart.GridLine, annotations *[]chart.AnnotationSeries, title, filename, xAxisName, yAxisName string) (error, string) {
	series := make([]chart.Series, 0)
	for _, v := range *charts {
		series = append(series, v)
	}
	for _, v := range *annotations {
		series = append(series, v)
	}
	
	var backgroundColor drawing.Color
	var fontsColor drawing.Color
	var colorMode string
	lightHour, _ := time.Parse("15:04:05", daySwitch)
	darkHour, _ := time.Parse("15:04:05", nightSwitch)
	now := time.Now().Hour()
	if now >= darkHour.Hour() || now < lightHour.Hour() {
		backgroundColor = chart.ColorBlack
		fontsColor = chart.ColorWhite
		colorMode = "dark"
	} else {
		backgroundColor = chart.ColorWhite
		fontsColor = chart.ColorBlack
		colorMode = "light"
	}

	graph := chart.Chart{
		Title:  title,
		Width:  1280,
		Height: 720,
		TitleStyle: chart.Style{
			FontColor: fontsColor,
		},
		Background: chart.Style{
			FillColor: backgroundColor,
			Padding: chart.Box{
				Top: 50,
			},
		},
		Canvas: chart.Style{
			FillColor: backgroundColor,
		},
		XAxis: chart.XAxis{
			Name: xAxisName,
			Style: chart.Style{
				StrokeColor: fontsColor,
			},
			ValueFormatter: chart.TimeDateValueFormatter,
			GridMajorStyle: chart.Style{
				StrokeColor: chart.ColorAlternateGray,
				StrokeWidth: 1.0,
			},
			TickStyle: chart.Style{
				FontColor: fontsColor,
				FontSize:  15,
			},
			GridLines: *gridLines,
		},
		YAxis: chart.YAxis{
			Name: yAxisName,
			Style: chart.Style{
				StrokeColor: fontsColor,
			},
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("%d", int(v.(float64)))
			},
			TickStyle: chart.Style{
				TextRotationDegrees: 45.0,
				FontColor:           fontsColor,
				FontSize:            15,
			},
		},
		Series: series,
	}

	graph.Elements = []chart.Renderable{chart.Legend(&graph, chart.Style{
		FontSize: 15,
	})}

	if filename==""{
		filename=title+"-"+colorMode+".png"
	}
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error while creating file: %v", err), ""
	}
	defer f.Close()
	err = graph.Render(chart.PNG, f)
	if err != nil {
		return fmt.Errorf("error while rendering graph: %v", err), ""
	}
	return nil, filename
}

// Converts dates to Float64 to fit the X Axis of the plots
func dateXAxis(date *[]chart.GridLine, newDate time.Time) *[]chart.GridLine {
	*date = append(*date, chart.GridLine{
		Value: chart.TimeToFloat64(newDate),
	})

	return date
}

// Creates series of points according to the national data
func nationToTimeseries(data *[]NationData, fieldName string, index int) (*[]time.Time, *[]float64, *[]chart.GridLine, error) {
	date := make([]time.Time, 0)
	values := make([]float64, 0)
	dateAxis := make([]chart.GridLine, 0)
	for i := index; i < len(*data); i++ {
		dateRead, err := time.Parse("2006-01-02T15:04:05", (*data)[i].Data)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("error converting date string to date: %v", err)
		}
		year, month, day := dateRead.Date()
		dateRead = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		date = append(date, dateRead)
		dateAxis = *dateXAxis(&dateAxis, dateRead)

		switch strings.ToLower(fieldName) {
		case "ricoverati_con_sintomi":
			values = append(values, float64((*data)[i].Ricoverati_con_sintomi))
			break
		case "terapia_intensiva":
			values = append(values, float64((*data)[i].Terapia_intensiva))
			break
		case "totale_ospedalizzati":
			values = append(values, float64((*data)[i].Totale_ospedalizzati))
			break
		case "isolamento_domiciliare":
			values = append(values, float64((*data)[i].Isolamento_domiciliare))
			break
		case "attualmente_positivi":
			values = append(values, float64((*data)[i].Totale_positivi))
			break
		case "nuovi_positivi":
			values = append(values, float64((*data)[i].Nuovi_positivi))
			break
		case "dimessi_guariti":
			values = append(values, float64((*data)[i].Dimessi_guariti))
			break
		case "deceduti":
			values = append(values, float64((*data)[i].Deceduti))
			break
		case "totale_casi":
			values = append(values, float64((*data)[i].Totale_casi))
			break
		case "tamponi":
			values = append(values, float64((*data)[i].Tamponi))
			break
		default:
			return nil, nil, nil, fmt.Errorf("wrong field name passed")
		}
	}

	return &date, &values, &dateAxis, nil
}

// Creates series of points according to the regional data
func regionToTimeseries(data *[]RegionData, fieldName string, index int, startRegionCodeIndex int) (*[]time.Time, *[]float64, *[]chart.GridLine, error) {
	if startRegionCodeIndex < 0 || startRegionCodeIndex > 21 {
		return nil, nil, nil, fmt.Errorf("region index out of range")
	}

	date := make([]time.Time, 0)
	values := make([]float64, 0)
	dateAxis := make([]chart.GridLine, 0)

	offset := index % 21
	var realIndex int
	if offset == startRegionCodeIndex {
		realIndex = index
	} else if offset > startRegionCodeIndex {
		realIndex = index + 21 - (offset - startRegionCodeIndex)
	} else {
		realIndex = index + (startRegionCodeIndex - offset)
	}

	for i := realIndex; i < len(*data); i += 21 {
		dateRead, err := time.Parse("2006-01-02T15:04:05", (*data)[i].Data)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("error converting date string to date: %v", err)
		}
		year, month, day := dateRead.Date()
		dateRead = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		date = append(date, dateRead)
		dateAxis = *dateXAxis(&dateAxis, dateRead)

		switch strings.ToLower(fieldName) {
		case "ricoverati_con_sintomi":
			values = append(values, float64((*data)[i].Ricoverati_con_sintomi))
			break
		case "terapia_intensiva":
			values = append(values, float64((*data)[i].Terapia_intensiva))
			break
		case "totale_ospedalizzati":
			values = append(values, float64((*data)[i].Totale_ospedalizzati))
			break
		case "isolamento_domiciliare":
			values = append(values, float64((*data)[i].Isolamento_domiciliare))
			break
		case "attualmente_positivi":
			values = append(values, float64((*data)[i].Totale_positivi))
			break
		case "nuovi_positivi":
			values = append(values, float64((*data)[i].Nuovi_positivi))
			break
		case "dimessi_guariti":
			values = append(values, float64((*data)[i].Dimessi_guariti))
			break
		case "deceduti":
			values = append(values, float64((*data)[i].Deceduti))
			break
		case "totale_casi":
			values = append(values, float64((*data)[i].Totale_casi))
			break
		case "tamponi":
			values = append(values, float64((*data)[i].Tamponi))
			break
		default:
			return nil, nil, nil, fmt.Errorf("wrong field name passed")
		}
	}

	return &date, &values, &dateAxis, nil
}

// Creates series of points according to the provincial data
func provinceToTimeseries(data *[]ProvinceData, fieldName string, provinceIndexes *[]int) (*[]time.Time, *[]float64, *[]chart.GridLine, error) {
	date := make([]time.Time, 0)
	values := make([]float64, 0)
	dateAxis := make([]chart.GridLine, 0)
	for _,v:=range *provinceIndexes {
		dateRead, err := time.Parse("2006-01-02T15:04:05", (*data)[v].Data)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("error converting date string to date: %v", err)
		}
		year, month, day := dateRead.Date()
		dateRead = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		date = append(date, dateRead)
		dateAxis = *dateXAxis(&dateAxis, dateRead)

		switch strings.ToLower(fieldName) {
		case "totale_casi":
			values = append(values, float64((*data)[v].Totale_casi))
			break
		case "nuovi_positivi":
			values = append(values, float64((*data)[v].NuoviCasi))
			break
		default:
			return nil, nil, nil, fmt.Errorf("wrong field name passed")
		}
	}

	return &date, &values, &dateAxis, nil
}
