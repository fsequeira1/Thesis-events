#!/bin/bash
clear
echo "stop snort"
service snortd stop
echo "rsyslog restart"
service rsyslog restart
echo "start snort"
service snortd start
echo "removing unix socket"
rm -i -f /tmp/snort
echo "starting syslogProbe"
syslogprobe &
