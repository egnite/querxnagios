package main

import (
	"fmt"
	"strconv"
	"math"
	"time"

	"github.com/egnite/querx"
	"github.com/egnite/querxnagios"
	"gopkg.in/jabdr/monitoringplugin.v1"
)

type querxCheck struct {
	params querxnagios.Parameters
	hostname string
	sensor querx.Sensor
	value float64
	warnRange monitoringplugin.Range
	critRange monitoringplugin.Range
}

func (check *querxCheck) HandleArguments(options monitoringplugin.PluginOpt) (monitoringplugin.PluginOpt, error) {
	var err error

	//Parse command line arguments
	params := querxnagios.Parameters{}
	params.Parse()
	check.params = params

	check.warnRange, err = monitoringplugin.NewRange(*params.Warning)
	if err != nil {
		return options, err
	}

	check.critRange, err = monitoringplugin.NewRange(*params.Critical)
	if  err != nil {
		return options, err
	}
	//fmt.Printf("WARN: %v  CRIT: %v\n", check.warnRange, check.critRange)

	//Initialize Querx and connect
	querx := querx.NewQuerx(*params.Hostname, *params.Port, false)
	err = querx.QueryCurrent()
	if err != nil {
		return options, fmt.Errorf("Could not establish connection to "+*params.Hostname)
	}

	//Get parameters for check
	sensor, err := querx.SensorByID(*params.SensorID)
	if err != nil {
		return options, fmt.Errorf("Failed to query sensor "+strconv.Itoa(*params.SensorID))
	}
	check.sensor = sensor
	check.value, err = querx.CurrentValue(sensor)
	check.hostname = querx.Current.Hostname
	if err != nil {
		return options, fmt.Errorf("Could not query current Readings from sensor "+strconv.Itoa(*params.SensorID))
	}

	options.Timeout = time.Duration(60) * time.Second
	options.PerformanceDataSpec = []monitoringplugin.PerformanceDataSpec{
		{
			Label:             sensor.Name,
			UnitOfMeasurement: monitoringplugin.NumberUnitSpecification,
			Minimum:	   math.Inf(-1),
			Maximum:	   math.Inf(1),
			Warning:           &check.warnRange,
			Critical: 	   &check.critRange,
		},
	}
	options.Check = check

	return options, nil
}

func (check *querxCheck) Run() (result monitoringplugin.CheckResult) {
	var cUpper, cLower, wUpper, wLower float64
	var message string
	var wAlertOnInside = false
	var cAlertOnInside = false

	//fmt.Printf("Run() %v\n", check)

	checkResult := monitoringplugin.NewDefaultCheckResult(nil)
	result = checkResult

	sensor := check.sensor
	params := check.params
	//fmt.Printf("Parameters %v\n", params)



	//Check critical values
	if params.UseDeviceLimits {
		//no critical values were provided
		cUpper = sensor.UpperLimit
		cLower = sensor.LowerLimit
	} else {
		cUpper = check.critRange.End
		cLower = check.critRange.Start
		cAlertOnInside = check.critRange.Invert
	}

	if params.WarningGiven {
		wLower = check.warnRange.Start
		wUpper = check.warnRange.End
		wAlertOnInside = check.warnRange.Invert
	}
	//Check warning values
	tUnit := monitoringplugin.NumberUnit(check.value)
	checkResult.SetPerformanceData(sensor.Name, tUnit)

	s := monitoringplugin.OK

	//Perform check for warning
	if params.WarningGiven {
		wOutOfRange := check.value > wUpper || check.value < wLower
		if wOutOfRange || (wAlertOnInside && !wOutOfRange) {
			s = monitoringplugin.WARNING
		}
	}

	//Perform check for critical
	cOutOfRange := check.value > cUpper || check.value < cLower
	if cOutOfRange || (cAlertOnInside && !cOutOfRange) {
		s = monitoringplugin.CRITICAL
	}

	switch s {
	case monitoringplugin.OK:
		message = "OK"
		break;
	case monitoringplugin.WARNING:
		message = "WARNING"
		break;
	case monitoringplugin.CRITICAL:
		message = "CRITICAL"
		break;
	default:
		message = "UNKOWN"
		break;
	}
	message += fmt.Sprintf(": [%s:%s] %3.1f %s", check.hostname, sensor.Name, check.value, sensor.Unit)
	checkResult.SetResult(s, message)

	return
}

func main() {
	plugin := monitoringplugin.NewPlugin(&querxCheck{})
	defer plugin.Exit()
	plugin.Start()
}
