package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	fmt.Println("Hello there")

	//----------------------------------------------------------------------------------

	var ns uint64
	clock_time_get(0, 1000, &ns) // direct wasi call example (see: https://github.com/WasmEdge/WasmEdge/blob/4f1059b4aa934c4de3c3edff4055c28e6f0b9311/lib/host/wasi/wasimodule.cpp)

	fmt.Println("Time from WASI: ", time.Unix(0, int64(ns)))

	//----------------------------------------------------------------------------------

	buf := make([]byte, 10)
	random_get(&buf[0], uint32(len(buf)))

	fmt.Println("Random from WASI: ", buf)

	//----------------------------------------------------------------------------------

	var argsCnt uint32
	var argsLen uint32
	args_sizes_get(&argsCnt, &argsLen)

	fmt.Println("Args sizes: ", argsCnt, argsLen)

	//----------------------------------------------------------------------------------

	buf = make([]byte, argsLen)
	p := make([]*byte, argsCnt)
	args_get(&p[0], &buf[0])

	fmt.Println("Args from WASI: ", string(*p[0]), strings.Split(string(buf), "\x00"))
	fmt.Println("Args OS: ", os.Args)

	//----------------------------------------------------------------------------------

	var list [2]__wasi_iovec_t
	buf = []byte("hello from wasi 1\n")
	list[0] = __wasi_iovec_t{
		buf:    &buf[0],
		bufLen: uint32(len(buf)),
	}
	buf = []byte("hello from wasi 2\n")
	list[1] = __wasi_iovec_t{
		buf:    &buf[0],
		bufLen: uint32(len(buf)),
	}
	var n uint32
	fd_write(stdout, &list[0], uint32(len(list)), &n)
	println("written: ", n)

	//----------------------------------------------------------------------------------

	subs := []__wasi_subscription_t_clock{
		{
			__wasi_subscription_t: __wasi_subscription_t{
				userData: 1,
				tag:      __wasi_eventtype_t_clock,
			},
			id:        0,
			timeout:   5 * 1_000_000_000, // 5 sec
			precision: 1000,
			flags:     0,
		},
		{
			__wasi_subscription_t: __wasi_subscription_t{
				userData: 2,
				tag:      __wasi_eventtype_t_clock,
			},
			id:        0,
			timeout:   6 * 1_000_000_000, // 6 sec
			precision: 1000,
			flags:     0,
		},
	}

	events := make([]__wasi_event_t, len(subs))
	var nevents uint32

	println("sleep 5 sec...")
	if err := poll_oneoff(&subs[0].__wasi_subscription_t, &events[0], uint32(len(subs)), &nevents /* Out */); err != 0 {
		println("poll_oneoff error: ", err)
	}

	fmt.Printf("Events from WASI: %+v\n", events[:nevents])

	//----------------------------------------------------------------------------------

	var sock int32

	if err := sock_open(__WASI_ADDRESS_FAMILY_INET4, __WASI_SOCK_TYPE_SOCK_STREAM, &sock); err != 0 {
		println("sock_open error: ", err)
	}

	buf = []byte{0, 0, 0, 0}
	addr := __wasi_address_t{
		buf:     &buf[0],
		buf_len: uint32(len(buf)),
	}

	const port = 8080

	if err := sock_bind(sock, &addr, port); err != 0 {
		println("sock_bind error: ", err)
	}

	if err := sock_listen(sock, 10); err != 0 {
		println("sock_listen error: ", err)
	}

	fmt.Printf("listen localhost:%d\n", port)

	{
		println("wait connection...")
		var conn int32
		if err := sock_accept(sock, &conn); err != 0 {
			println("sock_accept error: ", err)
			proc_exit(0)
		}
		println("Accepted new connection")
		var list [1]__wasi_iovec_t
		buf = make([]byte, 256)
		list[0] = __wasi_iovec_t{
			buf:    &buf[0],
			bufLen: uint32(len(buf)),
		}
		var n, f uint32
		if err := sock_recv(conn, &list[0], uint32(len(list)), 0, &n, &f); err != 0 {
			println("fd_read error: ", err)
			proc_exit(0)
		}
		fmt.Printf("Client send %d bytes:\n%s", n, string(buf))
	}

	proc_exit(0)
	println("will not be printed")
}

