module github.com/panzerdev/bt-thermostat

go 1.13

require (
	github.com/go-ble/ble v0.0.0-20190521171521-147700f13610
	github.com/jmoiron/sqlx v1.2.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/lib/pq v1.2.0
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.1.0
	github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4 // indirect
	github.com/prometheus/common v0.7.0 // indirect
	github.com/prometheus/procfs v0.0.5 // indirect
	google.golang.org/appengine v1.6.3 // indirect
)

replace github.com/go-ble/ble => github.com/panzerdev/ble v0.0.0-20190924180509-f484c7857a7a
