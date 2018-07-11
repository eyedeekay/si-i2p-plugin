FROM alpine:3.6
RUN apk update -U
RUN apk add go go-tools go-doc git make musl-dev musl musl-utils
RUN adduser -s /opt/bin/si-i2p-plugin -S -H -D sii2pplugin sii2pplugin
COPY . /opt/
WORKDIR /opt/
RUN make deps
RUN make release && chmod a+x /opt/bin/si-i2p-plugin
RUN apk del go go-tools go-doc git make musl-dev musl musl-utils
USER sii2pplugin
CMD /opt/bin/si-i2p-plugin \
    -proxy-addr=0.0.0.0 \
    -proxy-port=44443 \
    -ah-addr=sam-jumphelper \
    -ah-port=7054 \
    -bridge-addr=sam-host \
    -bridge-port=7656 \
    -in-tunnels=15 \
    -out-tunnels=15 \
    -in-backups=5 \
    -out-backups=5 \
    -internal-ah=false \
    -verbose=true \
    -addressbook=/opt/etc/si-i2p-plugin/addresses.csv
