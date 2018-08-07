FROM debian:stretch-slim

ARG V8_VERSION=6.9.253

RUN apt-get update && apt-get upgrade -yqq

RUN DEBIAN_FRONTEND=noninteractive \
    apt-get install bison \
                    cdbs \
                    curl \
                    flex \
                    g++ \
                    git \
                    python \
                    pkg-config -yqq

RUN git clone https://chromium.googlesource.com/chromium/tools/depot_tools.git

ENV PATH="/depot_tools:${PATH}"

RUN fetch v8 && \
    cd /v8 && \
    git checkout ${V8_VERSION} && \
    ./tools/dev/v8gen.py x64.release && \
    ninja -C out.gn/x64.release
 
# WHY: v8 only works when its working dir has debug crap. Otherwise it gives a nice
# Illegal instruction (core dumped)
RUN mkdir -p /usr/local/bin && \
	echo "#!/bin/sh\n cd /v8/out.gn/x64.release && ./d8 \$@" > /usr/local/bin/d8 && \
	chmod +x /usr/local/bin/d8
	
ENV GOPATH="/go"
ENV GOROOT="/goroot"
ENV GOVERSION="1.10"
ENV PATH="${GOROOT}/bin:${GOPATH}/bin:${PATH}"

RUN apt-get update && apt-get install -y wget make

RUN cd /tmp && \
	wget https://storage.googleapis.com/golang/go$GOVERSION.linux-amd64.tar.gz && \
	tar -xvf go$GOVERSION.linux-amd64.tar.gz && \
	mkdir -p $GOROOT && \
	mv ./go/* $GOROOT
