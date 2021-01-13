// +build !dynamic

package vpx

//#cgo CFLAGS: -I${SRCDIR}/include
//#cgo CXXFLAGS: -I${SRCDIR}/include
//#cgo linux,arm LDFLAGS: ${SRCDIR}/lib/libvpx_linux_armv7.a -lm
//#cgo linux,arm64 LDFLAGS: ${SRCDIR}/lib/libvpx_linux_arm64.a -lm
//#cgo linux,amd64 LDFLAGS: ${SRCDIR}/lib/libvpx_linux_x64.a -lm
//#cgo darwin,amd64 LDFLAGS: ${SRCDIR}/lib/libvpx_darwin_x64.a
//#cgo windows,amd64 LDFLAGS: ${SRCDIR}/lib/libvpx_windows_x64.a
import "C"
