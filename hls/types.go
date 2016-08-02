package hls

//Segments implements golang/sort interface to sort a Segment slice by Segment ID
type Segments []*Segment

func (s Segments) Len() int {
	return len(s)
}
func (s Segments) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s Segments) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}
