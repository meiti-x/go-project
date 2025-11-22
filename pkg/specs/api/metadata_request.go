package api

type MetadataRequest struct {
	UserId   string       `json:"user_id"`
	MetaData MetaDataBody `json:"meta_data"`
}

type MetaDataBody struct {
	Desc     string   `json:"desc"`
	Images   []string `json:"images"`
	Location string   `json:"location"`
}

type MetadataItemResponse struct {
	MetaDataBody
}
