package music

import (
	"bufio"
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-musicfox/netease-music/service"
	"github.com/h2non/bimg"
	"go.yaml.in/yaml/v4"

	"github.com/syhily/pandora/internal/config"
	"github.com/syhily/pandora/internal/http"
	"github.com/syhily/pandora/internal/s3"
)

// Processor handles music processing operations
type Processor struct {
	config     *config.PandoraConfig
	s3Client   *s3.Client
	httpClient *http.Client
}

// NewProcessor creates a new music processor
func NewProcessor(cfg *config.PandoraConfig) *Processor {
	return &Processor{
		config:     cfg,
		s3Client:   s3.NewClient(&cfg.MusicS3),
		httpClient: http.NewClient(),
	}
}

// Process processes a music ID and downloads/uploads it
func (p *Processor) Process(musicID int, useVIP bool) error {
	// Resolve the MUSIC information
	id := strconv.FormatInt(int64(musicID), 10)
	metadata := &NeteaseMusic{ID: id}

	// Try to resolve the song url
	var resp []byte
	var err error
	if useVIP {
		resp, err = p.getVIPURL(id)
		if err != nil {
			return fmt.Errorf("failed to get VIP URL: %w", err)
		}
	} else {
		urlService := service.SongUrlV1Service{ID: id, Level: service.Standard, SkipUNM: true}
		_, resp, err = urlService.SongUrl()
		if err != nil {
			return fmt.Errorf("failed to resolve song URL: %w", err)
		}
	}

	decode := json.NewDecoder(bytes.NewReader(resp))
	songUrl := &SongUrl{}
	if err := decode.Decode(songUrl); err != nil {
		return fmt.Errorf("failed to decode song URL response: %w", err)
	}

	if len(songUrl.Data) == 0 {
		return fmt.Errorf("failed to resolve the song url %s", id)
	}

	// Download the song and push to S3
	content, err := p.httpClient.DownloadFile(songUrl.Data[0].URL)
	if err != nil {
		return fmt.Errorf("failed to download song: %w", err)
	}

	musicKey := "musics/" + id + ".mp3"
	if err := p.s3Client.UploadObject(context.TODO(), musicKey, content); err != nil {
		return fmt.Errorf("failed to upload music: %w", err)
	}

	metadata.URL, _ = url.JoinPath(p.config.MusicS3.PublicDomain, musicKey)
	log.Println("Successfully upload the music", id)

	// Try to resolve the song details
	detailService := service.SongDetailService{Ids: id}
	_, resp = detailService.SongDetail()
	decode = json.NewDecoder(bytes.NewReader(resp))
	songDetail := &SongDetail{}
	if err := decode.Decode(songDetail); err != nil {
		return fmt.Errorf("failed to decode song detail response: %w", err)
	}

	if len(songDetail.Songs) == 0 {
		return fmt.Errorf("failed to resolve the song details %s", id)
	}

	metadata.Name = songDetail.Songs[0].Name
	metadata.Artist = songDetail.Songs[0].Ar[0].Name
	metadata.Album = songDetail.Songs[0].Al.Name

	// Download the album pic and push to S3
	picURL := songDetail.Songs[0].Al.PicURL
	pic, err := p.httpClient.DownloadFile(picURL)
	if err != nil {
		return fmt.Errorf("failed to download album pic: %w", err)
	}

	image := bimg.NewImage(pic)
	options := bimg.Options{
		Width:   300,
		Height:  300,
		Crop:    false,
		Quality: 100,
		Rotate:  0,
		Type:    bimg.JPEG,
	}
	pic, err = image.Process(options)
	if err != nil {
		return fmt.Errorf("failed to convert the images: %w", err)
	}

	picKey := "musics/" + id + ".jpg"
	if err := p.s3Client.UploadObject(context.TODO(), picKey, pic); err != nil {
		return fmt.Errorf("failed to upload album pic: %w", err)
	}

	metadata.Pic, _ = url.JoinPath(p.config.MusicS3.PublicDomain, picKey)
	log.Println("Successfully upload the album pic", id)

	// Try to resolve the song lyric
	lyricService := service.LyricService{ID: id}
	_, resp = lyricService.Lyric()
	decode = json.NewDecoder(bytes.NewReader(resp))
	lyric := &Lyric{}
	if err := decode.Decode(lyric); err != nil {
		return fmt.Errorf("failed to decode lyric response: %w", err)
	}

	metadata.Lyric = cmp.Or(lyric.Tlyric.Lyric, lyric.Klyric.Lyric, lyric.Lrc.Lyric, "[00:00.00]无歌词")

	// Save music metadata file into blog project.
	filename := filepath.Join(p.config.BlogRoot, "src", "content", "metas", "musics", id+".yml")
	if err := os.MkdirAll(filepath.Dir(filename), os.FileMode(0755)); err != nil {
		return fmt.Errorf("failed to create metadata directory: %w", err)
	}

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0644))
	if err != nil {
		return fmt.Errorf("failed to generate the metadata file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	encoder := yaml.NewEncoder(writer)
	encoder.SetIndent(2)
	if err := encoder.Encode(metadata); err != nil {
		return fmt.Errorf("failed to save image metadata: %w", err)
	}

	log.Println("Save the music metadata", id)
	return nil
}

// getVIPURL gets the VIP URL for a song
func (p *Processor) getVIPURL(id string) ([]byte, error) {
	url := "https://wyapi.toubiec.cn/api/music/url"
	payload := []byte(fmt.Sprintf(`{"id":"%s","level":"standard"}`, id))

	headers := map[string]string{
		"accept":             "*/*",
		"accept-language":    "en,zh;q=0.9,en-US;q=0.8,zh-CN;q=0.7",
		"cache-control":      "no-cache",
		"content-type":       "application/json",
		"cookie":             "server_name_session=8642b07d17bffff99272d773c8fbaf3d; hasSeenWelcome=true",
		"dnt":                "1",
		"origin":             "https://wyapi.toubiec.cn",
		"pragma":             "no-cache",
		"priority":           "u=1, i",
		"referer":            "https://wyapi.toubiec.cn/",
		"sec-ch-ua":          `"Chromium";v="142", "Google Chrome";v="142", "Not_A Brand";v="99"`,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": `"macOS"`,
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-origin",
		"user-agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36",
	}

	return p.httpClient.PostJSON(url, payload, headers)
}
