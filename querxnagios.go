package querxnagios

import (
	"fmt"
	"os"

	"github.com/pborman/getopt/v2"
)

const Version string = "1.3"

type Parameters struct {
	Hostname *string
	Port     *int
	//TLS             *bool
	//Auth            *bool
	//Username        *string
	//Password        *string
	Verbosity       *int
	Help            *bool
	Version         *bool
	Critical        *string
	Warning         *string
	Messages        []string
	Failed          bool
	UseDeviceLimits bool
	WarningGiven    bool
	SensorID        *int
}

func (p *Parameters) Parse() {
	p.Hostname = getopt.StringLong("hostname", 'H', "", "Hostname or IP address", "ip")
	p.Port = getopt.IntLong("port", 'P', 80, "HTTP port", "port")
	//p.TLS = getopt.BoolLong("tls", 't', "Connect via TLS")
	//p.Auth = getopt.BoolLong("auth", 'a', "Use Authentification")
	//p.Username = getopt.StringLong("user", 'u', "", "User name")
	//p.Password = getopt.StringLong("password", 'p', "", "Password")
	p.Verbosity = getopt.CounterLong("verbose", 'v', "Verbosity level")
	p.Help = getopt.BoolLong("help", 'h', "Print help screen")
	p.Version = getopt.BoolLong("version", 'V', "Print version")
	p.Critical = getopt.StringLong("critical", 'c', "", "Critical threshold", "range")
	p.Warning = getopt.StringLong("warning", 'w', "", "Warning threshold", "range")
	p.Messages = make([]string, 4)
	p.SensorID = getopt.IntLong("sensor", 's', 0, "Sensor ID", "id")
	p.Failed = false
	getopt.Parse()

	//Check, if a host name was given
	if !getopt.IsSet('H') {
		p.Messages = append(p.Messages, "Please provide a host name")
		p.Failed = true
	}

	if !getopt.IsSet('c') {
		p.Messages = append(p.Messages, "No critical values were given, querying device status")
		p.UseDeviceLimits = true
	}

	//Check, if a Port was given
	if !getopt.IsSet('P') {
		p.Messages = append(p.Messages, "No port given, default Port will be used")
	}

	if !getopt.IsSet('s') {
		p.Messages = append(p.Messages, "No sensor ID was given, using sensor ID 0")
	}

	if getopt.IsSet('w') {
		p.WarningGiven = true
	}

	//If started with --auth  we also need a username and a password
	/*if getopt.IsSet('a') {
		if !getopt.IsSet('u') || !getopt.IsSet('p') {
			p.Messages = append(p.Messages, "When using --auth, you need to provide a user name and a password")
			p.Failed = true
		}
	}
	*/

	if *p.Version {
		fmt.Println("Version: " + Version)
		os.Exit(0)
	}

	if *p.Help || p.Failed {
		getopt.PrintUsage(os.Stderr)
		for _, message := range p.Messages {
			fmt.Fprintln(os.Stderr, message)
		}
		os.Exit(1)
	}
}
