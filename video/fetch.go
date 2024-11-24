package video

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Video struct {
	Title string
	Link  string
}

const (
	youtubeAPI = "https://www.googleapis.com/youtube/v3/search"
)

var ids = []string{"UC8F_XIeG-FYJ8ezSB352MYg", "UC2HoU9eqcVX7vCeavbLrqGw"}

type LatestFetcher struct {
	key    string
	Ids    []string
	Latest map[string]Video
}

func NewVideoFetcher(key string) LatestFetcher {
	latest := make(map[string]Video)

	return LatestFetcher{
		key:    key,
		Ids:    ids,
		Latest: latest,
	}
}

func (lf *LatestFetcher) FetchLatestVideos() ([]Video, error) {
	var newvideos []Video
	for _, id := range ids {
		vd, err := FetchLatestVideo(id, lf.key)
		if lf.Latest[id].Title != vd.Title {
			newvideos = append(newvideos, vd)
			lf.Latest[id] = vd
		}
		if err != nil {
			return []Video{}, err
		}
	}
	return newvideos, nil
}

func FetchLatestVideo(id, key string) (Video, error) {
	url := fmt.Sprintf("%s?part=snippet&channelId=%s&maxResults=1&order=date&type=video&key=%s",
		youtubeAPI, id, key)
	resp, err := http.Get(url)
	if err != nil {
		return Video{}, err
	}
	defer resp.Body.Close()

	var data struct {
		Items []struct {
			ID struct {
				VideoID string `json:"videoId"`
			} `json:"id"`
			Snippet struct {
				Title string `json:"title"`
			} `json:"snippet"`
		} `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return Video{}, err
	}

	if len(data.Items) == 0 {
		return Video{}, fmt.Errorf("no videos found")
	}

	video := Video{
		Title: data.Items[0].Snippet.Title,
		Link:  fmt.Sprintf("https://www.youtube.com/watch?v=%s", data.Items[0].ID.VideoID),
	}
	return video, nil
}
