package dash

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

//Periods ...
type Periods []*Period

//Segments is a string of S elements that implements Sort interface
type Segments []*S

func (s Segments) Len() int {
	return len(s)
}
func (s Segments) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s Segments) Less(i, j int) bool {
	return s[i].T < s[j].T
}

//Subsets ...
type Subsets []*Subset

//AdaptationSets ...
type AdaptationSets []*AdaptationSet

//Representations ...
type Representations []*Representation

//CustomTime is a custom type of time.Time that implements XML marshaller and unmarshaller
type CustomTime struct {
	time.Time
}

//UnmarshalXMLAttr implementes UnmarshalerAttr interface for CustomTime
func (c *CustomTime) UnmarshalXMLAttr(attr xml.Attr) (err error) {
	c.Time, err = time.Parse(time.RFC3339Nano, attr.Value)
	return
}

//MarshalXMLAttr implementes MarshalerAttr interface for CustomTime
func (c *CustomTime) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	attr := xml.Attr{Name: name}
	attr.Value = c.Time.Format(time.RFC3339Nano)
	if attr.Value == "" {
		return attr, errors.New("Unable to format CustomTime")
	}
	return attr, nil
}

//CustomDuration is a custom type of time.Duration that implements XML marshaller and unmarshaller
type CustomDuration struct {
	time.Duration
}

//UnmarshalXMLAttr implementes UnmarshalerAttr interface for CustomDuration
func (c *CustomDuration) UnmarshalXMLAttr(attr xml.Attr) (err error) {
	//Removes 'PT' from attribute value before parsing duration
	c.Duration, err = time.ParseDuration(strings.ToLower(attr.Value[2:len(attr.Value)]))
	return
}

//MarshalXMLAttr implementes MarshalerAttr interface for CustomDuration
func (c *CustomDuration) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	attr := xml.Attr{Name: name}
	attr.Value = fmt.Sprintf("PT%s", strings.ToUpper(c.Duration.String()))
	return attr, nil
}

//CustomInt is a custom type for UIntVectorType that implements XML marshaller and unmarshaller
type CustomInt struct {
	CI []int
}

//UnmarshalXMLAttr implementes UnmarshalerAttr interface for CustomInt
func (c *CustomInt) UnmarshalXMLAttr(attr xml.Attr) (err error) {
	var ss []string
	var i int64
	if len(ss) > 0 {
		ss = strings.Split(attr.Value, " ")
		for _, s := range ss {
			i, err = strconv.ParseInt(s, 10, 8)
			if err != nil {
				return
			}
			c.CI = append(c.CI, int(i))
		}
	}
	return
}

//MarshalXMLAttr implementes MarshalerAttr interface for CustomInt
func (c *CustomInt) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	var ci []string
	var attr xml.Attr
	if len(c.CI) > 0 {
		for _, i := range c.CI {
			ci = append(ci, strconv.Itoa(i))
		}
		attr = xml.Attr{Name: name, Value: strings.Join(ci, " ")}
	}
	return attr, nil
}
