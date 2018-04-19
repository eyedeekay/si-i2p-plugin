#! /usr/bin/env sh

CFGFILE=/etc/si-i2p-plugin/settings.cfg
if [ -f /etc/si-i2p-plugin/settings.cfg ]; then
    . /etc/si-i2p-plugin/settings.cfg
    ps aux | grep -v 'grep' | grep $(cat $PIDFILE)
fi
