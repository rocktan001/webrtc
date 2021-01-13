// +build !dynamic

package x264

//#cgo CFLAGS: -I${SRCDIR}/include
//#cgo CXXFLAGS: -I${SRCDIR}/include
//#cgo linux,arm LDFLAGS: ${SRCDIR}/lib/libx264_linux_armv7.a -lm
//#cgo linux,arm64 LDFLAGS: ${SRCDIR}/lib/libx264_linux_arm64.a -lm
//#cgo linux,amd64 LDFLAGS: ${SRCDIR}/lib/libx264_linux_x64.a -lm
//#cgo darwin,amd64 LDFLAGS: ${SRCDIR}/lib/libx264_darwin_x64.a
//#cgo windows,amd64 LDFLAGS: ${SRCDIR}/lib/libx264__windows_x64.a
import "C"
