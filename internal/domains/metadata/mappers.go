package metadata

import (
	"agentic/commerce/config"
	"agentic/commerce/pkg/specs/api"
	"errors"
	"strconv"

	"github.com/samber/lo"
)

type IContentMapper interface {
	mapContentRequestToModel(*api.MetadataRequest, string) (*MetaDataModel, error)
	mapToMetadataList(res []MetaDataModel) []api.MetadataItemResponse
	mapToMetadataItem(res *MetaDataModel) *api.MetadataItemResponse
}

type contentMapper struct {
	IContentMapper
	config *config.Config
}

func NewContentMapper() IContentMapper {
	return &contentMapper{}
}

func (m *contentMapper) mapContentRequestToModel(req *api.MetadataRequest, uuid string) (*MetaDataModel, error) {
	if req == nil {
		return nil, errors.New("meta data request is empty")
	}

	metadata := map[string]interface{}{
		"desc":     req.MetaData.Desc,
		"images":   req.MetaData.Images,
		"location": req.MetaData.Location,
	}

	v, err := strconv.ParseInt(req.UserId, 10, 64)
	if err != nil {
		return nil, err
	}

	return &MetaDataModel{
		UUid:     lo.ToPtr(uuid),
		UserId:   lo.ToPtr(v),
		Metadata: metadata,
	}, nil

}

func (m *contentMapper) mapToMetadataItem(res *MetaDataModel) *api.MetadataItemResponse {
	if res == nil {
		return nil
	}
	desc, _ := res.Metadata["desc"].(string)
	location, _ := res.Metadata["location"].(string)

	var images []string
	if rawImages, ok := res.Metadata["images"].([]interface{}); ok {
		images = make([]string, 0, len(rawImages))
		for _, v := range rawImages {
			if s, ok := v.(string); ok {
				images = append(images, s)
			}
		}
	}
	return &api.MetadataItemResponse{
		MetaDataBody: api.MetaDataBody{
			Desc:     desc,
			Images:   images,
			Location: location,
		},
	}

}

func (m *contentMapper) mapToMetadataList(res []MetaDataModel) []api.MetadataItemResponse {
	if len(res) == 0 {
		return nil
	}

	out := make([]api.MetadataItemResponse, 0, len(res))

	for _, item := range res {
		desc, _ := item.Metadata["desc"].(string)
		location, _ := item.Metadata["location"].(string)

		var images []string
		if rawImages, ok := item.Metadata["images"].([]interface{}); ok {
			images = make([]string, 0, len(rawImages))
			for _, v := range rawImages {
				if s, ok := v.(string); ok {
					images = append(images, s)
				}
			}
		}

		out = append(out, api.MetadataItemResponse{
			MetaDataBody: api.MetaDataBody{
				Desc:     desc,
				Images:   images,
				Location: location,
			},
		})
	}

	return out
}
