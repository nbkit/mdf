package reg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/nbkit/mdf/utils"
	"github.com/robfig/cron"
	"io"
	"strings"

	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/nbkit/mdf/log"
)

//code 到 RegObject 的缓存
var clientRegMap map[string]*RegObject = make(map[string]*RegObject)

func setRegObjectCache(code string, item *RegObject) {
	clientRegMap[strings.ToLower(code)] = item
}
func getRegObjectCache(code string) *RegObject {
	if obj, ok := clientRegMap[strings.ToLower(code)]; ok {
		return obj
	}
	return nil
}

/**
获取注册中心地址
*/
func getRegistryHost() string {
	registry := utils.Config.App.Registry
	if registry == "" {
		registry = fmt.Sprintf("http://127.0.0.1:%s", utils.Config.App.Port)
	}
	return registry
}

/**
通过token获取上下文
*/
func GetTokenContext(tokenCode string) (*utils.TokenContext, error) {
	//1、权限注册中心、2、应用注册中心，3、本地
	authAddr := ""
	if ser, err := FindServerByCode(utils.Config.Auth.Code); ser != nil {
		authAddr = ser.Address
	} else {
		log.ErrorD(err)
	}
	if authAddr == "" {
		authAddr = fmt.Sprintf("http://127.0.0.1:%s", utils.Config.App.Port)
	}
	client := &http.Client{}
	client.Timeout = 2 * time.Second
	remoteUrl, _ := url.Parse(authAddr)
	remoteUrl.Path = fmt.Sprintf("/api/oauth/token/%s", tokenCode)
	req, err := http.NewRequest("GET", remoteUrl.String(), nil)
	if err != nil {
		return nil, log.ErrorD(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, log.ErrorD(err)
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, log.ErrorD(err)
	}
	var resBodyObj struct {
		Msg   string `json:"msg"`
		Token struct {
			AccessToken string `json:"access_token"`
			Type        string `json:"type"`
		} `json:"token"`
	}
	if err := json.Unmarshal(resBody, &resBodyObj); err != nil {
		return nil, log.ErrorD(err)
	}
	if resp.StatusCode != 200 || resBodyObj.Msg != "" {
		return nil, log.ErrorD(resBodyObj.Msg)
	}
	token := utils.NewTokenContext()
	token, _ = token.FromTokenString(fmt.Sprintf("%s %s", resBodyObj.Token.Type, resBodyObj.Token.AccessToken))
	return token, nil
}

var m_cronCache *cron.Cron

/**
由配置文件信息，注册
*/
func StartClient() {
	if m_cronCache == nil {
		m_cronCache := cron.New()
		m_cronCache.AddFunc("@every 120s", registerDefault)
		m_cronCache.Start()
	}
	registerDefault()
}
func registerDefault() {
	address := utils.Config.App.Address
	if address == "" {
		address = fmt.Sprintf("http://127.0.0.1:%s", utils.Config.App.Port)
	}
	Register(RegObject{
		Code:    utils.Config.App.Code,
		Name:    utils.Config.App.Name,
		Address: address,
		Configs: utils.Config,
	})
}
func Register(item RegObject) error {
	client := &http.Client{}
	client.Timeout = 3 * time.Second
	postBody, err := json.Marshal(item)
	if err != nil {
		return log.ErrorD(err)
	}
	regHost := getRegistryHost()
	remoteUrl, _ := url.Parse(regHost)
	remoteUrl.Path = "/api/regs/register"
	req, err := http.NewRequest("POST", remoteUrl.String(), bytes.NewBuffer([]byte(postBody)))
	if err != nil {
		return log.ErrorD(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return log.ErrorD(err)
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return log.ErrorD(err)
	}
	var resBodyObj struct {
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}
	if err := json.Unmarshal(resBody, &resBodyObj); err != nil {
		return log.ErrorD(err)
	}
	if resp.StatusCode != 200 || resBodyObj.Msg != "" {
		return log.ErrorD(resBodyObj.Msg)
	}
	log.Error().Any("Item", item).Any("RegHost", regHost).Msgf("成功注册")
	return nil
}
func DoHttpRequest(serverCode, method, path string, body io.Reader) ([]byte, error) {
	regs, err := FindServerByCode(serverCode)
	if err != nil {
		return nil, err
	}
	if regs == nil || regs.Address == "" {
		return nil, log.Error().String("serverCode", serverCode).Error("找不到服务")
	}
	serverUrl := regs.Address
	client := &http.Client{}
	remoteUrl, err := url.Parse(serverUrl)
	if err != nil {
		return nil, err
	}
	remoteUrl.Path = path
	req, err := http.NewRequest(method, remoteUrl.String(), body)
	if err != nil {
		return nil, log.ErrorD(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, log.ErrorD(err)
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, log.ErrorD(err)
	}
	if resp.StatusCode != 200 {
		var resBodyObj struct {
			Msg string `json:"msg"`
		}
		if err := json.Unmarshal(resBody, &resBodyObj); err != nil {
			return nil, err
		}
		return nil, log.ErrorD(resBodyObj.Msg)
	}
	return resBody, nil
}

/**
通过编码找到注册对象
*/
func FindServerByCode(serverCode string) (*RegObject, error) {
	if serverCode == "" {
		return nil, nil
	}
	//优先从缓存里取
	if cv := getRegObjectCache(serverCode); cv != nil {
		return cv, nil
	}
	client := &http.Client{}
	client.Timeout = 2 * time.Second
	remoteUrl, _ := url.Parse(getRegistryHost())
	remoteUrl.Path = fmt.Sprintf("/api/regs/%s", serverCode)
	req, err := http.NewRequest("GET", remoteUrl.String(), nil)

	if err != nil {
		return nil, log.ErrorD(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, log.ErrorD(err)
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, log.ErrorD(err)
	}
	var resBodyObj struct {
		Msg  string     `json:"msg"`
		Data *RegObject `json:"data"`
	}
	if err := json.Unmarshal(resBody, &resBodyObj); err != nil {
		return nil, log.ErrorD(err)
	}
	if resp.StatusCode != 200 || resBodyObj.Msg != "" {
		return nil, log.ErrorD(resBodyObj.Msg)
	}
	//设置缓存
	setRegObjectCache(resBodyObj.Data.Code, resBodyObj.Data)

	return resBodyObj.Data, nil
}
