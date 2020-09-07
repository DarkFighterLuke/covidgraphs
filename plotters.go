package covidgraphs

import (
	"fmt"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
	"os"
	"strings"
	"time"
)

// Returns a plot including total cases, healed and dead
func AndamentoNazionaleCompleto(data *[]NationData, title, filename string) (error, string) {
	xAxisName := ""
	yAxisName := ""

	fieldName := "Totale_casi"
	xTotale, yTotale, xNames, err := nationToTimeseries(data, "Totale_casi", 0)
	if err != nil {
		return fmt.Errorf("error while creating %v chart: %v", fieldName, err), ""
	}

	fieldName = "Dimessi_guariti"
	xGuariti, yGuariti, _, err := nationToTimeseries(data, fieldName, 0)
	if err != nil {
		return fmt.Errorf("error while creating %v chart: %v", fieldName, err), ""
	}

	fieldName = "Deceduti"
	xDeceduti, yDeceduti, _, err := nationToTimeseries(data, fieldName, 0)
	if err != nil {
		return fmt.Errorf("error while creating %v chart: %v", fieldName, err), ""
	}

	totale := chart.TimeSeries{
		Name: "Totale contagi",
		Style: chart.Style{
			StrokeColor: drawing.Color{
				R: 255,
				A: 255,
			},
			FillColor: drawing.Color{
				R: 255,
				A: 150,
			},
		},
		YAxis:   0,
		XValues: *xTotale,
		YValues: *yTotale,
	}

	guariti := chart.TimeSeries{
		Name: "Guariti",
		Style: chart.Style{
			StrokeColor: chart.ColorGreen,
			FillColor:   chart.ColorGreen.WithAlpha(175),
		},
		YAxis:   0,
		XValues: *xGuariti,
		YValues: *yGuariti,
	}

	deceduti := chart.TimeSeries{
		Name: "Morti",
		Style: chart.Style{
			StrokeColor: chart.ColorAlternateGray,
			FillColor:   chart.ColorAlternateGray.WithAlpha(200),
		},
		YAxis:   0,
		XValues: *xDeceduti,
		YValues: *yDeceduti,
	}

	series := make([]chart.TimeSeries, 3)
	series[0] = totale
	series[1] = guariti
	series[2] = deceduti

	annotations := make([]chart.AnnotationSeries, 0)

	err, fileName := timeseriesChart(&series, xNames, &annotations, title, filename, xAxisName, yAxisName)
	if err != nil {
		return fmt.Errorf("%v", err), ""
	}
	return nil, fileName
}

