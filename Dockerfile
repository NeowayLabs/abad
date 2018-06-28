FROM ubuntu:18.04

ENV DEPOTPATH=/depot	
ENV PATH=${PATH}:${DEPOTPATH}

RUN apt-get update -y  && \
	apt-get install -y wget python git lsb-release sudo && \
	git clone https://chromium.googlesource.com/chromium/tools/depot_tools.git	${DEPOTPATH} && \
	cd /tmp && fetch v8

WORKDIR /tmp/v8

# RUN ./build/install-build-deps.sh


	