// ArgsGet uint32_t body(uint32_t ArgvPtr, uint32_t ArgvBufPtr);

//go:wasm-module wasi_snapshot_preview1
//export args_get
func args_get(ArgvPtr **byte, ArgvBufPtr *byte) uint32

// ArgsSizesGet uint32_t body(uint32_t /* Out */ ArgcPtr, uint32_t /* Out */ ArgvBufSizePtr);

//go:wasm-module wasi_snapshot_preview1
//export args_sizes_get
func args_sizes_get(ArgcPtr *uint32 /* Out */, ArgvBufPtr *uint32 /* Out */) uint32

// EnvironGet uint32_t body(uint32_t EnvPtr, uint32_t EnvBufPtr);
// EnvironSizesGet uint32_t body(uint32_t /* Out */ EnvCntPtr, uint32_t /* Out */ EnvBufSizePtr);
// ClockResGet uint32_t body(uint32_t ClockId, uint32_t /* Out */ ResolutionPtr);

// ClockTimeGet uint32_t body(uint32_t ClockId, uint64_t Precision, uint32_t /* Out */ TimePtr);

//go:wasm-module wasi_snapshot_preview1
//export clock_time_get
func clock_time_get(ClockId uint32, Precision uint64, TimePtr *uint64 /* Out */) uint32

// fd_advise uint32_t body(int32_t Fd, uint64_t Offset, uint64_t Len, uint32_t Advice);
// fd_allocate uint32_t body(int32_t Fd, uint64_t Offset, uint64_t Len);
// fd_close uint32_t body(int32_t Fd);
// fd_datasync uint32_t body(int32_t Fd);
// fd_fdstat_get uint32_t body(int32_t Fd, uint32_t /* Out */ FdStatPtr);
// fd_fdstat_set_flags uint32_t body(int32_t Fd, uint32_t FsFlags);
// fd_fdstat_set_rights uint32_t body(int32_t Fd, uint64_t FsRightsBase, uint64_t FsRightsInheriting);
// fd_filestat_get uint32_t body(int32_t Fd, uint32_t /* Out */ FilestatPtr);
// fd_filestat_set_size uint32_t body(int32_t Fd, uint64_t Size);
// fd_filestat_set_times uint32_t body(int32_t Fd, uint64_t ATim, uint64_t MTim, uint32_t FstFlags);
// fd_pread uint32_t body(int32_t Fd, uint32_t IOVsPtr, uint32_t IOVsLen, uint64_t Offset, uint32_t /* Out */ NReadPtr);
// fd_prestat_get uint32_t body(int32_t Fd, uint32_t /* Out */ PreStatPtr);
// fd_prestat_dir_name uint32_t body(int32_t Fd, uint32_t PathBufPtr, uint32_t PathLen);
// fd_pwrite uint32_t body(int32_t Fd, uint32_t IOVSPtr, uint32_t IOVSLen, uint64_t Offset, uint32_t /* Out */ NWrittenPtr);
// fd_read uint32_t body(int32_t Fd, uint32_t IOVSPtr, uint32_t IOVSLen, uint32_t /* Out */ NReadPtr);

//go:wasm-module wasi_snapshot_preview1
//export fd_read
func fd_read(Fd int32, IOVSPtr *__wasi_iovec_t, IOVSLen uint32, NReadPtr *uint32) uint32

// fd_read_dir uint32_t body(int32_t Fd, uint32_t BufPtr, uint32_t BufLen, uint64_t Cookie, uint32_t /* Out */ NReadPtr);
// fd_renumber uint32_t body(int32_t Fd, int32_t ToFd);
// fd_seek int32_t body(int32_t Fd, int64_t Offset, uint32_t Whence, uint32_t /* Out */ NewOffsetPtr);
// fd_sync uint32_t body(int32_t Fd);
// fd_tell uint32_t body(int32_t Fd, uint32_t /* Out */ OffsetPtr);

// fd_write uint32_t body(int32_t Fd, uint32_t IOVSPtr, uint32_t IOVSLen, uint32_t /* Out */ NWrittenPtr);

const stdout = 1

type __wasi_size_t = uint32

type __wasi_iovec_t struct {
	buf    *byte
	bufLen __wasi_size_t
}

