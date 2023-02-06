package contribution

import (
	"Go-Live/consts"
	receive "Go-Live/interaction/receive/contribution/article"
	response "Go-Live/interaction/response/contribution/article"
	"Go-Live/models/common"
	"Go-Live/models/contribution/article"
	"Go-Live/models/contribution/article/classification"
	"Go-Live/models/contribution/article/comments"
	"Go-Live/utils/conversion"
	"encoding/json"
	"fmt"
	"github.com/dlclark/regexp2"
)

func CreateArticleContribution(data *receive.CreateArticleContributionReceiveStruct, userID uint) (results interface{}, err error) {
	//进行内容判断
	for _, v := range data.Label {
		vRune := []rune(v) //避免中文占位问题
		if len(vRune) > 7 {
			return nil, fmt.Errorf("标签长度不能大于7位")
		}
	}

	coverImg, _ := json.Marshal(common.Img{
		Src: data.Cover,
		Tp:  data.CoverUploadType,
	})
	//正则匹配替换url
	//取url前缀
	prefix, err := conversion.SwitchTypeAsUrlPrefix(data.ArticleContributionUploadType)
	if err != nil {
		return nil, fmt.Errorf("保存资源方式不存在")
	}
	//正则匹配替换
	reg := regexp2.MustCompile(`(?<=(img[^>]*src="))[^"]*?`+prefix, 0)
	match, err := reg.Replace(data.Content, consts.UrlPrefixSubstitution, -1, -1)
	data.Content = match
	//插入数据
	articlesContribution := article.ArticlesContribution{
		Uid:                userID,
		ClassificationID:   data.ClassificationID,
		Title:              data.Title,
		Cover:              coverImg,
		Timing:             conversion.BoolTurnInt8(*data.Timing),
		TimingTime:         data.TimingTime,
		Label:              conversion.MapConversionString(data.Label),
		Content:            data.Content,
		ContentStorageType: data.ArticleContributionUploadType,
		IsComments:         conversion.BoolTurnInt8(*data.Comments),
		Heat:               0,
	}

	if *data.Timing {
		//发布视频后进行的推送相关（待开发）
	}
	if !articlesContribution.Create() {
		return nil, fmt.Errorf("保存失败")
	}
	return "保存成功", nil
}

func GetArticleContributionListByUser(data *receive.GetArticleContributionListByUserReceiveStruct, userID uint) (results interface{}, err error) {
	articlesContribution := new(article.ArticlesContributionList)
	if !articlesContribution.GetListByUid(data.UserID) {
		return nil, fmt.Errorf("查询失败")
	}
	return response.GetArticleContributionListByUserResponse(articlesContribution), nil
}

func GetArticleContributionByID(data *receive.GetArticleContributionByIDReceiveStruct, userID uint) (results interface{}, err error) {
	articlesContribution := new(article.ArticlesContribution)
	if !articlesContribution.GetInfoByID(data.ArticleID) {
		return nil, fmt.Errorf("查询失败")
	}
	return response.GetArticleContributionByIDResponse(articlesContribution), nil
}

func ArticlePostComment(data *receive.ArticlesPostCommentReceiveStruct, userID uint) (results interface{}, err error) {
	ct := comments.Comment{
		PublicModel: common.PublicModel{ID: data.ContentID},
	}
	CommentFirstID := ct.GetCommentFirstID()

	ctu := comments.Comment{
		PublicModel: common.PublicModel{ID: data.ContentID},
	}
	CommentUserID := ctu.GetCommentUserID()
	comment := comments.Comment{
		Uid:            userID,
		ContributionID: data.ArticleID,
		Context:        data.Content,
		CommentID:      data.ContentID,
		CommentUserID:  CommentUserID,
		CommentFirstID: CommentFirstID,
	}
	if !comment.Create() {
		return nil, fmt.Errorf("发布失败")
	}
	return "发布成功", nil
}

func GetArticleComment(data *receive.GetArticleCommentReceiveStruct) (results interface{}, err error) {
	articlesContribution := new(article.ArticlesContribution)
	if !articlesContribution.GetArticleComments(data.ArticleID, data.PageInfo) {
		return nil, fmt.Errorf("查询失败")
	}
	return response.GetArticleContributionCommentsResponse(articlesContribution), nil
}

func GetArticleClassificationList() (results interface{}, err error) {
	cn := new(classification.ClassificationsList)
	err = cn.FindAll()
	if err != nil {
		return nil, fmt.Errorf("查询失败")
	}
	return response.GetArticleClassificationListResponse(cn), nil
}

func GetArticleTotalInfo() (results interface{}, err error) {
	//查询文章数量
	articleNm := new(int64)
	al := new(article.ArticlesContributionList)
	al.GetAllCount(articleNm)
	//查询文章分类信息
	cn := make(classification.ClassificationsList, 0)
	err = cn.FindAll()
	if err != nil {
		return nil, fmt.Errorf("查询失败")
	}
	cnNum := int64(len(cn))

	return response.GetArticleTotalInfoResponse(&cn, articleNm, cnNum), nil
}
