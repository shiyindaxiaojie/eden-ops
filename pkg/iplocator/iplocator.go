package iplocator

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// IPLocator IP定位器
type IPLocator struct {
	accessKey string
	secretKey string
	baseURL   string
}

// NewIPLocator 创建IP定位器
func NewIPLocator(accessKey, secretKey, baseURL string) *IPLocator {
	return &IPLocator{
		accessKey: accessKey,
		secretKey: secretKey,
		baseURL:   baseURL,
	}
}

// IPInfo IP信息
type IPInfo struct {
	IP          string  `json:"ip"`
	Country     string  `json:"country"`
	CountryCode string  `json:"country_code"`
	City        string  `json:"city"`
	Region      string  `json:"region"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Timezone    string  `json:"timezone"`
}

// Locate 定位IP地址
func (l *IPLocator) Locate(ip string) (*IPInfo, error) {
	url := fmt.Sprintf("%s/v1/%s?access_key=%s", l.baseURL, ip, l.accessKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var info IPInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}

	return &info, nil
}
