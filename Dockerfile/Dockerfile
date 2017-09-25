FROM golang
RUN adduser --system --no-create-home --disabled-password --disabled-login --group sii2pplugin
COPY . /opt/src
RUN mkdir -p /opt/bin
RUN chown -R sii2pplugin /opt/src
WORKDIR /opt/src
RUN go get github.com/eyedeekay/gosam
RUN go build -o /opt/bin/si-i2p-proxy ./src
CMD /opt/bin/si-i2p-proxy
