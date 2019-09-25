package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

const TimeFormat = "04-15-02-01-06"

type TimeValue struct {
	time.Time
}

func (t TimeValue) ToBytes() []byte {
	b := make([]byte, 5)
	timeString := t.Format(TimeFormat)
	timeArray := strings.Split(timeString, "-")
	for i, v := range timeArray {
		value, err := strconv.Atoi(v)
		if err != nil {
			log.Println(err)
		}
		b[i] = byte(value)
	}
	return b
}

func TimeValueFromBytes(b []byte) (TimeValue, error) {
	if len(b) != 5 {
		return TimeValue{}, fmt.Errorf("Input byte array to short. Expected 5 and got %v", len(b))
	}
	timeString := fmt.Sprintf("%02d-%02d-%02d-%02d-%02d", b[0], b[1], b[2], b[3], b[4])
	t, err := time.Parse(TimeFormat, timeString)
	if err != nil {
		return TimeValue{}, err
	}
	return TimeValue{t}, nil
}

type PinValue uint32

func (pv PinValue) ToBytes() []byte {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(pv))
	return bs
}

type TempValue struct {
	CurrentTemp          float32
	DesiredTemp          float32
	EngergySavingTemp    float32
	HeatingTemp          float32
	OffsetTemp           float32
	WindowDetecThreshold int
	WindowDetecTimer     int
}

func (tv TempValue) String() string {
	buf := &bytes.Buffer{}
	tab := tabwriter.NewWriter(buf, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintf(tab, "CurrentTemp\t%.02f\n", tv.CurrentTemp)
	fmt.Fprintf(tab, "DesiredTemp\t%.02f\n", tv.DesiredTemp)
	fmt.Fprintf(tab, "EngergySavingTemp\t%.02f\n", tv.EngergySavingTemp)
	fmt.Fprintf(tab, "HeatingTemp\t%.02f\n", tv.HeatingTemp)
	fmt.Fprintf(tab, "OffsetTemp\t%.02f\n", tv.OffsetTemp)
	fmt.Fprintf(tab, "WindowDetecThreshold\t%v\n", tv.WindowDetecThreshold)
	fmt.Fprintf(tab, "WindowDetecTimer\t%v\n", tv.WindowDetecTimer)
	tab.Flush()
	return buf.String()
}

func (tv TempValue) ToBytes() []byte {
	b := make([]byte, 7)
	b[0] = tempToBytes(tv.CurrentTemp)
	b[1] = tempToBytes(tv.DesiredTemp)
	b[2] = tempToBytes(tv.EngergySavingTemp)
	b[3] = tempToBytes(tv.HeatingTemp)
	b[4] = tempToBytes(tv.OffsetTemp)
	b[5] = byte(tv.WindowDetecThreshold)
	b[6] = byte(tv.WindowDetecTimer)
	log.Println("out",b)
	return b
}

func TempValueFromBytes(b []byte) (TempValue, error) {
	if len(b) != 7 {
		return TempValue{}, fmt.Errorf("Input byte array to short. Expected 7 and got %v", len(b))
	}
	return TempValue{
		CurrentTemp:          tempFromBytes(b[0]),
		DesiredTemp:          tempFromBytes(b[1]),
		EngergySavingTemp:    tempFromBytes(b[2]),
		HeatingTemp:          tempFromBytes(b[3]),
		OffsetTemp:           tempFromBytes(b[4]),
		WindowDetecThreshold: int(b[5]),
		WindowDetecTimer:     int(b[6]),
	}, nil
}

func tempFromBytes(b byte) float32 {
	return float32(b) / float32(2)
}

func tempToBytes(f float32) byte {
	return byte(f * 2)
}

type BatteryValue byte

func (bv BatteryValue) String() string {
	return fmt.Sprintf("%d %%", byte(bv))
}
