FROM scratch
COPY bin/si-i2p-plugin-static /opt/si-i2p-plugin-static
WORKDIR /opt
COPY si-i2p-plugin/bash /bin/sh
COPY si-i2p-plugin/busybox /bin/busybox
COPY si-i2p-plugin/passwd /etc/passwd
RUN /bin/busybox chown -R sii2pplugindocker /opt
RUN /bin/busybox rm /bin/busybox
USER sii2pplugindocker
ENTRYPOINT ./si-i2p-plugin-static

