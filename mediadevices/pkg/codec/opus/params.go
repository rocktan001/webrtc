package opus

import (
	"github.com/rocktan001/webrtc/mediadevices/pkg/codec"
	"github.com/rocktan001/webrtc/mediadevices/pkg/io/audio"
	"github.com/rocktan001/webrtc/mediadevices/pkg/prop"
	"github.com/rocktan001/webrtc/mediadevices/pkg/wave/mixer"
)

// Params stores opus specific encoding parameters.
type Params struct {
	codec.BaseParams
	// ChannelMixer is a mixer to be used if number of given and expected channels differ.
	ChannelMixer mixer.ChannelMixer
}

// NewParams returns default opus codec specific parameters.
func NewParams() (Params, error) {
	return Params{}, nil
}

// RTPCodec represents the codec metadata
func (p *Params) RTPCodec() *codec.RTPCodec {
	return codec.NewRTPOpusCodec(48000)
}

// BuildAudioEncoder builds opus encoder with given params
func (p *Params) BuildAudioEncoder(r audio.Reader, property prop.Media) (codec.ReadCloser, error) {
	return newEncoder(r, property, *p)
}
