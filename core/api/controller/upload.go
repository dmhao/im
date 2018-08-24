package controller

import (
	"context"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"im/core/log"
	"image"
	"image/jpeg"
	"mime/multipart"
	"path"
	"strconv"
	"strings"
	"time"
)

type imageUploadRes struct {
	Url    string
	Width  int
	Height int
}

type videoUploadRes struct {
	Thumb  string
	Length float64
	Size   int
	Url    string
}

type audioUploadRes struct {
	Length float64
	Size   int
	Url    string
}

func UploadImage(c *gin.Context) {
	imageType := c.PostForm("ImageType")
	fileHeader, err := c.FormFile("Image")
	if imageType == "" {
		mContext{c}.ErrorResponse(ParamsError, ParamsErrorMsg)
		return
	}
	if err != nil {
		log.Warnln("图片临时文件获取失败", err)
		mContext{c}.ErrorResponse(ParamsDataError, ParamsDataErrorMsg)
		return
	}
	fd, err := fileHeader.Open()
	if err != nil {
		log.Warnln("上传图片文件打开失败", err)
		mContext{c}.ErrorResponse(ServiceError, ServiceErrorMsg)
		return
	}
	defer fd.Close()

	imageExt := path.Ext(fileHeader.Filename)
	var img image.Image
	if imageExt == ".jpg" || imageExt == ".jpeg" {
		img, err = jpeg.Decode(fd)
		if err != nil {
			mContext{c}.ErrorResponse(ServiceError, ServiceErrorMsg)
			return
		}
	}
	if img != nil {
		b := img.Bounds()
		width := b.Max.X
		height := b.Max.Y
		filePath := makeImageFilePath(imageType, imageExt, width, height)
		qnOssImageUpload(c, fileHeader, filePath, width, height)
	} else {
		mContext{c}.ErrorResponse(ServiceError, ServiceErrorMsg)
		return
	}
}

func makeImageFilePath(imageType string, imageExt string, width int, height int) string {
	timeNow := time.Now().UnixNano()
	currentMonth := time.Now().Format("2006-01-02")
	fileName := strings.Join([]string{
		strconv.FormatInt(timeNow, 10),
		"-",
		strconv.Itoa(width),
		"x",
		strconv.Itoa(height),
		imageExt}, "")
	filePath := path.Join(imageType, currentMonth, fileName)
	return filePath
}

