package googleflights

import "net/url"

// URLOptions holds query parameters outside the TFS protobuf message.
type URLOptions struct {
	Language string
	Region   string
	Currency string
	TFU      string // Optional, already encoded independently of TFS.
}

// BuildURL encodes the query and adds localization and optional UI state.
// It does not validate search settings or fetch results.
func BuildURL(query TFS, options URLOptions) (*url.URL, error) {
	u, err := url.Parse(baseURL + "/search")
	if err != nil {
		return nil, err
	}
	params := u.Query()
	params.Set("tfs", EncodeTFS(query))
	for key, value := range map[string]string{
		"hl": options.Language, "gl": options.Region,
		"curr": options.Currency, "tfu": options.TFU,
	} {
		if value != "" {
			params.Set(key, value)
		}
	}
	u.RawQuery = params.Encode()
	return u, nil
}
