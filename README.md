# Docker + Wasm + TinyGo example

## Usage
1. [Turn on the Docker+Wasm integration](https://docs.docker.com/desktop/wasm/#turn-on-the-dockerwasm-integration)
2. `git clone https://github.com/tsukanov-as/docker_wasm_tinygo_hello.git`
3. `cd docker_wasm_tinygo_hello`
4. `docker-compose up --build`

## Reading list
0. https://hacks.mozilla.org/2019/03/standardizing-wasi-a-webassembly-system-interface/
1. https://wazero.io/specs/#wasi
2. https://wazero.io/languages/tinygo/
3. https://tinygo.org/docs/reference/lang-support/stdlib/
4. https://github.com/tinygo-org/tinygo/issues/2704
5. https://github.com/golang/go/issues/58141
6. https://github.com/docker/roadmap/issues/426
7. [WasmEdge WASI](https://github.com/WasmEdge/WasmEdge/blob/4f1059b4aa934c4de3c3edff4055c28e6f0b9311/lib/host/wasi/wasimodule.cpp)
