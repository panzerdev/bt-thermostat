package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
)

var (
	currentTemp    *prometheus.GaugeVec
	desiredTemp    *prometheus.GaugeVec
	currentBattery *prometheus.GaugeVec
	cnt            *prometheus.CounterVec
)

const (
	ADDRESS  = "address"
	LOCATION = "location"
	CODE     = "code"
)

func (t Thermostat) Labels() prometheus.Labels {
	return prometheus.Labels{
		ADDRESS:  t.Address,
		LOCATION: t.Location,
	}
}

func promInit() {
	currentTemp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "current_room_temp",
			Help: "The current temp resported by the thermostat",
		},
		[]string{
			// Which MAC address
			ADDRESS,
			// In what room
			LOCATION,
		},
	)
	prometheus.MustRegister(currentTemp)

	desiredTemp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "desired_room_temp",
			Help: "The desired temp resported by the thermostat",
		},
		[]string{
			// Which MAC address
			ADDRESS,
			// In what room
			LOCATION,
		},
	)
	prometheus.MustRegister(desiredTemp)

	currentBattery = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "current_room_battery",
			Help: "The current battery percentage resported by the thermostat",
		},
		[]string{
			// Which MAC address
			ADDRESS,
			// In what room
			LOCATION,
		},
	)
	prometheus.MustRegister(currentBattery)

	cnt = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "thermostat_value_read",
			Help: "Total number of statuscodes for read from thermostat",
		},
		[]string{ADDRESS, LOCATION, CODE},
	)
	prometheus.MustRegister(cnt)
}

type PromGaugeType int

const (
	CurrentTemp PromGaugeType = iota
	DesiredTemp
	Battery
)

func setGauge(target PromGaugeType, t Thermostat, val float32) {
	var gauge *prometheus.GaugeVec
	switch target {
	case CurrentTemp:
		gauge = currentTemp
	case DesiredTemp:
		gauge = desiredTemp
	case Battery:
		gauge = currentBattery
	}
	if gauge == nil{
		log.Println("couldn't fine gauge for PromGaugeType: ",target)
		return
	}

	gauge.With(t.Labels()).Set(float64(val))
}