//go:wasm-module wasi_snapshot_preview1
//export fd_write
func fd_write(Fd int32, IOVSPtr *__wasi_iovec_t, IOVSLen uint32, NWrittenPtr *uint32) uint32

// path_create_directory uint32_t body(int32_t Fd, uint32_t PathPtr, uint32_t PathLen);
// path_filestat_get uint32_t body(int32_t Fd, uint32_t Flags, uint32_t PathPtr, uint32_t PathLen, uint32_t /* Out */ FilestatPtr);
// path_filestat_set_times uint32_t body(int32_t Fd, uint32_t Flags, uint32_t PathPtr, uint32_t PathLen, uint64_t ATim, uint64_t MTim, uint32_t FstFlags);
// path_link uint32_t body(int32_t OldFd, uint32_t OldFlags, uint32_t OldPathPtr, uint32_t OldPathLen, int32_t NewFd, uint32_t NewPathPtr, uint32_t NewPathLen);
// path_open uint32_t body(int32_t DirFd, uint32_t DirFlags, uint32_t PathPtr, uint32_t PathLen, uint32_t OFlags, uint64_t FsRightsBase, uint64_t FsRightsInheriting, uint32_t FsFlags, uint32_t /* Out */ FdPtr);
// path_read_link uint32_t body(int32_t Fd, uint32_t PathPtr, uint32_t PathLen, uint32_t BufPtr, uint32_t BufLen, uint32_t /* Out */ NReadPtr);
// path_remove_directory uint32_t body(int32_t Fd, uint32_t PathPtr, uint32_t PathLen);
// path_rename uint32_t body(int32_t Fd, uint32_t OldPathPtr, uint32_t OldPathLen, int32_t NewFd, uint32_t NewPathPtr, uint32_t NewPathLen);
// path_symlink uint32_t body(uint32_t OldPathPtr, uint32_t OldPathLen, int32_t Fd, uint32_t NewPathPtr, uint32_t NewPathLen);
// path_unlink_file uint32_t body(int32_t Fd, uint32_t PathPtr, uint32_t PathLen);

// poll_oneoff uint32_t body(uint32_t InPtr, uint32_t OutPtr, uint32_t NSubscriptions, uint32_t /* Out */ NEventsPtr);

const (
	__wasi_eventtype_t_clock    __wasi_eventtype_t = 0
	__wasi_eventtype_t_fd_read  __wasi_eventtype_t = 1
	__wasi_eventtype_t_fd_write __wasi_eventtype_t = 2
)

type __wasi_eventtype_t = byte

type __wasi_subscription_t struct {
	userData uint64
	tag      __wasi_eventtype_t
}

type __wasi_subscription_t_clock struct {
	__wasi_subscription_t
	id        uint32
	timeout   uint64
	precision uint64
	flags     uint16
}

type __wasi_subscription_t_fd_readwrite struct {
	__wasi_subscription_t
	fd int32
}

type __wasi_event_t struct {
	userData  uint64
	errno     uint16
	eventType __wasi_eventtype_t

	// only used for fd_read or fd_write events
	_ struct {
		nBytes uint64
		flags  uint16
	}
}

//go:wasm-module wasi_snapshot_preview1
//export poll_oneoff
func poll_oneoff(InPtr *__wasi_subscription_t, OutPtr *__wasi_event_t, NSubscriptions uint32, NEventsPtr *uint32 /* Out */) uint32

// epoll_oneoff uint32_t body(uint32_t InPtr, uint32_t OutPtr, uint32_t NSubscriptions, uint32_t /* Out */ NEventsPtr);

// proc_exit Expect<void> body(uint32_t Status);

//go:wasm-module wasi_snapshot_preview1
//export proc_exit
func proc_exit(Status uint32)

// proc_raise uint32_t body(uint32_t Signal);
// sched_yield uint32_t body(const Runtime::CallingFrame &Frame);

// random_get uint32_t body(uint32_t BufPtr, uint32_t BufLen);

//go:wasm-module wasi_snapshot_preview1
//export random_get
func random_get(buf *byte, len uint32) (errno uint32)

// sock_open uint32_t body(uint32_t AddressFamily, uint32_t SockType, uint32_t /* Out */ RoFdPtr);

const (
	__WASI_ADDRESS_FAMILY_UNSPEC = 0
	__WASI_ADDRESS_FAMILY_INET4  = 1
	__WASI_ADDRESS_FAMILY_INET6  = 2
)

