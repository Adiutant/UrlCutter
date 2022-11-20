package http_server

import (
	"fmt"
	gin "github.com/gin-gonic/gin"
	logrus "github.com/sirupsen/logrus"
	"hash/crc32"
	"math/rand"
	"net/http"
	"net/url"
	"shortUrl/model"
	"time"
)

type Request struct {
	Url string `json:"url"`
}

const lowercase = "abcdefghijklmnopqrstunwxyz"
const uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const numbers = "0123456789"
const symbols = "-_=+/.|\\"

type UrlServer struct {
	serverInstance *gin.Engine
	mostViewed     map[string]model.UrlInfo
	logger         *logrus.Logger
	config         model.UrlServerConfig
	whitelist      string
}

func MakeUrlServer() (*UrlServer, error) {
	wl := lowercase + uppercase + numbers + symbols
	wl = Shuffle(wl)
	logger := logrus.New()
	config := model.UrlServerConfig{}
	engine := gin.Default()
	err := config.ReadConfig("urlserver")
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &UrlServer{
		serverInstance: engine,
		mostViewed:     make(map[string]model.UrlInfo),
		logger:         logger,
		config:         config,
		whitelist:      wl,
	}, nil

}

func (s *UrlServer) SetRoutes() {
	s.serverInstance.POST("/cut", func(c *gin.Context) {
		req := Request{}
		err := c.BindJSON(&req)
		if err != nil {
			s.logger.Error(err)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		if shortUrl, err := s.makeShortUrl(req.Url); err == nil {
			s.mostViewed[shortUrl] = model.UrlInfo{Url: req.Url, ShortUrl: shortUrl, TimesCalled: 0, Position: 0}
			c.String(http.StatusOK, "Short url is %s", "localhost:8080/"+shortUrl)
		} else {
			s.logger.Error(err)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	})
	s.serverInstance.GET("/:shortpath", func(c *gin.Context) {
		urlInfo, exist := s.mostViewed[c.Param("shortpath")]
		if !exist {
			s.logger.Error(fmt.Errorf("short url is not found"))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.Redirect(http.StatusPermanentRedirect, urlInfo.Url)
	})
}

func (s *UrlServer) Run() error {
	err := s.serverInstance.Run(":8080")
	if err != nil {
		s.logger.Error(err)
		return err
	}
	return nil
}
func (s *UrlServer) makeShortUrl(urlString string) (string, error) {
	shortUrl := make([]rune, 0)
	url, err := url.Parse(urlString)
	if err != nil {

		return "", err
	}
	if url.String() == "" {
		return "", fmt.Errorf("Error in building short url")
	}
	shortUrlIndexes := crc32.Checksum([]byte(url.String()), crc32.IEEETable)
	for shortUrlIndexes != 0 {
		position := shortUrlIndexes % 8
		position %= uint32(len(s.whitelist))
		shortUrl = append(shortUrl, []rune(s.whitelist)[position])
		shortUrlIndexes /= 8
	}

	return string(shortUrl), nil
}
func Shuffle(source string) string {
	slice := []rune(source)
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for n := len(slice); n > 0; n-- {
		randIndex := r.Intn(n)
		slice[n-1], slice[randIndex] = slice[randIndex], slice[n-1]
	}
	return string(slice)
}
