package contribution

import (
	receive "Go-Live/interaction/receive/contribution/video"
	response "Go-Live/interaction/response/contribution/video"
	"Go-Live/models/common"
	"Go-Live/models/contribution/video"
	"Go-Live/utils/conversion"
	"encoding/json"
	"fmt"
)

func CreateVideoContribution(data *receive.CreateVideoContributionReceiveStruct, userID uint) (results interface{}, err error) {
	videoSrc, _ := json.Marshal(common.Img{
		Src: data.Video,
		Tp:  data.VideoUploadType,
	})
	coverImg, _ := json.Marshal(common.Img{
		Src: data.Cover,
		Tp:  data.CoverUploadType,
	})
	videoContribution := video.VideosContribution{
		Uid:           userID,
		Title:         data.Title,
		Video:         videoSrc,
		Cover:         coverImg,
		VideoDuration: data.VideoDuration,
		Reprinted:     conversion.BoolTurnInt8(*data.Reprinted),
		Timing:        conversion.BoolTurnInt8(*data.Timing),
		TimingTime:    data.TimingTime,
		Label:         conversion.MapConversionString(data.Label),
		Introduce:     data.Introduce,
		Heat:          0,
	}
	if *data.Timing {
		//发布视频后进行的推送相关（待开发）
	}
	if !videoContribution.Create() {
		return nil, fmt.Errorf("保存失败")
	}
	return "保存成功", nil
}

func GetVideoContributionByID(data *receive.GetVideoContributionByIDReceiveStruct, userID uint) (results interface{}, err error) {
	videoInfo := new(video.VideosContribution)
	err = videoInfo.FindByID(data.VideoID)
	if err != nil {
		return nil, err
	}
	recommendList := new(video.VideosContributionList)
	err = recommendList.GetRecommendList()
	if err != nil {
		return nil, err
	}
	res := response.GetVideoContributionByIDResponse(videoInfo, recommendList)
	return res, nil
}