func qnOssImageUpload(c *gin.Context, fileHeader *multipart.FileHeader, filePath string, width int, height int) {
	fdTmp, err := fileHeader.Open()
	if err != nil {
		log.Warnln("上传图片文件打开失败", err)
		mContext{c}.ErrorResponse(ServiceError, ServiceErrorMsg)
		return
	}
	defer fdTmp.Close()

	putPolicy := storage.PutPolicy{
		Scope:      qnImageBucket,
		ReturnBody: `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)"}`,
	}
	mac := qbox.NewMac(qnAccessKey, qnSecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	formUploader := storage.NewFormUploader(&cfg)
	ret := struct {
		Key    string
		Hash   string
		FSize  int
		Bucket string
	}{}
	err = formUploader.Put(context.Background(), &ret, upToken, filePath, fdTmp, fileHeader.Size, nil)
	if err != nil {
		log.Warnln("图片上传到oss失败", err)
		mContext{c}.ErrorResponse(ServiceError, ServiceErrorMsg)
		return
	}

	rsp := imageUploadRes{
		Url:    path.Join(qnImageHost, filePath),
		Width:  width,
		Height: height,
	}
	mContext{c}.SuccessResponse(rsp)
}

func UploadVideo(c *gin.Context) {
	videoType := c.PostForm("VideoType")
	fileHeader, err := c.FormFile("Video")
	if videoType == "" {
		mContext{c}.ErrorResponse(ParamsError, ParamsErrorMsg)
		return
	}
	if err != nil {
		log.Warnln("视频临时文件获取失败", err)
		mContext{c}.ErrorResponse(ParamsDataError, ParamsDataErrorMsg)
		return
	}
	filePath := makeVideoFilePath(videoType, videoExt)
	qnOssVideoUpload(c, fileHeader, filePath)
}

func makeVideoFilePath(videoType string, videoExt string) string {
	timeNow := time.Now().UnixNano()
	currentMonth := time.Now().Format("2006-01-02")
	fileName := strings.Join([]string{
		strconv.FormatInt(timeNow, 10),
		videoExt}, "")
	filePath := path.Join(videoType, currentMonth, fileName)
	return filePath
}

func qnOssVideoUpload(c *gin.Context, fileHeader *multipart.FileHeader, videoPath string) {
	fdTmp, err := fileHeader.Open()
	if err != nil {
		log.Warnln("上传视频打开失败", err)
		mContext{c}.ErrorResponse(ParamsDataError, ParamsDataErrorMsg)
		return
	}
	defer fdTmp.Close()

	//视频帧截图文件  存储设置
	videoThumbPath := path.Join(videoThumbImageType, strings.Replace(videoPath, videoExt, videoThumbExt, 1))
	videoThumbSave := qnImageBucket + ":" + videoThumbPath
	videoThumbPersistentOps := "vframe/jpg/offset/2/w/480/h/360|saveas/" +
		base64.StdEncoding.EncodeToString([]byte(videoThumbSave))

	//视频转码文件	  存储设置
	videoSave := qnFormatVideoBucket + ":" + videoPath
	videoPersistentOps := "avthumb/mp4/r/24/vcodec/libx264|saveas/" +
		base64.StdEncoding.EncodeToString([]byte(videoSave))

	//视频源文件  存储设置
	persistentOps := strings.Join([]string{videoPersistentOps, videoThumbPersistentOps}, ";")
	putPolicy := storage.PutPolicy{
		Scope: qnVideoBucket,
		ReturnBody: `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)",
			"duration":"$(avinfo.format.duration)", "size":"$(avinfo.format.size)"}`,
		PersistentOps: persistentOps,
	}
	mac := qbox.NewMac(qnAccessKey, qnSecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	formUploader := storage.NewFormUploader(&cfg)
	ret := struct {
		Key      string
		Hash     string
		FSize    int
		Bucket   string
		Duration string
	}{}
	err = formUploader.Put(context.Background(), &ret, upToken, videoPath, fdTmp, fileHeader.Size, nil)
	if err != nil {
		log.Warnln("视频上传到oss失败", err)
		mContext{c}.ErrorResponse(ServiceError, ServiceErrorMsg)
		return
	}

	length, _ := strconv.ParseFloat(ret.Duration, 10)
	rsp := videoUploadRes{
		Thumb:  path.Join(qnImageHost, videoThumbPath),
		Length: length,
		Size:   ret.FSize,
		Url:    path.Join(qnVideoHost, videoPath),
	}
	mContext{c}.SuccessResponse(rsp)
}

func UploadAudio(c *gin.Context) {
	audioType := c.PostForm("AudioType")
	fileHeader, err := c.FormFile("Audio")
	if audioType == "" {
		mContext{c}.ErrorResponse(ParamsError, ParamsErrorMsg)
		return
	}
	if err != nil {
		log.Warnln("音频临时文件获取失败", err)
		mContext{c}.ErrorResponse(ParamsDataError, ParamsDataErrorMsg)
		return
	}
	audioFilePath := makeAudioFilePath(audioType, audioExt)
	qnOssAudioUpload(c, fileHeader, audioFilePath)
}

func makeAudioFilePath(audioType string, audioExt string) string {
	timeNow := time.Now().UnixNano()
	currentMonth := time.Now().Format("2006-01-02")
	fileName := strings.Join([]string{
		strconv.FormatInt(timeNow, 10),
		audioExt}, "")
	filePath := path.Join(audioType, currentMonth, fileName)
	return filePath
}

func qnOssAudioUpload(c *gin.Context, fileHeader *multipart.FileHeader, audioPath string) {
	fdTmp, err := fileHeader.Open()
	if err != nil {
		log.Warnln("音频文件打开失败", err)
		mContext{c}.ErrorResponse(ParamsDataError, ParamsDataErrorMsg)
		return
	}
	defer fdTmp.Close()

	putPolicy := storage.PutPolicy{
		Scope: qnAudioBucket,
		ReturnBody: `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)",
			"duration":"$(avinfo.format.duration)", "size":"$(avinfo.format.size)"}`,
	}
	mac := qbox.NewMac(qnAccessKey, qnSecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	formUploader := storage.NewFormUploader(&cfg)
	ret := struct {
		Key      string
		Hash     string
		FSize    int
		Bucket   string
		Duration string
	}{}
	err = formUploader.Put(context.Background(), &ret, upToken, audioPath, fdTmp, fileHeader.Size, nil)
	if err != nil {
		log.Warnln("音频上传到oss失败", err)
		mContext{c}.ErrorResponse(ServiceError, ServiceErrorMsg)
		return
	}

	length, _ := strconv.ParseFloat(ret.Duration, 10)
	rsp := audioUploadRes{
		Url:    path.Join(qnAudioHost, audioPath),
		Length: length,
		Size:   ret.FSize,
	}
	mContext{c}.SuccessResponse(rsp)
}
