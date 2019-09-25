package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/panzerdev/bt-thermostat/ble"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type EnvConfig struct {
	DbConnectionString string `required:"true" envconfig:"DB_STRING"`
	Port               string `required:"true" envconfig:"PORT"`
}

func main() {
	var env EnvConfig
	envconfig.MustProcess("", &env)

	promInit()
	dbHandler := GetDb(env.DbConnectionString)
	go startCollecting(dbHandler)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatalln(http.ListenAndServe(":"+env.Port, nil))
}

func startCollecting(handler *DbHandler) {
	for {
		thermostats, err := handler.GetThermostats()
		if err != nil {
			time.Sleep(time.Second * 15)
			continue
		}
		log.Printf("Time to measure %v thermostats\n", len(thermostats))
		for _, th := range thermostats {
			log.Printf("Start reading %v - MAC: %v", th.Location, th.Address)
			tryReadThermostat(th, handler)
			log.Printf("End reading %v - MAC: %v", th.Location, th.Address)
		}
		log.Println("Measuring done")
		time.Sleep(time.Minute * 5)
	}
}

func tryReadThermostat(th Thermostat, handler *DbHandler) {
	lb := th.Labels()
	for tries := 1; tries < 4; tries++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		tempV, bv, err := readThermostat(th, ctx)
		cancel()

		if err != nil {
			lb[CODE] = "500"
			cnt.With(lb).Inc()
			continue
		}
		setGauge(Battery, th, float32(*bv))
		setGauge(CurrentTemp, th, tempV.CurrentTemp)
		setGauge(DesiredTemp, th, tempV.DesiredTemp)

		lb[CODE] = "200"
		cnt.With(lb).Inc()
		handler.InsertMeasurement(ThermostatData{
			LocationId:  th.Id,
			Temp:        tempV.CurrentTemp,
			DesiredTemp: tempV.DesiredTemp,
			Battery:     int(*bv),
		})
		return
	}
	currentTemp.With(th.Labels()).Set(0)
	log.Printf("Cound not read  %v - %v \n", th.Location, th.Address)
}

func readThermostat(th Thermostat, ctx context.Context) (*TempValue, *BatteryValue, error) {
	bleAdaptor := ble.NewClientAdaptor(th.Address)
	err := bleAdaptor.Connect(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("connect failed: %w", err)
	}

	bs := PinValue(0).ToBytes()
	err = bleAdaptor.WriteCharacteristic(UUID_PIN, bs)
	if err != nil {
		return nil, nil, fmt.Errorf("write UUID_PIN failed: %w", err)
	}

	defer func() {
		bleAdaptor.Disconnect()
		bleAdaptor.Finalize()
	}()

	tempV, err := ReadTemperatures(bleAdaptor)
	if err != nil {
		return nil, nil, fmt.Errorf("read temperatures %w", err)
	}

	bv, err := ReadBattery(bleAdaptor)
	if err != nil {
		return nil, nil, fmt.Errorf("read battery %w", err)
	}

	return tempV, bv, nil
}

func ReadTemperatures(bt ble.BLEConnector) (*TempValue, error) {
	c, err := bt.ReadCharacteristic(UUID_TEMP)
	if err != nil {
		return nil, err
	}

	tempV, err := TempValueFromBytes(c)
	if err != nil {
		return nil, err
	}
	return &tempV, nil
}

func ReadBattery(bt ble.BLEConnector) (*BatteryValue, error) {
	c, err := bt.ReadCharacteristic(UUID_BATTERY)
	if err != nil {
		return nil, err
	}

	res := BatteryValue(c[0])
	return &res, nil
}
