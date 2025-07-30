package model

import (
	"yanblog/utils"
	"yanblog/utils/errmsg"

	"fmt"
	"context"
	"mime/multipart"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

// 从配置文件中读取七牛云相关配置
var Zone = utils.ServerConfig.Qiniu.Zone
var AccessKey = utils.ServerConfig.Qiniu.AccessKey
var SecretKey = utils.ServerConfig.Qiniu.SecretKey
var Bucket = utils.ServerConfig.Qiniu.Bucket
var ImgUrl = utils.ServerConfig.Qiniu.Server 

// UpLoadFile 上传文件到七牛云
// 参数: file - 要上传的文件, fileSize - 文件大小
// 返回: 文件访问URL和状态码
func UpLoadFile(file multipart.File, fileSize int64) (string, int) {

	putPolicy := storage.PutPolicy{
		Scope: Bucket,
	}
	
	// 打印调试信息
	fmt.Println(Bucket)
	fmt.Println(ImgUrl)
	fmt.Println(AccessKey)
	fmt.Println(SecretKey)

	mac := qbox.NewMac(AccessKey, SecretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{
		Zone:          SetZone(Zone),
		UseCdnDomains: false,
		UseHTTPS:      false,
	}

	putExtra := storage.PutExtra{}

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	err := formUploader.PutWithoutKey(context.Background(), &ret, upToken, file, fileSize, &putExtra)
	if err != nil {
		return "", errmsg.ERROR
	}
	
	url := ImgUrl + ret.Key
	return url, errmsg.SUCCESS
}

// SetZone 根据配置设置七牛云存储区域
// 参数: zone - 区域代码(0:华东 1:华南 2:北美 3:新加坡)
// 返回: 对应的存储区域配置
func SetZone(zone int) *storage.Zone {
	switch zone {
	case 0:
		return &storage.ZoneHuadong   // 华东
	case 1:
		return &storage.ZoneHuanan    // 华南
	case 2:
		return &storage.ZoneBeimei    // 北美
	case 3:
		return &storage.ZoneXinjiapo  // 新加坡
	default:
		return &storage.ZoneHuadong   // 默认华东
	}
}