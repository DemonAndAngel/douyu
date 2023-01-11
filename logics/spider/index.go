package spider

import (
	"bytes"
	"douyu/utils/helpers"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

const (
	GetIndexUri = `https://www.douyu.com/%s`
)

type GetIndexResponse struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

func GetIndex(roomId string) (resp GetIndexResponse, err error) {
	code, body, err := helpers.Request("GET", fmt.Sprintf(GetIndexUri, roomId), nil, nil)
	if err != nil {
		return
	}
	if code != http.StatusOK {
		err = errors.New(fmt.Sprintf("GetIndex code not 200;code:%d", code))
		return
	}
	// 解析页面
	dom, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return
	}
	resp.Title = dom.Find("#js-player-title").Find(".Title-header").Text()
	resp.Name = dom.Find("#js-player-title").Find(".Title-anchorNameH2").Text()
	return
}
