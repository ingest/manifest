package hls

//MediaPlaylist represents a Media Playlist and its tags.
//
//TODO:(sliding window) - add field for sliding window to represent either the max amount of segments
//or the max duration of a window (TBD). Also would be useful to add variable to track the current first and last sequence numbers
//as a helper to adding and removing segments and tracking MediaSequence, DiscontinuitySequence etc
//
type MediaPlaylist struct {
	*Variant                          // Variant is embedded, contains information on how the master playlist represented this media playlist.
	Version               int         // Version is required, is written #EXT-X-VERSION: <int>.
	Segments              Segments    // Segments are represented by #EXT-INF\n <duration>.
	TargetDuration        int         // TargetDuration is required, is written #EXT-X-TARGETDURATION: <int>. MUST BE >= largest EXT-INF duration
	MediaSequence         int         //Represents tag #EXT-X-MEDIA-SEQUENCE. Number of the first media sequence in the playlist.
	DiscontinuitySequence int         //Represents tag #EXT-X-DISCONTINUITY-SEQUENCE. If present, MUST appear before the first Media Segment. MUST appear before any EXT-X-DISCONTINUITY Media Segment tag.
	EndList               bool        //Represents tag #EXT-X-ENDLIST. Indicates no more media segments will be added to the playlist.
	Type                  string      //Possible Values: EVENT or VOD. Represents tag #EXT-X-PLAYLIST-TYPE. If EVENT - segments can only be added to the end of playlist. If VOD - playlist cannot change. If segments need to be removed from playlist, this tag MUST NOT be present
	IFramesOnly           bool        //Represents tag #EXT-X-I-FRAMES-ONLY. If present, segments MUST begin with either a Media Initialization Section or have a EXT-X-MAP tag.
	AllowCache            bool        //Possible Values: YES or NO. Represents tag #EXT-X-ALLOW-CACHE. Versions 3 - 6 only.
	IndependentSegments   bool        //Represents tag #EXT-X-INDEPENDENT-SEGMENTS. Applies to every Media Segment in the playlist.
	StartPoint            *StartPoint //Represents tag #EXT-X-START
}
