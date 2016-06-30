package main

import (
	"fmt"
	"strconv"

	"github.com/egnite/querx"
	"github.com/egnite/querxnagios"
	"github.com/olorin/nagiosplugin"
)

func main() {
	var cUpper, cLower, wUpper, wLower, value float64
	var wAlertOnInside = false
	var cAlertOnInside = false

	//prepare Nagios plugin
	check := nagiosplugin.NewCheck()
	defer check.Finish()

	//Parse command line arguments
	params := querxnagios.Parameters{}
	params.Parse()

	//Initialize Querx and connect
	querx := querx.NewQuerx(*params.Hostname, *params.Port, false)
	err := querx.QueryCurrent()
	if err != nil {
		check.AddResult(nagiosplugin.UNKNOWN, "Could not establish connection to "+*params.Hostname)
	}

	//Get parameters for check
	sensor, err := querx.SensorByID(*params.SensorID)
	if err != nil {
		check.AddResult(nagiosplugin.UNKNOWN, "Failed to querx sensor "+strconv.Itoa(*params.SensorID))
	}
	value, err = querx.CurrentValue(sensor)
	if err != nil {
		check.AddResult(nagiosplugin.UNKNOWN, "Could not query current Readings from sensor "+strconv.Itoa(*params.SensorID))
	}

	//Check critical values
	if params.UseDeviceLimits {
		//no critical values were provided
		cUpper = sensor.UpperLimit
		cLower = sensor.LowerLimit
	} else {
		//parse critical values
		critical, err := nagiosplugin.ParseRange(*params.Critical)
		if err != nil {
			check.AddResult(nagiosplugin.UNKNOWN, "Failed to parse critical threshold: "+*params.Critical)
		} else {
			cUpper = critical.End
			cLower = critical.Start
			cAlertOnInside = critical.AlertOnInside
		}
	}

	if params.WarningGiven {
		warning, err := nagiosplugin.ParseRange(*params.Warning)
		if err != nil {
			check.AddResult(nagiosplugin.UNKNOWN, "Failed to parse warning threshold: "+*params.Warning)
		} else {
			wLower = warning.Start
			wUpper = warning.End
			wAlertOnInside = warning.AlertOnInside
		}
	}
	//Check warning values
	check.AddPerfDatum(sensor.Name, "", value)
	//Perform check for critical
	cOutOfRange := value > cUpper || value < cLower
	message := fmt.Sprintf("[%s:%s] %3.2f %s", querx.Current.Hostname, sensor.Name, value, sensor.Unit)
	if cOutOfRange || (cAlertOnInside && !cOutOfRange) {
		check.AddResult(nagiosplugin.CRITICAL, message)
	}

	//Perform check for warning

	if params.WarningGiven {
		wOutOfRange := value > wUpper || value < wLower
		if wOutOfRange || (wAlertOnInside && !wOutOfRange) {
			check.AddResult(nagiosplugin.WARNING, message)
		}
	}
	//add standard result
	check.AddResult(nagiosplugin.OK, message)
}
