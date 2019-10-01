# querxnagios v1.2
Interact with Querx Smart Ethernet Sensors using Go (golang)

## compile
```
git clone https://github.com/egnite/querxnagios.git
cd querxnagios/check_querx
go get github.com/egnite/querx
go get github.com/egnite/querxnagios
GOARCH=386 go build
cp check_querx linux-32-bit/
GOARCH=amd64 go build
cp check_querx linux-64-bit/
tar cvzf egnite-querx_nagios_plugin-1.2.tar.gz linux-32-bit/ linux-64-bit/
cp check_querx /var/lib/nagios/plugins
```

## More information on Querx Smart Ethernet sensors
Find more information on Querx on the [product page](http://sensors.egnite.de)
and on [egnite's website](http://www.egnite.de)

## Tutorial for ICINGA2
A in-depth tutorial on server room monitoring with Icinga2 can be found
on on the product page as well. Please be patient while we translate it
from German to English.
