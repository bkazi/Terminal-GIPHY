package terminalGiphy

type Response struct {
	Data []data `json:"data"`
}

type data struct {
	Images images `json:"images"`
}

type images struct {
	FixedHeight map[string]interface{} `json:"fixed_height"`
	Preview     map[string]interface{} `json:"fixed_height_small_still"`
}
