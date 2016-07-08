FROM debian:jessie

RUN DEBIAN_FRONTEND=noninteractive apt-get update -y \
    && DEBIAN_FRONTEND=noninteractive apt-get -yy -q \
 install \
 iptables \
 ca-certificates \
 file \
 util-linux \
 net-tools \
 procps \
 iproute2 \
 iputils-tracepath \
 traceroute \
 && DEBIAN_FRONTEND=noninteractive apt-get autoremove -y

COPY sigmactler /usr/local/bin/sigmactler
RUN chmod a+rx  /usr/local/bin/sigmactler
COPY kubectl /usr/local/bin/kubectl
RUN chmod a+rx  /usr/local/bin/kubectl

