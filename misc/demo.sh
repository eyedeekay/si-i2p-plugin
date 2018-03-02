#! /bin/sh

chown -R i2pd:i2pd /var/lib/i2pd
ln -sf /usr/share/i2pd/certificates /var/lib/i2pd/certificates
ln -sf /etc/i2pd/subscriptions.txt /var/lib/i2pd/subscriptions.txt
su - -c "i2pd i2pd --service --loglevel=info
    --conf=/etc/i2pd/i2pd.conf
    --tunconf=/etc/i2pd/tunnels.conf
    --log=/var/log/i2pd/log" &
sleep 5

#i2pd --service --loglevel=error --conf=/etc/i2pd/i2pd.conf tunconf=/etc/i2pd/tunnels.conf --log=/dev/null &
python3 /usr/bin/reflect.py

tail -f /var/log/i2pd/log
