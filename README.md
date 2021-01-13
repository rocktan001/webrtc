<h1 align="center">
  <br>
  webrtc desktop
  <br>
</h1>
<h4 align="center">Go implementation of the <a href="https://developer.mozilla.org/en-US/docs/Web/API/MediaDevices">MediaDevices</a> API</h4>
<p align="center">
  <a href="https://pion.ly/slack"><img src="https://img.shields.io/badge/join-us%20on%20slack-gray.svg?longCache=true&logo=slack&colorB=brightgreen" alt="Slack Widget"></a>
  <a href="https://github.com/pion/mediadevices/actions"><img src="https://github.com/pion/mediadevices/workflows/CI/badge.svg?branch=master" alt="Build status"></a> 
  <a href="https://pkg.go.dev/github.com/pion/mediadevices"><img src="https://godoc.org/github.com/pion/mediadevices?status.svg" alt="GoDoc"></a>
  <a href="https://codecov.io/gh/pion/mediadevices"><img src="https://codecov.io/gh/pion/mediadevices/branch/master/graph/badge.svg" alt="Coverage Status"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT"></a>
</p>
<br>

`webrtc` 参考pion ，增加本地桌面音视频功能，目前提供编码有h264 ,opus。持续更新中.... 联系：116072620@qq.com

## Install

`go get -u github.com/rocktan001/webrtc@v1.1.0`

## Usage

The following snippet shows how to capture a camera stream and store a frame as a jpeg image:

```go
package main

import (
	"fmt"
	"github.com/rocktan001/webrtc/mediadevices"
	"time"

	"github.com/pion/webrtc/v3"
	"github.com/rocktan001/webrtc/mediadevices/pkg/codec/opus" // This is required to use opus audio encoder
	"github.com/rocktan001/webrtc/mediadevices/pkg/codec/x264" // This is required to use h264 video encoder
	"github.com/rocktan001/webrtc/mediadevices/pkg/prop"

	"github.com/rocktan001/goutil"
	_ "github.com/rocktan001/webrtc/mediadevices/pkg/driver/loopback" // This is required to register microphone adapter
	_ "github.com/rocktan001/webrtc/mediadevices/pkg/driver/screen"
	"github.com/rocktan001/webrtc/signal"
)

func main() {
	go Bye()

	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:www.rocktan001.com:3478"},
			},
			{
				URLs:           []string{"turn:www.rocktan001.com:3478"},
				Username:       "rocktan001",
				Credential:     "F96AEB124C",
				CredentialType: webrtc.ICECredentialTypePassword,
			},
		},
	}

	// Create a new RTCPeerConnection
	x264Params, err := x264.NewParams()
	if err != nil {
		panic(err)
	}
	x264Params.BitRate = 1_000_000 // 500kbps
	// x264Params.Preset = x264.PresetUltrafast
	x264Params.Preset = x264.PresetMedium
	opusParams, err := opus.NewParams()
	if err != nil {
		panic(err)
	}
	codecSelector := mediadevices.NewCodecSelector(
		mediadevices.WithVideoEncoders(&x264Params),
		mediadevices.WithAudioEncoders(&opusParams),
	)

	mediaEngine := webrtc.MediaEngine{}
	codecSelector.Populate(&mediaEngine)
	api := webrtc.NewAPI(webrtc.WithMediaEngine(&mediaEngine))
	peerConnection, err := api.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())
		if connectionState.String() == "disconnected" {
			panic(nil)
		}
		if connectionState.String() == "failed" {
			panic(nil)
		}

	})

	func() {

		s, err := mediadevices.GetDisplayMedia(mediadevices.MediaStreamConstraints{
			Video: func(c *mediadevices.MediaTrackConstraints) {
				// c.FrameFormat = prop.FrameFormat(frame.FormatI420)
				// c.Width = prop.Int(640)
				// c.Height = prop.Int(480)
			},
			Audio: func(c *mediadevices.MediaTrackConstraints) {
				c.SampleRate = prop.Int(48000)
				c.IsBigEndian = prop.BoolExact(false)
				c.ChannelCount = prop.Int(2)
				c.IsFloat = prop.BoolExact(false)
				c.SampleSize = prop.Int(2)
				c.IsInterleaved = prop.BoolExact(true)
				c.Latency = prop.DurationExact(time.Millisecond * 20)
			},
			Codec: codecSelector,
		})
		if err != nil {
			panic(err)
		}
		for _, track := range s.GetTracks() {
			track.OnEnded(func(err error) {
				fmt.Printf("Track (ID: %s) ended with error: %v\n",
					track.ID(), err)
			})
			fmt.Println(track.ID(), "  ", track.StreamID())
			_, err = peerConnection.AddTransceiverFromTrack(track,
				webrtc.RtpTransceiverInit{
					Direction: webrtc.RTPTransceiverDirectionSendrecv,
				},
			)
			if err != nil {
				panic(err)
			}
		}

	}()

	//=========================================================
	// Create an offer to send to the browser
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(offer)
	if err != nil {
		panic(err)
	}

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete

	// Output the answer in base64 so we can paste it in browser
	// fmt.Println(signal.Encode(*peerConnection.LocalDescription()))
	goutil.Redis_json_set("remoteSessionDescription", signal.Encode(*peerConnection.LocalDescription()))
	// fmt.Println(signal.Encode(*peerConnection.LocalDescription()))
	answer := webrtc.SessionDescription{}
	// signal.Decode(signal.MustReadStdin(), &offer)
	goutil.Redis_json_sub("webrtc-start")
	fmt.Println("webrtc-start")
	signal.Decode(goutil.Redis_json_get("localDescription"), &answer)
	// fmt.Println(answer)
	peerConnection.SetRemoteDescription(answer)

	// Block forever
	select {}

}

func Bye() {
	for {
		goutil.Redis_json_sub("webrtc-bye")
		panic("leave")
	}
}



```

## 调试界面 
http://www.rocktan001.com:63000/

## 运行本地server
	cd ${src}/github.com/rocktan01/webrtc/main
	go build server.go
	go build webrtc-desktop-pion-offer.go
	./server.exe

## License
MIT License - see [LICENSE](LICENSE) for full text