// Returns a plot with the national data according to the specified fields
func VociNazione(data *[]NationData, fieldName []string, nationIndex int, title, filename string) (error, string) {
	xAxisName := ""
	yAxisName := ""

	var err error
	var xValues *[]time.Time
	var yValues *[]float64
	var xNames *[]chart.GridLine
	series := make([]chart.TimeSeries, 0)
	var color drawing.Color
	var alpha uint8 = 200
	var fileName string
	for _, v := range fieldName {
		switch strings.ToLower(v) {
		case "ricoverati_con_sintomi":
			xValues, yValues, xNames, err = nationToTimeseries(data, v, nationIndex)
			color = drawing.Color{38, 224, 175, 255}

			series = append(series, chart.TimeSeries{
				Name: v,
				Style: chart.Style{
					StrokeColor: color,
					FillColor:   color.WithAlpha(alpha),
				},
				YAxis:   0,
				XValues: *xValues,
				YValues: *yValues,
			})
			break
		case "terapia_intensiva":
			xValues, yValues, xNames, err = nationToTimeseries(data, v, nationIndex)
			color = drawing.Color{88, 22, 115, 255}

			series = append(series, chart.TimeSeries{
				Name: v,
				Style: chart.Style{
					StrokeColor: color,
					FillColor:   color.WithAlpha(alpha),
				},
				YAxis:   0,
				XValues: *xValues,
				YValues: *yValues,
			})
			break
		case "totale_ospedalizzati":
			xValues, yValues, xNames, err = nationToTimeseries(data, v, nationIndex)
			color = drawing.Color{171, 213, 255, 255}

			series = append(series, chart.TimeSeries{
				Name: v,
				Style: chart.Style{
					StrokeColor: color,
					FillColor:   color.WithAlpha(alpha),
				},
				YAxis:   0,
				XValues: *xValues,
				YValues: *yValues,
			})
			break
		case "isolamento_domiciliare":
			xValues, yValues, xNames, err = nationToTimeseries(data, v, nationIndex)
			color = drawing.Color{171, 213, 255, 255}

			series = append(series, chart.TimeSeries{
				Name: v,
				Style: chart.Style{
					StrokeColor: color,
					FillColor:   color.WithAlpha(alpha),
				},
				YAxis:   0,
				XValues: *xValues,
				YValues: *yValues,
			})
			break
		case "attualmente_positivi":
			xValues, yValues, xNames, err = nationToTimeseries(data, v, nationIndex)
			color = drawing.Color{237, 164, 17, 255}

			series = append(series, chart.TimeSeries{
				Name: v,
				Style: chart.Style{
					StrokeColor: color,
					FillColor:   color.WithAlpha(alpha),
				},
				YAxis:   0,
				XValues: *xValues,
				YValues: *yValues,
			})
			break
		case "nuovi_positivi":
			xValues, yValues, xNames, err = nationToTimeseries(data, v, nationIndex)
			color = drawing.Color{18, 4, 217, 255}

			series = append(series, chart.TimeSeries{
				Name: v,
				Style: chart.Style{
					StrokeColor: color,
					FillColor:   color.WithAlpha(alpha),
				},
				YAxis:   0,
				XValues: *xValues,
				YValues: *yValues,
			})
			break
		case "dimessi_guariti":
			xValues, yValues, xNames, err = nationToTimeseries(data, v, nationIndex)
			color = drawing.Color{38, 224, 175, 255}

			series = append(series, chart.TimeSeries{
				Name: v,
				Style: chart.Style{
					StrokeColor: color,
					FillColor:   color.WithAlpha(alpha),
				},
				YAxis:   0,
				XValues: *xValues,
				YValues: *yValues,
			})
			break
		case "deceduti":
			xValues, yValues, xNames, err = nationToTimeseries(data, v, nationIndex)
			color = chart.ColorAlternateGray

			series = append(series, chart.TimeSeries{
				Name: v,
				Style: chart.Style{
					StrokeColor: color,
					FillColor:   color.WithAlpha(alpha),
				},
				YAxis:   0,
				XValues: *xValues,
				YValues: *yValues,
			})
			break
		case "totale_casi":
			xValues, yValues, xNames, err = nationToTimeseries(data, v, nationIndex)
			if err != nil {
				fmt.Println(err)
			}
			color = drawing.Color{R: 255, A: 255}

			series = append(series, chart.TimeSeries{
				Name: v,
				Style: chart.Style{
					StrokeColor: color,
					FillColor:   color.WithAlpha(alpha),
				},
				YAxis:   0,
				XValues: *xValues,
				YValues: *yValues,
			})
			break
		case "tamponi":
			xValues, yValues, xNames, err = nationToTimeseries(data, v, nationIndex)
			color = drawing.Color{175, 232, 169, 255}

			series = append(series, chart.TimeSeries{
				Name: v,
				Style: chart.Style{
					StrokeColor: color,
					FillColor:   color.WithAlpha(alpha),
				},
				YAxis:   0,
				XValues: *xValues,
				YValues: *yValues,
			})
			break
		default:
			return fmt.Errorf("wrong field name passed"), ""
		}

		annotations := make([]chart.AnnotationSeries, 0)

		err, fileName = timeseriesChart(&series, xNames, &annotations, title, filename, xAxisName, yAxisName)
		if err != nil {
			return fmt.Errorf("%v", err), ""
		}
	}

	return nil, fileName
}

