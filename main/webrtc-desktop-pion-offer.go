package main

import (
	// "context"
	"fmt"
	"github.com/rocktan001/webrtc/mediadevices"
	"time"
	// "github.com/rocktan001/webrtc/mediadevices/examples/internal/signal"
	// "github.com/rocktan001/webrtc/mediadevices/pkg/frame"
	"github.com/rocktan001/webrtc/mediadevices/pkg/prop"
	"github.com/pion/webrtc/v3"

	// If you don't like x264, you can also use vpx by importing as below
	// "github.com/rocktan001/webrtc/mediadevices/pkg/codec/vpx" // This is required to use VP8/VP9 video encoder
	// or you can also use openh264 for alternative h264 implementation
	// "github.com/rocktan001/webrtc/mediadevices/pkg/codec/openh264"
	// or if you use a raspberry pi like, you can use mmal for using its hardware encoder
	// "github.com/rocktan001/webrtc/mediadevices/pkg/codec/mmal"
	"github.com/rocktan001/webrtc/mediadevices/pkg/codec/opus" // This is required to use opus audio encoder
	"github.com/rocktan001/webrtc/mediadevices/pkg/codec/x264" // This is required to use h264 video encoder

	// Note: If you don't have a camera or microphone or your adapters are not supported,
	//       you can always swap your adapters with our dummy adapters below.
	// _ "github.com/rocktan001/webrtc/mediadevices/pkg/driver/videotest"
	// _ "github.com/rocktan001/webrtc/mediadevices/pkg/driver/audiotest"
	// _ "github.com/rocktan001/webrtc/mediadevices/pkg/driver/camera"     // This is required to register camera adapter
	// _ "github.com/rocktan001/webrtc/mediadevices/pkg/driver/microphone" // This is required to register microphone adapter
	_ "github.com/rocktan001/webrtc/mediadevices/pkg/driver/loopback" // This is required to register microphone adapter
	_ "github.com/rocktan001/webrtc/mediadevices/pkg/driver/screen"
	"github.com/rocktan001/goutil"
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
	// config := webrtc.Configuration{
	// 	ICEServers: []webrtc.ICEServer{
	// 		{
	// 			URLs: []string{"stun:www.rocktan001.com:3478"},
	// 		},
	// 		// {
	// 		// 	URLs:           []string{"turn:hk01.cn.menghaocheng.com:3478"},
	// 		// 	Username:       "menghaocheng",
	// 		// 	Credential:     "mypasswd",
	// 		// 	CredentialType: webrtc.ICECredentialTypePassword,
	// 		// },
	// 		{
	// 			URLs:           []string{"turn:gz01.cn.coturn.menghaocheng.com:3478"},
	// 			Username:       "menghaocheng",
	// 			Credential:     "mypasswd",
	// 			CredentialType: webrtc.ICECredentialTypePassword,
	// 		},
	// 		{
	// 			URLs:           []string{"turn:sz01.cn.coturn.menghaocheng.com:63078"},
	// 			Username:       "menghaocheng",
	// 			Credential:     "mypasswd",
	// 			CredentialType: webrtc.ICECredentialTypePassword,
	// 		},
	// 		// {
	// 		// 	URLs:           []string{"turn:hk01.cn.menghaocheng.com:3478"},
	// 		// 	Username:       "menghaocheng",
	// 		// 	Credential:     "mypasswd",
	// 		// 	CredentialType: webrtc.ICECredentialTypePassword,
	// 		// },
	// 		// {
	// 		// 	URLs:           []string{"turn:la01.us.coturn.menghaocheng.com:3478"},
	// 		// 	Username:       "menghaocheng",
	// 		// 	Credential:     "mypasswd",
	// 		// 	CredentialType: webrtc.ICECredentialTypePassword,
	// 		// },
	// 	},
	// }
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

	// _, iceConnectedCtxCancel := context.WithCancel(context.Background())
	// iceConnectedCtx, iceConnectedCtxCancel := context.WithCancel(context.Background())
	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())
		if connectionState.String() == "disconnected" {
			panic(nil)
		}
		if connectionState.String() == "failed" {
			panic(nil)
		}
		// if connectionState == webrtc.ICEConnectionStateConnected {
		// 	// iceConnectedCtxCancel()
		// }
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
