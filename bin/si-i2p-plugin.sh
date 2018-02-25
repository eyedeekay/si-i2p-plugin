#! /usr/bin/env sh

set -e

CFGFILE=/etc/si-i2p-plugin/settings.cfg

if [ -f /etc/si-i2p-plugin/settings.cfg ]; then
        . /etc/si-i2p-plugin/settings.cfg
        si2p_config="-bridge-addr=$sam_addr \
            -bridge-port=$sam_port \
            -proxy-addr=$proxy_addr \
            -proxy-port=$proxy_port \
            -http-proxy=$use_proxy \
            -conn-debug=$debug_info \
            -address=$testing_url"
fi

if [ ! -d "$working_dir" ]; then
        echo "Working directory $working_dir does not exist"
        exit 1
fi

cd "$working_dir" && si-i2p-plugin $si2p_config