// Returns a plot including national total cases
func TotalePositiviNazione(data *[]NationData, placeAnnotations bool, title, filename string) (error, string) {
	xAxisName := ""
	yAxisName := "Contagiati"

	fieldName := "Totale_casi"
	xTotale, yTotale, xNames, err := nationToTimeseries(data, "Totale_casi", 0)
	if err != nil {
		return fmt.Errorf("error while creating %v chart: %v", fieldName, err), ""
	}

	totale := chart.TimeSeries{
		Name: "Totale contagi",
		Style: chart.Style{
			StrokeColor: drawing.Color{
				R: 255,
				A: 255,
			},
			FillColor: drawing.Color{
				R: 255,
				A: 150,
			},
		},
		YAxis:   0,
		XValues: *xTotale,
		YValues: *yTotale,
	}

	series := make([]chart.TimeSeries, 1)
	series[0] = totale

	annotations := make([]chart.AnnotationSeries, 0)
	if placeAnnotations {
		deltas, e := calculateDeltaPerDay(data, fieldName, 0)
		if e != nil {
			return fmt.Errorf("%v", e), ""
		}
		annotations = append(annotations, deltaAnnotations(deltas, xTotale, yTotale))
	}

	err, fileName := timeseriesChart(&series, xNames, &annotations, title, filename, xAxisName, yAxisName)
	if err != nil {
		return fmt.Errorf("%v", err), ""
	}
	return nil, fileName
}

// Returns a plot including national total healed
func TotaleGuaritiNazione(data *[]NationData, placeAnnotations bool, title, filename string) (error, string) {
	xAxisName := ""
	yAxisName := "Guariti"

	fieldName := "Dimessi_guariti"
	xGuariti, yGuariti, xNames, err := nationToTimeseries(data, fieldName, 0)
	if err != nil {
		return fmt.Errorf("error while creating %v chart: %v", fieldName, err), ""
	}

	guariti := chart.TimeSeries{
		Name: "Guariti",
		Style: chart.Style{
			StrokeColor: chart.ColorGreen,
			FillColor:   chart.ColorGreen.WithAlpha(200),
		},
		YAxis:   0,
		XValues: *xGuariti,
		YValues: *yGuariti,
	}

	series := make([]chart.TimeSeries, 1)
	series[0] = guariti

	annotations := make([]chart.AnnotationSeries, 0)
	deltas, e := calculateDeltaPerDay(data, fieldName, 0)
	if e != nil {
		return fmt.Errorf("%v", e), ""
	}
	annotations = append(annotations, deltaAnnotations(deltas, xGuariti, yGuariti))

	err, fileName := timeseriesChart(&series, xNames, &annotations, title, filename, xAxisName, yAxisName)
	if err != nil {
		return fmt.Errorf("%v", err), ""
	}
	return nil, fileName
}

// Returns a plot including national total deaths
func TotaleDecedutiNazione(data *[]NationData, placeAnnotations bool, title, filename string) (error, string) {
	xAxisName := ""
	yAxisName := "Morti"

	fieldName := "Deceduti"
	xDeceduti, yDeceduti, xNames, err := nationToTimeseries(data, fieldName, 0)
	if err != nil {
		return fmt.Errorf("error while creating %v chart: %v", fieldName, err), ""
	}

	deceduti := chart.TimeSeries{
		Name: "Morti",
		Style: chart.Style{
			StrokeColor: chart.ColorAlternateGray,
			FillColor:   chart.ColorAlternateGray.WithAlpha(200),
		},
		YAxis:   0,
		XValues: *xDeceduti,
		YValues: *yDeceduti,
	}

	series := make([]chart.TimeSeries, 1)
	series[0] = deceduti

	annotations := make([]chart.AnnotationSeries, 0)
	if placeAnnotations {
		deltas, e := calculateDeltaPerDay(data, fieldName, 0)
		if e != nil {
			return fmt.Errorf("%v", e), ""
		}
		annotations = append(annotations, deltaAnnotations(deltas, xDeceduti, yDeceduti))
	}

	err, fileName := timeseriesChart(&series, xNames, &annotations, title, filename, xAxisName, yAxisName)
	if err != nil {
		return fmt.Errorf("%v", err), ""
	}
	return nil, fileName
}

// Returns a plot including national current positive cases
func AttualmentePositiviNazione(data *[]NationData, placeAnnotations bool, title, filename string) (error, string) {
	xAxisName := ""
	yAxisName := "Positivi ancora in vita"

	fieldName := "attualmente_positivi"
	xPositivi, yPositivi, xNames, err := nationToTimeseries(data, fieldName, 0)
	if err != nil {
		return fmt.Errorf("error while creating %v chart: %v", fieldName, err), ""
	}

	positivi := chart.TimeSeries{
		Name: "Positivi ancora in vita",
		Style: chart.Style{
			StrokeColor: drawing.Color{237, 164, 17, 255},
			FillColor:   drawing.Color{237, 164, 17, 200},
		},
		YAxis:   0,
		XValues: *xPositivi,
		YValues: *yPositivi,
	}

	series := make([]chart.TimeSeries, 1)
	series[0] = positivi

	annotations := make([]chart.AnnotationSeries, 0)
	if placeAnnotations {
		deltas, e := calculateDeltaPerDay(data, fieldName, 0)
		if e != nil {
			return fmt.Errorf("%v", e), ""
		}
		annotations = append(annotations, deltaAnnotations(deltas, xPositivi, yPositivi))
	}

	err, fileName := timeseriesChart(&series, xNames, &annotations, title, filename, xAxisName, yAxisName)
	if err != nil {
		return fmt.Errorf("%v", err), ""
	}
	return nil, fileName
}

