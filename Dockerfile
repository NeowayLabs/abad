FROM ubuntu:18.04

ENV DEPOTPATH=/depot	
ENV PATH=${PATH}:${DEPOTPATH}

RUN apt-get update -y  && \
	apt-get install -y wget python git pkg-config && \
	git clone https://chromium.googlesource.com/chromium/tools/depot_tools.git	${DEPOTPATH} && \
	cd /tmp && fetch v8

WORKDIR /tmp/v8

RUN gclient sync

# Google V8 Deps copied from Google install-build-deps.sh script
# Why copied ? because the script wont work =), I lost a day with it.
# Why not fix properly and help them ? Because fuck bloated software and I dont want
# to be part of it, just want to compile V8.
COPY hack/install-build-deps.sh /hack/v8/install-build-deps.sh

RUN /hack/v8/install-build-deps.sh
# Uncomment to run the default script
# RUN ./build/install-build-deps.sh --no-prompt

# RUN ./tools/dev/gm.py x64.release

# WORKDIR /abad


	