const (
	__WASI_SOCK_TYPE_SOCK_ANY    = 0
	__WASI_SOCK_TYPE_SOCK_DGRAM  = 1
	__WASI_SOCK_TYPE_SOCK_STREAM = 2
)

//go:wasm-module wasi_snapshot_preview1
//export sock_open
func sock_open(AddressFamily uint32, SockType uint32, RoFdPtr *int32) (errno uint32)

// sock_bind uint32_t body(int32_t Fd, uint32_t AddressPtr, uint32_t Port);

type __wasi_address_t struct {
	buf     *byte
	buf_len __wasi_size_t
}

// type __wasi_sockaddr_in_t struct {
// 	sin_family   uint16
// 	sin_port     uint16
// 	sin_addr     __wasi_address_t
// 	sin_zero_len __wasi_size_t
// 	sin_zero     *byte
// }

//go:wasm-module wasi_snapshot_preview1
//export sock_bind
func sock_bind(Fd int32, AddressPtr *__wasi_address_t, Port uint32) (errno uint32)

// sock_listen uint32_t body(int32_t Fd, int32_t Backlog);

//go:wasm-module wasi_snapshot_preview1
//export sock_listen
func sock_listen(Fd int32, Backlog int32) (errno uint32)

// sock_accept uint32_t body(int32_t Fd, uint32_t /* Out */ RoFdPtr);

//go:wasm-module wasi_snapshot_preview1
//export sock_accept
func sock_accept(Fd int32, RoFdPtr *int32 /* Out */) (errno uint32)

// sock_connect uint32_t body(int32_t Fd, uint32_t AddressPtr, uint32_t Port);
// sock_recv uint32_t body(int32_t Fd, uint32_t RiDataPtr, uint32_t RiDataLen, uint32_t RiFlags, uint32_t /* Out */ RoDataLenPtr, uint32_t /* Out */ RoFlagsPtr);

//go:wasm-module wasi_snapshot_preview1
//export sock_recv
func sock_recv(Fd int32, RiDataPtr *__wasi_iovec_t, RiDataLen uint32, RiFlags uint32, RoDataLenPtr *uint32 /* Out */, RoFlagsPtr *uint32 /* Out */) uint32

// sock_recvfrom uint32_t body(int32_t Fd, uint32_t RiDataPtr, uint32_t RiDataLen, uint32_t AddressPtr, uint32_t RiFlags, uint32_t /* Out */ PortPtr, uint32_t /* Out */ RoDataLenPtr, uint32_t /* Out */ RoFlagsPtr);
// sock_send uint32_t body(int32_t Fd, uint32_t SiDataPtr, uint32_t SiDataLen, uint32_t SiFlags, uint32_t /* Out */ SoDataLenPtr);
// sock_sendto uint32_t body(int32_t Fd, uint32_t SiDataPtr, uint32_t SiDataLen, uint32_t AddressPtr, int32_t Port, uint32_t SiFlags, uint32_t /* Out */ SoDataLenPtr);
// sock_shutdown uint32_t body(int32_t Fd, uint32_t SdFlags);

const (
	__WASI_SDFLAGS_RD = 1
	__WASI_SDFLAGS_WR = 2
)

//go:wasm-module wasi_snapshot_preview1
//export sock_shutdown
func sock_shutdown(Fd int32, SdFlags uint32) (errno uint32)

// sock_getopt uint32_t body(int32_t Fd, uint32_t SockOptLevel, uint32_t SockOptName, uint32_t FlagPtr, uint32_t FlagSizePtr);
// sock_setopt uint32_t body(int32_t Fd, uint32_t SockOptLevel, uint32_t SockOptName, uint32_t FlagPtr, uint32_t FlagSizePtr);
// sock_get_local_addr uint32_t body(int32_t Fd, uint32_t AddressPtr, uint32_t PortPtr);
// sock_get_peeraddr uint32_t body(int32_t Fd, uint32_t AddressPtr, uint32_t PortPtr);

// get_addr_info uint32_t body(uint32_t NodePtr, uint32_t NodeLen, uint32_t ServicePtr, uint32_t ServiceLen, uint32_t HintsPtr, uint32_t ResPtr, uint32_t MaxResLength, uint32_t ResLengthPtr);
