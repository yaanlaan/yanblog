package v1

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"time"
	"yanblog/model"
	"yanblog/utils"
	"yanblog/utils/errmsg"

	"github.com/gin-gonic/gin"
)

// Url 定义 sitemap 中的 url 节点
type Url struct {
	Loc        string  `xml:"loc"`
	LastMod    string  `xml:"lastmod"`
	ChangeFreq string  `xml:"changefreq"`
	Priority   float32 `xml:"priority"`
}

// UrlSet 定义 sitemap 的根节点
type UrlSet struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	Urls    []Url    `xml:"url"`
}

// GetSitemap 生成并返回 sitemap.xml
func GetSitemap(c *gin.Context) {
	articles, code := model.GetSitemapData()
	if code != errmsg.SUCCESS {
		c.Status(http.StatusInternalServerError)
		return
	}

	baseUrl := utils.ServerConfig.Server.SiteUrl
	// 如果未配置 SiteUrl，默认使用 localhost
	if baseUrl == "" {
		baseUrl = "http://localhost:5173"
	}

	urlSet := UrlSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		Urls:  make([]Url, 0),
	}

	// 1. 添加静态页面
	staticPages := []string{"", "/category", "/archive", "/about"}
	for _, page := range staticPages {
		urlSet.Urls = append(urlSet.Urls, Url{
			Loc:        fmt.Sprintf("%s%s", baseUrl, page),
			LastMod:    time.Now().Format("2006-01-02"),
			ChangeFreq: "daily",
			Priority:   0.8,
		})
	}

	// 2. 添加所有文章页面
	for _, art := range articles {
		urlSet.Urls = append(urlSet.Urls, Url{
			Loc:        fmt.Sprintf("%s/article/info/%d", baseUrl, art.ID),
			LastMod:    art.UpdatedAt.Format("2006-01-02"),
			ChangeFreq: "weekly",
			Priority:   0.6,
		})
	}

	// 生成 XML
	c.Header("Content-Type", "application/xml")
	c.XML(http.StatusOK, urlSet)
}
