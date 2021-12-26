package utils

// .版本 2
// .支持库 spec
// 局_ApiUrl ＝ 局_ApiUrl ＋ “/robot/stat/data”
// 调试输出 (局_ApiUrl)
// 局_Json.创建 ()
// 局_Json.置文本 (“botype”, botype)
// 局_Json.置文本 (“botid”, botid)
// ' 局_Json.置文本 (“botime”, 时间_取现行时间戳 (真))
// 局_Json.置文本 (“botoken”, botoken)
// 局_Text ＝ 局_Json.到文本 (“.”)
// 局_Text ＝ AnsiToUtf8Str (局_Text)
import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"gitee.com/lyhuilin/QN/constvar"

	"gitee.com/lyhuilin/my_admin/model/bot"
	"gitee.com/lyhuilin/util"
	"github.com/spf13/viper"
)

var lastUpdateStat time.Time

// var b bot.BotOnline
func UpdateRobotStatToMyAdmin(botid string, isOnline bool) (retText string, err error) {
	if !isOnline {
		// 不在线
		return
	}

	// now := time.Now()

	myadminApiServerUrl := viper.GetString("myadmin_api_server_url")
	myadminApiServerToken := viper.GetString("myadmin_api_server_token")
	if len(myadminApiServerUrl) <= 0 || len(myadminApiServerToken) <= 0 {
		return
	}
	apiURL := fmt.Sprintf("%s/robot/stat/data", myadminApiServerUrl)

	var b bot.BotOnline
	b.Botype = constvar.APP_NAME
	b.Botid = botid
	// b.LastloginTime = ""
	b.Botoken = myadminApiServerToken

	dataBody := util.JsonEncode(b)

	rep, err := http.NewRequest("POST", apiURL, strings.NewReader(dataBody))
	if err != nil {
		return
	}
	rep.Header.Add("Content-Type", "application/json")
	// rep.Header.Add("HLTYClient", "haowu_video")

	httpClient := &http.Client{}
	// httpClient.Timeout = app.Timeout

	response, err := httpClient.Do(rep)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		err = fmt.Errorf("请求错误:%d", response.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	if len(body) > 0 {
		retText = string(body)
	}
	return
}
