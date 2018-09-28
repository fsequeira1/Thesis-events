# Build windows version
env GOOS=windows GOARCH=amd64 go build -o bin/syslogprobe.exe syslogprobe.go
# Build linux version
env GOOS=linux GOARCH=amd64 go build -o bin/syslogprobe syslogprobe.go
# Build mac version
env GOOS=darwin GOARCH=amd64 go build -o bin/syslogprobe_mac syslogprobe.go
# Build RPi version
#env GOOS=linux GOARCH=arm GOARM=6 go build -o syslogprobe_rpi syslogprobe.go
# Build Alpine docker
env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o gounixsocketlistener syslogprobe.go

