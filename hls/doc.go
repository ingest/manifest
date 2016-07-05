//Package hls implements the Manifest interface of package m3u8 to encode/parse
//playlists used in HTTP Live Streaming. Comments explaining type attributes are
//related to HLS Spec 'MUST' and 'MUST NOT' recommendations, and should be considered
//when creating your MediaPlaylist and MasterPlaylist objects for encoding.
//
//Example usage:
//
//Encoding Manifest
//  import "stash.redspace.com/ing/m3u8/hls"
//
//  func main(){
//    //Will start a MediaPlaylist object for hls version 7
//    p := hls.NewMediaPlaylist(7)
//    p.TargetDuration = 10
//    p.EndList = true
//    segment := &hls.Segment{
//      URI:       "segmenturi.ts",
//      Inf:       &Inf{Duration: 9.052},
//      Byterange: &Byterange{Length: 400},
//    }
//    p.Segments = append(p.Segments, segment)
//    reader, err := p.Encode()
//    if err!=nil{
//    //handle error
//    }
//    if err := ioutil.WriteFile("path/to/file", p.Bytes(), 0666); err != nil {
//    //handle error
//    }
//  }
//
//
//
//Decoding Manifest
//  import "stash.redspace.com/ing/m3u8/hls"
//
//  func main(){
//    f, err := os.Open("path/to/file.m3u8")
//    if err != nil {
//      //handle error
//    }
//    playlist := &hls.MasterPlaylist{}
//    if err = playlist.Parse(bufio.NewReader(f)); err!=io.EOF{
//      //handle error
//    }
//    //manipulate playlist
//  }
//
package hls
