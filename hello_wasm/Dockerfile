FROM --platform=$BUILDPLATFORM tinygo/tinygo:0.26.0 AS build

USER root
RUN apt-get update && apt-get install -y ca-certificates openssl
ARG cert_location=/usr/local/share/ca-certificates
RUN openssl s_client -showcerts -connect github.com:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/github.crt
RUN openssl s_client -showcerts -connect proxy.golang.org:443 </dev/null 2>/dev/null|openssl x509 -outform PEM >  ${cert_location}/proxy.golang.crt
RUN update-ca-certificates

COPY ./src /src
RUN cd /src && tinygo build -o tinygo_hello.wasm -target wasm main.go

FROM scratch
COPY --from=build /src/tinygo_hello.wasm /app.wasm
ENTRYPOINT [ "app.wasm" ]
