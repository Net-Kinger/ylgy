package main

import (
	"bufio"
	"context"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

var rootContext context.Context

func init() {
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		DialTLSContext:         nil,
		TLSClientConfig:        nil,
		TLSHandshakeTimeout:    0,
		DisableKeepAlives:      false,
		DisableCompression:     false,
		MaxIdleConns:           100,
		MaxIdleConnsPerHost:    100,
		MaxConnsPerHost:        100,
		IdleConnTimeout:        300 * time.Second,
		ResponseHeaderTimeout:  0,
		ExpectContinueTimeout:  0,
		TLSNextProto:           nil,
		ProxyConnectHeader:     nil,
		GetProxyConnectHeader:  nil,
		MaxResponseHeaderBytes: 0,
		WriteBufferSize:        0,
		ReadBufferSize:         0,
		ForceAttemptHTTP2:      false,
	}
	httpClient := &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}
	rootContext = context.WithValue(context.Background(), "client", httpClient)
	rootContext, _ = context.WithTimeout(rootContext, 300*time.Second)
}
func main() {
	reader := bufio.NewReader(os.Stdin)
	defer os.Stdin.Close()
	go func() {
		tag, _ := reader.ReadByte()
		if tag == 's' {
			os.Exit(0)
		}
	}()
	r := gin.Default()
	r.GET("/setToken/:count/:token", setToken)
	r.Run(":80")
}

func setToken(gcontext *gin.Context) {
	param := gcontext.Param("token")
	count := gcontext.Param("count")
	if param == "" {
		gcontext.JSON(
			200,
			gin.H{
				"error": "空Token",
			})
		return
	}
	atoi, err := strconv.Atoi(count)
	if err != nil {
		gcontext.JSON(
			200,
			gin.H{
				"error": "Count异常",
			})
		return
	}
	value := rootContext.Value("client")
	if v, ok := value.(*http.Client); ok {
		req, err := http.NewRequest("GET", `https://cat-match.easygame2021.com/sheep/v1/game/game_over?rank_score=1&rank_state=1&rank_time=`+strconv.Itoa(time.Now().Day())+`&rank_role=1&skin=1`, nil)
		if err != nil {
			gcontext.JSON(
				200,
				gin.H{
					"error": err.Error(),
				})
			return
		}
		req.Header = map[string][]string{
			"t":               {param},
			"Accept-Encoding": {`gzip,compress,br,deflate`},
			"content-type":    {`application/json`},
			"Referer":         {`https://servicewechat.com/wx141bfb9b73c970a9/17/page-frame.html`},
			"User-Agent":      {`Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/600.1.00 (KHTML, like Gecko) Mobile/1888888 MicroMessenger/8.0.30(0x00111100) NetType/WIFI Language/zh_CN`},
		}
		timeout, cancelFunc := context.WithTimeout(rootContext, 100*time.Second)
		for i := 0; i < atoi; i++ {
			select {
			case <-timeout.Done():
				gcontext.String(200, "%s", "105Bug 通道关闭")
				return
			case <-time.NewTicker(1000 * time.Millisecond).C:
				go func() {
					_, err := v.Do(req)
					if err != nil {
						cancelFunc()
					}
				}()
			}
		}
	} else {
		gcontext.JSON(
			http.StatusOK,
			gin.H{
				"error": "断言*http.Client失败",
			})
	}
	return
}
