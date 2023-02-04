FROM --platform=$BUILDPLATFORM tinygo/tinygo:0.26.0 AS build
COPY --chown=tinygo ./src /src
RUN tinygo build -o /src/tinygo_hello.wasm -target wasm /src/main.go

FROM scratch
COPY --from=build /src/tinygo_hello.wasm /app.wasm
ENTRYPOINT [ "app.wasm" ]