// Returns a plot including national new cases
func NuoviPositiviNazione(data *[]NationData, placeAnnotations bool, title, filename string) (error, string) {
	xAxisName := ""
	yAxisName := "Nuovi positivi"

	fieldName := "Nuovi_positivi"
	xNuoviPositivi, yNuoviPositivi, xNames, err := nationToTimeseries(data, fieldName, 0)
	if err != nil {
		return fmt.Errorf("error while creating %v chart: %v", fieldName, err), ""
	}

	positivi := chart.TimeSeries{
		Name: "Nuovi Positivi",
		Style: chart.Style{
			StrokeColor: drawing.Color{47, 21, 214, 255},
			FillColor:   drawing.Color{47, 21, 214, 200},
		},
		YAxis:   0,
		XValues: *xNuoviPositivi,
		YValues: *yNuoviPositivi,
	}

	series := make([]chart.TimeSeries, 1)
	series[0] = positivi

	annotations := make([]chart.AnnotationSeries, 0)
	if placeAnnotations {
		deltas, e := calculateDeltaPerDay(data, fieldName, 0)
		if e != nil {
			return fmt.Errorf("%v", e), ""
		}
		annotations = append(annotations, deltaAnnotations(deltas, xNuoviPositivi, yNuoviPositivi))
	}

	err, fileName := timeseriesChart(&series, xNames, &annotations, title, filename, xAxisName, yAxisName)
	if err != nil {
		return fmt.Errorf("%v", err), ""
	}
	return nil, fileName
}

