package spider

import (
	"douyu/utils/helpers"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"net/http"
	"net/url"
	"time"
)

const (
	GetH5PlayUri = `https://www.douyu.com/lapi/live/getH5Play/%s`
	SignMode1    = `Sign1`
	SignMode2    = `Sign2`
)

func GetH5Play(roomId, signMode string) (resp GetH5PlayResponse, mode string, err error) {
	did := viper.GetString("appConfig.spider.did")
	tt := time.Now().Unix()
	params := url.Values{}
	params.Set("did", did)
	params.Set("tt", cast.ToString(tt))
	if signMode == "" {
		// 先用sign1
		mode = SignMode1
		sign, v := helpers.ParseSignStr(helpers.Sign1(cast.ToInt64(roomId), did, tt))
		params.Set("sign", sign)
		params.Set("v", v)
		resp, err = getH5PlayRequest(roomId, params)
		if err != nil {
			// 这里使用sign2
			mode = SignMode2
			sign, v = helpers.ParseSignStr(helpers.Sign2(cast.ToInt64(roomId), did, tt))
			params.Set("sign", sign)
			params.Set("v", v)
			resp, err = getH5PlayRequest(roomId, params)
			return
		}
		return
	} else if signMode == SignMode1 {
		mode = SignMode1
		sign, v := helpers.ParseSignStr(helpers.Sign1(cast.ToInt64(roomId), did, tt))
		params.Set("sign", sign)
		params.Set("v", v)
		resp, err = getH5PlayRequest(roomId, params)
		return
	} else if signMode == SignMode2 {
		mode = SignMode2
		sign, v := helpers.ParseSignStr(helpers.Sign2(cast.ToInt64(roomId), did, tt))
		params.Set("sign", sign)
		params.Set("v", v)
		resp, err = getH5PlayRequest(roomId, params)
		return
	} else {
		mode = ""
		err = errors.New("sign mode error")
		return
	}
}

type GetH5PlayResponse struct {
	Error int    `json:"error"`
	Msg   string `json:"msg"`
	Data  struct {
		RoomId       int    `json:"room_id"`
		IsMixed      bool   `json:"is_mixed"`
		MixedLive    string `json:"mixed_live"`
		MixedUrl     string `json:"mixed_url"`
		RtmpCdn      string `json:"rtmp_cdn"`
		RtmpUrl      string `json:"rtmp_url"`
		RtmpLive     string `json:"rtmp_live"`
		ClientIp     string `json:"client_ip"`
		InNA         int    `json:"inNA"`
		RateSwitch   int    `json:"rateSwitch"`
		Rate         int    `json:"rate"`
		CdnsWithName []struct {
			Name   string `json:"name"`
			Cdn    string `json:"cdn"`
			IsH265 bool   `json:"isH265"`
		} `json:"cdnsWithName"`
		Multirates []struct {
			Name       string `json:"name"`
			Rate       int    `json:"rate"`
			HighBit    int    `json:"highBit"`
			Bit        int    `json:"bit"`
			DiamondFan int    `json:"diamondFan"`
		} `json:"multirates"`
		IsPassPlayer int         `json:"isPassPlayer"`
		Eticket      interface{} `json:"eticket"`
		Online       int         `json:"online"`
		MixedCDN     string      `json:"mixedCDN"`
		P2P          int         `json:"p2p"`
		StreamStatus int         `json:"streamStatus"`
		Smt          int         `json:"smt"`
		P2PMeta      struct {
			BestPlayBufferMs       int     `json:"best_play_buffer_ms"`
			DnsTmoms               int     `json:"dns_tmoms"`
			Dyxp2PApiDomain        string  `json:"dyxp2p_api_domain"`
			Dyxp2PBakEdge          string  `json:"dyxp2p_bak_edge"`
			Dyxp2PDomain           string  `json:"dyxp2p_domain"`
			Dyxp2PSubsDomain       string  `json:"dyxp2p_subs_domain"`
			Dyxp2PSugEgde          string  `json:"dyxp2p_sug_egde"`
			FakePeerFactor         float64 `json:"fake_peer_factor"`
			HttpdnsAgingSec        int     `json:"httpdns_aging_sec"`
			HttpdnsFailnumMax      int     `json:"httpdns_failnum_max"`
			MaxPlayBufferMs        int     `json:"max_play_buffer_ms"`
			MinPlayBufferMs        int     `json:"min_play_buffer_ms"`
			MinPlayDelayTimeSecond int     `json:"min_play_delay_time_second"`
			Pcdn                   string  `json:"pcdn"`
			PcdnTable              string  `json:"pcdn_table"`
			PlayFastScale          int     `json:"play_fast_scale"`
			PlaySlowScale          int     `json:"play_slow_scale"`
			SdkapiIplist           string  `json:"sdkapi_iplist"`
			SdkapiIpv6List         string  `json:"sdkapi_ipv6list"`
			StreamProps            []struct {
				Sid      string `json:"sid"`
				TxSecret string `json:"txSecret"`
			} `json:"stream_props"`
			UseP2PAgent                int    `json:"use_p2p_agent"`
			UsingAliHttpdns            int    `json:"using_ali_httpdns"`
			UsingHttpdns               int    `json:"using_httpdns"`
			UsingSyshost               int    `json:"using_syshost"`
			UsingTctHttpdns            int    `json:"using_tct_httpdns"`
			Xp2PApiDomain              string `json:"xp2p_api_domain"`
			Xp2PDomain                 string `json:"xp2p_domain"`
			Xp2PEmergencyWndMs         int    `json:"xp2p_emergency_wnd_ms"`
			Xp2PEnableRepullFullStream int    `json:"xp2p_enable_repull_full_stream"`
			Xp2PEnableWebsocket        int    `json:"xp2p_enable_websocket"`
			Xp2PMaxCacheSliceNum       int    `json:"xp2p_max_cache_slice_num"`
			Xp2PSubsDomain             string `json:"xp2p_subs_domain"`
			Xp2PTxDelay                int    `json:"xp2p_txDelay"`
			Xp2PTxExpire               int    `json:"xp2p_txExpire"`
			Xp2PTxSecret               string `json:"xp2p_txSecret"`
			Xp2PTxTime                 string `json:"xp2p_txTime"`
			ZeroPackageDomain          string `json:"zero_package_domain"`
		} `json:"p2pMeta"`
		P2PCid               int64  `json:"p2pCid"`
		P2PCids              string `json:"p2pCids"`
		Player1              string `json:"player_1"`
		H265P2P              int    `json:"h265_p2p"`
		H265P2PCid           int    `json:"h265_p2p_cid"`
		H265P2PCids          string `json:"h265_p2p_cids"`
		Acdn                 string `json:"acdn"`
		Av1Url               string `json:"av1_url"`
		RtcStreamUrl         string `json:"rtc_stream_url"`
		RtcStreamConfig      string `json:"rtc_stream_config"`
		PictureQualitySwitch int    `json:"pictureQualitySwitch"`
	} `json:"data"`
}

func getH5PlayRequest(roomId string, params url.Values) (resp GetH5PlayResponse, err error) {
	code, body, err := helpers.Request("POST", fmt.Sprintf(GetH5PlayUri, roomId), nil, params)
	if err != nil {
		return
	}
	if code != http.StatusOK {
		err = errors.New(fmt.Sprintf("GetH5Play code not 200;code:%d", code))
		return
	}
	_ = json.Unmarshal(body, &resp)
	return
}
