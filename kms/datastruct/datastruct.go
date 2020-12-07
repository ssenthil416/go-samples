package datastruct

// RequestKMS KMS Request
type RequestKMS struct {
	KeyRingID string            `json:"krid"`
	ClientID  string            `json:"clientid"`
	KeyID     string            `json:"keyid"`
	Labels    map[string]string `json:"labels"`
	DataSet   map[string]string `json:"dataset"`
}
