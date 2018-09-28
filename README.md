# syslogProbe
Probe that consumes events from Syslog and sends them to Kafka

# Demo instructions

Compile:

1 - Golang installed with environment properly configured

2 - Rsyslog properly configured (sending snort events to unix-socket /tmp/snort)

3 - Snort installed with rules configured and sending alerts to rsyslog

4 - Execute syslogcompiler.sh

5 - chmod+x syslogExec.sh

Execute syslogprobe:

1 - Start rsyslog (in case it isn't running)
		service rsyslog start

2 - Begin capture (change interface accordingly):
		sudo snort -D -i eno16777984 -u snort -g snort -c /etc/snort/snort.conf

3 - Start Probe (DEMO machine --> /root/projects/src/src/github.com/vgraveto/syslogProbe)
	./syslogExec.sh


DON'T FORGET TO COPY SYSLOGPROBE TO /USR/BIN

(IN CASE EVERYTHING FAILS REBOOT AND REPEAT STEPS Execute "syslogprobe", IF THE PROBLEM REMAINS LAY OVER AND CRY)