// Returns a plot with the data of a specified region according to the specified fields
func VociRegione(data *[]RegionData, fieldName []string, regionIndex int, regionCode int, title, filename string) (error, string) {
	xAxisName := ""
	yAxisName := ""

	var err error
	var xValues *[]time.Time
	var yValues *[]float64
	var xNames *[]chart.GridLine
	series := make([]chart.TimeSeries, 0)
	var color drawing.Color
	var alpha uint8 = 200
	var fileName string
	for _, v := range fieldName {
		switch strings.ToLower(v) {
		case "ricoverati_con_sintomi":
			xValues, yValues, xNames, err = regionToTimeseries(data, v, regionIndex, regionCode)
			color = drawing.Color{38, 224, 175, 255}

			series = append(series, chart.TimeSeries{
				Name: v,
				Style: chart.Style{
					StrokeColor: color,
					FillColor:   color.WithAlpha(alpha),
				},
				YAxis:   0,
				XValues: *xValues,
				YValues: *yValues,
			})
			break
		case "terapia_intensiva":
			xValues, yValues, xNames, err = regionToTimeseries(data, v, regionIndex, regionCode)
			color = drawing.Color{88, 22, 115, 255}

			series = append(series, chart.TimeSeries{
				Name: v,
				Style: chart.Style{
					StrokeColor: color,
					FillColor:   color.WithAlpha(alpha),
				},
				YAxis:   0,
				XValues: *xValues,
				YValues: *yValues,
			})
			break
		case "totale_ospedalizzati":
			xValues, yValues, xNames, err = regionToTimeseries(data, v, regionIndex, regionCode)
			color = drawing.Color{171, 213, 255, 255}

			series = append(series, chart.TimeSeries{
				Name: v,
				Style: chart.Style{
					StrokeColor: color,
					FillColor:   color.WithAlpha(alpha),
				},
				YAxis:   0,
				XValues: *xValues,
				YValues: *yValues,
			})
			break
		case "isolamento_domiciliare":
			xValues, yValues, xNames, err = regionToTimeseries(data, v, regionIndex, regionCode)
			color = drawing.Color{171, 213, 255, 255}

			series = append(series, chart.TimeSeries{
				Name: v,
				Style: chart.Style{
					StrokeColor: color,
					FillColor:   color.WithAlpha(alpha),
				},
				YAxis:   0,
				XValues: *xValues,
				YValues: *yValues,
			})
			break
		case "attualmente_positivi":
			xValues, yValues, xNames, err = regionToTimeseries(data, v, regionIndex, regionCode)
			color = drawing.Color{237, 164, 17, 255}

			series = append(series, chart.TimeSeries{
				Name: v,
				Style: chart.Style{
					StrokeColor: color,
					FillColor:   color.WithAlpha(alpha),
				},
				YAxis:   0,
				XValues: *xValues,
				YValues: *yValues,
			})
			break
		case "nuovi_positivi":
			xValues, yValues, xNames, err = regionToTimeseries(data, v, regionIndex, regionCode)
			color = drawing.Color{18, 4, 217, 255}

			series = append(series, chart.TimeSeries{
				Name: v,
				Style: chart.Style{
					StrokeColor: color,
					FillColor:   color.WithAlpha(alpha),
				},
				YAxis:   0,
				XValues: *xValues,
				YValues: *yValues,
			})
			break
		case "dimessi_guariti":
			xValues, yValues, xNames, err = regionToTimeseries(data, v, regionIndex, regionCode)
			color = drawing.Color{38, 224, 175, 255}

			series = append(series, chart.TimeSeries{
				Name: v,
				Style: chart.Style{
					StrokeColor: color,
					FillColor:   color.WithAlpha(alpha),
				},
				YAxis:   0,
				XValues: *xValues,
				YValues: *yValues,
			})
			break
		case "deceduti":
			xValues, yValues, xNames, err = regionToTimeseries(data, v, regionIndex, regionCode)
			color = chart.ColorAlternateGray

			series = append(series, chart.TimeSeries{
				Name: v,
				Style: chart.Style{
					StrokeColor: color,
					FillColor:   color.WithAlpha(alpha),
				},
				YAxis:   0,
				XValues: *xValues,
				YValues: *yValues,
			})
			break
		case "totale_casi":
			xValues, yValues, xNames, err = regionToTimeseries(data, v, regionIndex, regionCode)
			color = drawing.Color{R: 255, A: 255}

			series = append(series, chart.TimeSeries{
				Name: v,
				Style: chart.Style{
					StrokeColor: color,
					FillColor:   color.WithAlpha(alpha),
				},
				YAxis:   0,
				XValues: *xValues,
				YValues: *yValues,
			})
			break
		case "tamponi":
			xValues, yValues, xNames, err = regionToTimeseries(data, v, regionIndex, regionCode)
			color = drawing.Color{175, 232, 169, 255}

			series = append(series, chart.TimeSeries{
				Name: v,
				Style: chart.Style{
					StrokeColor: color,
					FillColor:   color.WithAlpha(alpha),
				},
				YAxis:   0,
				XValues: *xValues,
				YValues: *yValues,
			})
			break
		default:
			return fmt.Errorf("wrong field name passed"), ""
		}

		annotations := make([]chart.AnnotationSeries, 0)

		err, fileName = timeseriesChart(&series, xNames, &annotations, title, filename, xAxisName, yAxisName)
		if err != nil {
			return fmt.Errorf("%v", err), ""
		}
	}

	return nil, fileName
}


// Returns a plot with the data of a specified province according to the specified fields
func VociProvince(data *[]ProvinceData, fieldName []string, provinceIndex int, provinceCode int, title, filename string) (error, string) {
	xAxisName := ""
	yAxisName := ""

	var err error
	var xValues *[]time.Time
	var yValues *[]float64
	var xNames *[]chart.GridLine
	series := make([]chart.TimeSeries, 0)
	var color drawing.Color
	var alpha uint8 = 200
	var fileName string
	for _, v := range fieldName {
		switch strings.ToLower(v) {
		case "totale_casi":
			xValues, yValues, xNames, err = provinceToTimeseries(data, v, provinceIndex, provinceCode)
			color = drawing.Color{R: 255, A: 255}

			series = append(series, chart.TimeSeries{
				Name: v,
				Style: chart.Style{
					StrokeColor: color,
					FillColor:   color.WithAlpha(alpha),
				},
				YAxis:   0,
				XValues: *xValues,
				YValues: *yValues,
			})
			break
		case "nuovi_positivi":
			xValues, yValues, xNames, err = provinceToTimeseries(data, v, provinceIndex, provinceCode)
			color = drawing.Color{18, 4, 217, 255}

			series = append(series, chart.TimeSeries{
				Name: v,
				Style: chart.Style{
					StrokeColor: color,
					FillColor:   color.WithAlpha(alpha),
				},
				YAxis:   0,
				XValues: *xValues,
				YValues: *yValues,
			})
		default:
			return fmt.Errorf("wrong field name passed"), ""
		}

		annotations := make([]chart.AnnotationSeries, 0)

		err, fileName = timeseriesChart(&series, xNames, &annotations, title, filename, xAxisName, yAxisName)
		if err != nil {
			return fmt.Errorf("%v", err), ""
		}
	}
	return nil, fileName
}

