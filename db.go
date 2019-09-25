package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

/*
create table if not exists location_thermostats(
	id smallserial NOT null primary key,
	sensor_id varchar(50) NOT NULL UNIQUE,
	location_name varchar(50) NOT NULL
);

create table if not exists public.thermostats_data (
	created_at timestamp NOT null primary key DEFAULT now(),
	location_id int not null references location(id),
	temperature numeric NOT NULL,
	desired_temperature numeric NOT null,
	battery numeric NOT null
);
*/

type Thermostat struct {
	Id       int    `db:"id"`
	Address  string `db:"sensor_id"`
	Location string `db:"location_name"`
}

type ThermostatData struct {
	LocationId  int     `db:"location_id"`
	Temp        float32 `db:"temperature"`
	DesiredTemp float32 `db:"desired_temperature"`
	Battery     int     `db:"battery"`
}

type DbHandler struct {
	db *sqlx.DB
}

func GetDb(connection string) *DbHandler {
	return &DbHandler{
		db: sqlx.MustConnect("postgres", connection),
	}
}

func (db *DbHandler) GetThermostats() ([]Thermostat, error) {
	var th []Thermostat
	err := db.db.Select(&th, "SELECT id, sensor_id, location_name from location_thermostats;")
	return th, err
}

func (db *DbHandler) InsertMeasurement(data ThermostatData) {
	_, err := db.db.NamedExec(`INSERT into thermostats_data (location_id, temperature, desired_temperature, battery) VALUES (:location_id, :temperature, :desired_temperature, :battery)`, data)
	if err != nil{
		log.Println("error writing thermostat data to db", err)
	}
}
