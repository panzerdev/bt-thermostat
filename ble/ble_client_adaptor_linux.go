package ble

import (
	blelib "github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
)

func defaultDevice() (d blelib.Device, err error) {
	return linux.NewDevice()
}
