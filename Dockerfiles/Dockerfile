FROM alpine:edge
RUN apk update -U
RUN apk add go go-tools go-doc git make musl-dev musl musl-utils
RUN adduser -s /opt/bin/si-i2p-plugin -S -H -D sii2pplugin sii2pplugin
RUN git clone https://github.com/eyedeekay/si-i2p-plugin /opt/
RUN chown -R sii2pplugin /opt/
WORKDIR /opt/
RUN go get github.com/eyedeekay/gosam
RUN make release
RUN apk del go go-tools go-doc git make musl-dev musl musl-utils
USER sii2pplugin
CMD /opt/bin/si-i2p-plugin -proxy-addr=0.0.0.0 -bridge-addr=thirdeye-host
