//Package dash implements the Manifest interface to encode/parse MPEG DASH MPD files.
//
//
//Usage:
//  import (
// 	  "bufio"
// 	  "bytes"
// 	  "io/ioutil"
// 	  "os"
// 	  "time"
//
// 	  "stash.redspace.com/ing/manifest/dash"
//  )
//
//  func main() {
//    //initiate a new MPD with profile and minBufferTime
// 	  mpd := dash.NewMPD("urn:mpeg:dash:profile:isoff-live:2011", time.Second*2)
// 	  period := &dash.Period{
// 		AdaptationSets: dash.AdaptationSets{
// 			&dash.AdaptationSet{SegmentAlignment: true,
// 				MaxWidth:     1280,
// 				MaxHeight:    720,
// 				MaxFrameRate: "24",
// 				Representations: dash.Representations{&dash.Representation{
// 					ID: "1", MimeType: "video/mp4", Codecs: "avc1.4d01f", Width: 1280, Height: 720,
// 					FrameRate: "24", StartWithSAP: 1, Bandwidth: 980104,
// 					SegmentTemplate: &dash.SegmentTemplate{
// 						Timescale: 12288,
// 						Duration:  24576,
// 						Media:     "video_$Number$.mp4"}},
// 				}},
// 			&dash.AdaptationSet{SegmentAlignment: true,
// 				Representations: dash.Representations{&dash.Representation{
// 					ID: "1", MimeType: "audio/mp4", Codecs: "mp4a.40.29", AudioSamplingRate: "48000",
// 					StartWithSAP: 1, Bandwidth: 33434,
// 					AudioChannelConfig: []*dash.Descriptor{
// 						&dash.Descriptor{SchemeIDURI: "audio_channel_configuration:2011", Value: "2"}},
// 					SegmentTemplate: &dash.SegmentTemplate{
// 						Timescale:          48000,
// 						Duration:           94175,
// 						Media:              "audio_$Number$.mp4",
// 						InitializationAttr: "BBB_32k_init.mp4"}},
// 				}},
// 		}}
// 	  mpd.Periods = append(mpd.Periods, period)
//
// 	  reader, err := mpd.Encode()
// 	  if err != nil {
// 		 panic(err)
// 	  }
//
// 	  buf := new(bytes.Buffer)
// 	  buf.ReadFrom(reader)
//
// 	  if err := ioutil.WriteFile("./output.mpd", buf.Bytes(), 0666); err != nil {
// 		 panic(err)
// 	  }
//
// 	  f , err := os.Open("./output.mpd")
// 	  if err != nil {
// 		 panic(err)
// 	  }
// 	  defer f.Close()
//
// 	  newMPD := &dash.MPD{}
// 	  if err := newMPD.Parse(bufio.NewReader(f)); err != nil {
// 		 panic(err)
// 	  }
//    //manipulate playlist. Ex: add encryption
// 	  	cp := dash.NewContentProtection("mp4.urn.test", "cenc", "", "", "")
//    	cp2 := dash.NewContentProtection("mp4.urn.test", "cenc", "1234", "psshhashstring", "")
//    	cp3 := dash.NewContentProtection("playready.uri", "MSPR 2.0", "", "", "msprhashstring")
//      cp4 := dash.NewContentProtection("playready.uri", "MSPR 2.0", "", "", "")
//    	cp4.SetTrackEncryptionBox(8, "16bytekeyidentifier")
//    	newMPD.Periods[0].AdaptationSets[0].CENCContentProtections =
//    		append(newMPD.Periods[0].AdaptationSets[0].CENCContentProtections, cp, cp2, cp3, cp4)
//
// 	  newReader, err := newMPD.Encode()
// 	  if err != nil {
// 		 panic(err)
// 	  }
//
// 	  newBuf := new(bytes.Buffer)
// 	  newBuf.ReadFrom(newReader)
//
// 	  if err := ioutil.WriteFile("./newOutput.mpd", newBuf.Bytes(), 0666); err != nil {
// 		 panic(err)
// 	  }
//  }
package dash
