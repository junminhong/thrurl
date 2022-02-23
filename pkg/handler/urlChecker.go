package handler

import (
	"bytes"
	"encoding/json"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type client struct {
	ClientId      string `json:"clientId"`
	ClientVersion string `json:"clientVersion"`
}
type threatInfo struct {
	ThreatTypes      []string           `json:"threatTypes"`
	PlatformTypes    []string           `json:"platformTypes"`
	ThreatEntryTypes []string           `json:"threatEntryTypes"`
	ThreatEntries    []threatEntriesUrl `json:"threatEntries"`
}
type safeUrlReq struct {
	Client     client     `json:"client"`
	ThreatInfo threatInfo `json:"threatInfo"`
}
type threatEntriesUrl struct {
	Url string `json:"url"`
}

type safeUrlResponse struct {
	Matches []struct {
		ThreatType      string `json:"threatType"`
		PlatformType    string `json:"platformType"`
		ThreatEntryType string `json:"threatEntryType"`
		Threat          struct {
			Url string `json:"url"`
		} `json:"threat"`
	} `json:"matches"`
}

// SafeUrlCheck bonus function
func SafeUrlCheck(url string) (bool, string) {
	// 透過google safe browser check source url是不是惡意網站
	// protect all member computer security
	safeUrlReq := safeUrlReq{
		Client: client{ClientId: "thrurl", ClientVersion: "1.5.2"},
		ThreatInfo: threatInfo{
			ThreatTypes:      []string{"MALWARE", "SOCIAL_ENGINEERING", "UNWANTED_SOFTWARE"},
			PlatformTypes:    []string{"ANY_PLATFORM"},
			ThreatEntryTypes: []string{"URL"},
			ThreatEntries: []threatEntriesUrl{
				{Url: url},
			},
		},
	}
	data, _ := json.Marshal(safeUrlReq)
	response, err := http.Post(viper.GetString("APP.SAFE_BROWSING_API"), "application/json", bytes.NewBuffer(data))
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err.Error())
		return true, ""
	}
	safeUrlResponse := &safeUrlResponse{}
	json.Unmarshal(body, safeUrlResponse)
	if len(safeUrlResponse.Matches) == 0 {
		return false, ""
	}
	return true, safeUrlResponse.Matches[0].ThreatType
}

func UrlLifeCheck(source string) bool {
	url, _ := url.ParseRequestURI(source)
	if url == nil {
		return false
	}
	return true
}