// Returns total cases for the given province
func TotalePositiviProvincia(data *[]ProvinceData, provinceIndex int, provinceCode int, title, filename string) (error, string) {
	xAxisName := ""
	yAxisName := "Contagiati"

	fieldName := "Totale_casi"
	xTotale, yTotale, xNames, err := provinceToTimeseries(data, "Totale_casi", provinceIndex, provinceCode)
	if err != nil {
		return fmt.Errorf("error while creating %v chart: %v", fieldName, err), ""
	}

	totale := chart.TimeSeries{
		Name: "Totale contagi",
		Style: chart.Style{
			StrokeColor: drawing.Color{
				R: 255,
				A: 255,
			},
			FillColor: drawing.Color{
				R: 255,
				A: 150,
			},
		},
		YAxis:   0,
		XValues: *xTotale,
		YValues: *yTotale,
	}

	series := make([]chart.TimeSeries, 1)
	series[0] = totale

	annotations := make([]chart.AnnotationSeries, 0)

	err, fileName := timeseriesChart(&series, xNames, &annotations, title, filename, xAxisName, yAxisName)
	if err != nil {
		return fmt.Errorf("%v", err), ""
	}
	return nil, fileName
}

// Returns new cases for the given province
func NuoviPositiviProvincia(data *[]ProvinceData, provinceName string, placeAnnotations bool, title, filename string) (error, string) {
	xAxisName := ""
	yAxisName := "Nuovi positivi"

	provinceIndex, err := FindFirstOccurrenceProvince(data, "denominazione_provincia", provinceName)
	if err != nil {
		return err, ""
	}
	fieldName := "NuoviCasi"
	xNuoviPositivi, yNuoviPositivi, xNames, err := provinceToTimeseries(data, fieldName, 0, provinceIndex)
	if err != nil {
		return fmt.Errorf("error while creating %v chart: %v", fieldName, err), ""
	}

	positivi := chart.TimeSeries{
		Name: "Nuovi Positivi",
		Style: chart.Style{
			StrokeColor: drawing.Color{47, 21, 214, 255},
			FillColor:   drawing.Color{47, 21, 214, 200},
		},
		YAxis:   0,
		XValues: *xNuoviPositivi,
		YValues: *yNuoviPositivi,
	}

	series := make([]chart.TimeSeries, 1)
	series[0] = positivi

	annotations := make([]chart.AnnotationSeries, 0)
	if placeAnnotations {
		deltas := intToDeltasArray(data, provinceIndex)
		annotations = append(annotations, deltaAnnotations(deltas, xNuoviPositivi, yNuoviPositivi))
	}

	err, fileName := timeseriesChart(&series, xNames, &annotations, title, filename, xAxisName, yAxisName)
	if err != nil {
		return fmt.Errorf("%v", err), ""
	}
	return nil, fileName
}

// Checks if a given plot already exist by title
func IsGraphExisting(filename string) bool{
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

// Returns a well formatted filename
func FilenameCreator(plotTitle string) (filename string) {
	var colorMode string

	lightHour, _ := time.Parse("15:04:05", daySwitch)
	darkHour, _ := time.Parse("15:04:05", nightSwitch)
	now := time.Now().Hour()
	if now >= darkHour.Hour() || now < lightHour.Hour() {
		colorMode = "dark"
	} else {
		colorMode = "light"
	}

	filename = plotTitle + "-" + colorMode + ".png"
	return
}
