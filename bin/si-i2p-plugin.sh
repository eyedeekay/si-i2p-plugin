#! /usr/bin/env sh

set -e

CFGFILE=/etc/si-i2p-plugin/settings.cfg
LOGFILE=/var/log/si-i2p-plugin/si-i2p-plugin.log

if [ -f /etc/si-i2p-plugin/settings.cfg ]; then
        . /etc/si-i2p-plugin/settings.cfg
        si2p_config="-bridge-addr=$sam_addr \
            -bridge-port=$sam_port \
            -proxy-addr=$proxy_addr \
            -proxy-port=$proxy_port \
            -http-proxy=$use_proxy \
            -conn-debug=$debug_info \
            -url=$testing_url \
            -addresshelper=$jump_urls"
fi

if [ ! -d "$working_dir" ]; then
        echo "Working directory $working_dir does not exist"
        exit 1
fi

echo "si-i2p-plugin $si2p_config 1> \"$LOGFILE\" 2> \"$LOGFILE\""

cd "$working_dir" && si-i2p-plugin $si2p_config 1> "$LOGFILE" 2> "$LOGFILE"
