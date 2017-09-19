#! /usr/bin/env sh

set -e

CFGFILE=/etc/si-i2p-plugin/settings.cfg

if [ -f /etc/si-i2p-plugin/settings.cfg ]; then
        . /etc/si-i2p-plugin/settings.cfg
        si2p_config="$sam_addr $sam_port $proxy_addr $proxy_port $debug_info $testing_url"
fi

if [ ! -d "$working_dir" ]; then
        echo "Working directory $working_dir does not exist"
        exit 1
fi

cd "$working_dir" && si-i2p-plugin $si2p_config
