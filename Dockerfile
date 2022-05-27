FROM --platform=linux/arm/v6 debian:bullseye AS lib_builder

RUN apt-get update -y && apt-get install -y build-essential cmake git

WORKDIR /lib
RUN git clone https://github.com/jgarff/rpi_ws281x.git
RUN mkdir rpi_ws281x/build/

WORKDIR /lib/rpi_ws281x/build/
RUN cmake -D BUILD_SHARED=OFF -D BUILD_TEST=OFF ..
RUN cmake --build .
RUN make install



FROM golang:1.18.2-bullseye AS app_builder

RUN dpkg --add-architecture armel
RUN apt-get update -y && apt-get install -y crossbuild-essential-armel libgtk-3-dev:armel libwebkit2gtk-4.0-dev:armel

COPY --from=lib_builder /usr/local/lib/ /usr/local/lib/
COPY --from=lib_builder /usr/local/include/ws2811 /usr/local/include/ws2811

ENV GOOS=linux GOARCH=arm GOARM=6
ENV CGO_ENABLED=1 CC=arm-linux-gnueabi-gcc CXX=arm-linux-gnueabi-g++
ENV PKG_CONFIG_PATH="/usr/lib/arm-linux-gnueabi/pkgconfig:${PKG_CONFIG_PATH}"

VOLUME /app

WORKDIR /app
ENTRYPOINT /app/build.sh