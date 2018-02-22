#! /bin/sh
i2pd --service --loglevel=error --conf=/etc/i2pd/i2pd.conf tunconf=/etc/i2pd/tunnels.conf --log=/dev/null &
python3 /usr/bin/reflect.py
