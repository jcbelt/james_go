package guestApi

import (
	"bytes"
	"encoding/xml"
	"net/http"
    "net/url"
	"io/ioutil"
	"code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"
    "strconv"
)


// Structs for parsed XML response data
type Open struct {
    Day int `xml:"day"`
    Time string `xml:"time"`
}
type Close struct {
    Day int `xml:"day"`
    Time string `xml:"time"`
}
type Period struct {
    Open Open `xml:"open"`
    Close Close `xml:"close"`
}
type Schedule struct {
    Periods  []Period `xml:"period"`
}
type Restaurant struct {
	Name string `xml:"biz_name"`
    ID string `xml:"biz_id"`
	Schedule Schedule `xml:"schedule"`
}



// Generic Api Response
type ApiResponse struct {
    XMLName xml.Name `xml:"response"`
    Status int `xml:status`
    Message string `xml:message`
}

// Decode the XML into data structures
func (response *GetRestaurantsResponse) decodeXML(xmlData []byte) (err error){

    reader := bytes.NewReader(xmlData)
    decoder := xml.NewDecoder(reader)
    decoder.CharsetReader = charset.NewReader
    err = decoder.Decode(response)
    
    return
}


// Generic Api Request
type ApiRequest struct {
    Url url.URL
}

// Send the request 
func (request *ApiRequest) Send() (xml []byte, err error){

    resp, err := http.Get(request.Url.String());
    if err != nil {
        return 
    }

    defer resp.Body.Close()
    xml, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        return
    }

    return
}



// Request/Response structs for GetRestaurants API
type GetRestaruantsRequest struct {
    ApiRequest
    Lat int
    Lon int
}
type GetRestaurantsResponse struct {
    ApiResponse
    Restaurants []Restaurant `xml:"data>restaurant"`
}


// Send the request 
func (request *GetRestaruantsRequest) Send() (response GetRestaurantsResponse, err error){
    values := url.Values{}
    values.Add("lat", strconv.Itoa(request.Lat))
    values.Add("lon", strconv.Itoa(request.Lon))

    request.Url.Scheme ="https"
    request.Url.Host = "nowaitapp.com"
    request.Url.Path = "consumerAPI/public/getRestaurants"
    request.Url.RawQuery = values.Encode()

    xml, err := request.ApiRequest.Send()

    err = response.decodeXML(xml)
    return
}