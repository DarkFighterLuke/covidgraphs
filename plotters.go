package datagraphs

import (
	"fmt"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
	"time"
)

func andamentoNazionaleCompleto(data *[]nationData) (error, string) {
	title := "Andamento nazionale"
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

	err, filename := timeseriesChart(&series, xNames, &annotations, title, xAxisName, yAxisName)
	if err != nil {
		return fmt.Errorf("%v", err), ""
	}
	return nil, filename
}

func totalePositiviNazione(data *[]nationData, placeAnnotations bool) (error, string) {
	title := "Totale Contagi"
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

	err, filename := timeseriesChart(&series, xNames, &annotations, title, xAxisName, yAxisName)
	if err != nil {
		return fmt.Errorf("%v", err), ""
	}
	return nil, filename
}

func totaleGuaritiNazione(data *[]nationData, placeAnnotations bool) (error, string) {
	title := "Totale Guariti"
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

	err, filename := timeseriesChart(&series, xNames, &annotations, title, xAxisName, yAxisName)
	if err != nil {
		return fmt.Errorf("%v", err), ""
	}
	return nil, filename
}

func totaleDecedutiNazione(data *[]nationData, placeAnnotations bool) (error, string) {
	title := "Totale Morti"
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

	err, filename := timeseriesChart(&series, xNames, &annotations, title, xAxisName, yAxisName)
	if err != nil {
		return fmt.Errorf("%v", err), ""
	}
	return nil, filename
}

func attualmentePositiviNazione(data *[]nationData, placeAnnotations bool) (error, string) {
	title := "Positivi Ancora in Vita"
	xAxisName := ""
	yAxisName := "Positivi ancora in vita"

	fieldName := "Totale_attualmente_positivi"
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

	err, filename := timeseriesChart(&series, xNames, &annotations, title, xAxisName, yAxisName)
	if err != nil {
		return fmt.Errorf("%v", err), ""
	}
	return nil, filename
}

func nuoviPositiviNazione(data *[]nationData, placeAnnotations bool) (error, string) {
	title := "Nuovi Positivi"
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

	err, filename := timeseriesChart(&series, xNames, &annotations, title, xAxisName, yAxisName)
	if err != nil {
		return fmt.Errorf("%v", err), ""
	}
	return nil, filename
}

func vociRegione(data *[]regionData, fieldName []string, regionName string) (error, string) {
	title := "Dati regione " + regionName
	xAxisName := ""
	yAxisName := ""

	regionCode, err := findFirstOccurrenceRegion(data, "denominazione_regione", regionName)
	if err != nil {
		return fmt.Errorf("wrong region name: %v", err), ""
	}

	var xValues *[]time.Time
	var yValues *[]float64
	var xNames *[]chart.GridLine
	series := make([]chart.TimeSeries, 0)
	var color drawing.Color
	var alpha uint8 = 200
	var filename string
	for _, v := range fieldName {
		switch v {
		case "Ricoverati_con_sintomi":
			xValues, yValues, xNames, err = regionToTimeseries(data, v, 0, regionCode)
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
		case "Terapia_intensiva":
			xValues, yValues, xNames, err = regionToTimeseries(data, v, 0, regionCode)
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
		case "Totale_ospedalizzati":
			xValues, yValues, xNames, err = regionToTimeseries(data, v, 0, regionCode)
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
		case "Isolamento_domiciliare":
			xValues, yValues, xNames, err = regionToTimeseries(data, v, 0, regionCode)
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
		case "Totale_attualmente_positivi":
			xValues, yValues, xNames, err = regionToTimeseries(data, v, 0, regionCode)
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
		case "Nuovi_positivi":
			xValues, yValues, xNames, err = regionToTimeseries(data, v, 0, regionCode)
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
		case "Dimessi_guariti":
			xValues, yValues, xNames, err = regionToTimeseries(data, v, 0, regionCode)
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
		case "Deceduti":
			xValues, yValues, xNames, err = regionToTimeseries(data, v, 0, regionCode)
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
		case "Totale_casi":
			xValues, yValues, xNames, err = regionToTimeseries(data, v, 0, regionCode)
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
		case "Tamponi":
			xValues, yValues, xNames, err = regionToTimeseries(data, v, 0, regionCode)
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

		err, filename = timeseriesChart(&series, xNames, &annotations, title, xAxisName, yAxisName)
		if err != nil {
			return fmt.Errorf("%v", err), ""
		}
	}

	return nil, filename
}

func totalePositiviProvincia(data *[]provinceData, provinceName string) (error, string) {
	title := "Totale Contagi " + provinceName
	xAxisName := ""
	yAxisName := "Contagiati"

	provinceIndex, err := findFirstOccurrenceProvince(data, "denominazione_provincia", provinceName)
	if err != nil {
		return err, ""
	}
	fieldName := "Totale_casi"
	xTotale, yTotale, xNames, err := provinceToTimeseries(data, "Totale_casi", 0, provinceIndex)
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

	err, filename := timeseriesChart(&series, xNames, &annotations, title, xAxisName, yAxisName)
	if err != nil {
		return fmt.Errorf("%v", err), ""
	}
	return nil, filename
}

func nuoviPositiviProvincia(data *[]provinceData, provinceName string, placeAnnotations bool) (error, string) {
	title := "Nuovi Positivi " + provinceName
	xAxisName := ""
	yAxisName := "Nuovi positivi"

	provinceIndex, err := findFirstOccurrenceProvince(data, "denominazione_provincia", provinceName)
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

	err, filename := timeseriesChart(&series, xNames, &annotations, title, xAxisName, yAxisName)
	if err != nil {
		return fmt.Errorf("%v", err), ""
	}
	return nil, filename
}
