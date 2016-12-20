package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"text/template"
)

type Command struct {
	Help *template.Template
	Func func(Context, string, []string) error
}

func newIyashiCommand() Command {
	tmpl := template.Must(template.New("iyashi").Parse(
		"`@{{ .Iyashi.AuthTest.User }} 癒やし <文言>` でflickrから画像を返すよ! 文言はスペース区切りで複数指定可",
	))
	return Command{
		Help: tmpl,
		Func: iyashiFunc,
	}
}

func newHelpCommand() Command {
	tmpl := template.Must(template.New("help").Parse(
		"this message",
	))
	return Command{
		Help: tmpl,
		Func: helpFunc,
	}
}

func iyashiFunc(ctx Context, cmd string, args []string) error {
	if len(args) == 0 {
		args = append(args, "猫")
	}
	args = append(args, "-hentai", "-porn", "-sexy", "-fuck")
	keyword := strings.Join(args, " ")

	query := map[string]string{
		"api_key":        ctx.Iyashi.flickrApiToken,
		"format":         "json",
		"nojsoncallback": "1",
		"method":         "flickr.photos.search",
		"text":           keyword,
		"safe_mode":      "1",
		"media":          "photo",
	}

	res1, err := flickrSearch(query)
	if err != nil {
		return err
	}
	page := rand.Intn(res1.Photos.Pages) + 1

	res2, err := flickrSearch(merge(query, map[string]string{
		"page": fmt.Sprintf("%d", page),
	}))
	if err != nil {
		return err
	}
	if len(res2.Photos.Photo) == 0 {
		return fmt.Errorf("見つかんないよ(´・ω・｀)")
	}

	photo := res2.Photos.Photo[rand.Intn(len(res2.Photos.Photo))]

	imgUrl := fmt.Sprintf(
		"https://farm%d.staticflickr.com/%s/%s_%s.jpg",
		photo.Farm,
		photo.Server,
		photo.Id,
		photo.Secret,
	)

	ctx.ReplyWithoutPermalink(imgUrl)

	return nil
}

func helpFunc(ctx Context, cmd string, args []string) error {
	if 0 < len(args) {
		var doc bytes.Buffer
		if v, ok := ctx.Iyashi.dispatchMap[args[0]]; ok {
			if err := v.Help.Execute(&doc, ctx); err != nil {
				return err
			}
			m := "```\n"
			m += fmt.Sprintf("%s --- %s\n", args[0], doc.String())
			m += "```"
			ctx.Reply(m)
		} else {
			ctx.Reply(fmt.Sprintf("command not found: %s", args[0]))
		}
		return nil
	}

	m := "command list\n```"
	for k, v := range ctx.Iyashi.dispatchMap {
		var doc bytes.Buffer
		if err := v.Help.Execute(&doc, ctx); err != nil {
			return err
		}
		m += fmt.Sprintf("%s --- %s\n", k, doc.String())
	}
	m += "```"
	ctx.Reply(m)
	return nil
}

// flickr
type FlickrSearchResponse struct {
	Photos struct {
		Page    int    `json:"page"`
		Pages   int    `json:"pages"`
		PerPage int    `json:"perpage"`
		Total   string `json:"total"`
		Photo   []struct {
			Id       string `json:"id"`
			Owner    string `json:"owner"`
			Secret   string `json:"secret"`
			Server   string `json:"server"`
			Farm     int    `json:"farm"`
			Title    string `json:"title"`
			Ispublic int    `json:"ispublic"`
			Isfriend int    `json:"isfriend"`
			Isfamily int    `json:"isfamily"`
		} `json:"photo"`
	} `json:"photos"`
}

func flickrSearch(query map[string]string) (FlickrSearchResponse, error) {
	var res FlickrSearchResponse
	resp, err := get(
		"https://api.flickr.com/services/rest/",
		query,
	)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

// tumblr
type TumblrSearchResponse struct {
	Response struct {
		Posts []struct {
			Photos []struct {
				OriginalSize struct {
					Url string `json:"url"`
				} `json:"original_size"`
			} `json:"photos"`
		} `json:"posts"`
		TotalPosts int `json:"total_posts"`
	} `json:"response"`
}

func tumblrSearch(token, tumblrId string, offset int) (TumblrSearchResponse, error) {
	var res TumblrSearchResponse
	url := fmt.Sprintf(
		"http://api.tumblr.com/v2/blog/%s.tumblr.com/posts/photo?api_key=%s&offset=%d",
		tumblrId,
		token,
		offset,
	)
	resp, err := http.Get(url)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return res, fmt.Errorf("failed unmarshal json: %v", err)
	}
	return res, nil
}

func newTumblrCommand(token, tumblrId, command string) Command {
	tmpl := template.Must(template.New(command).Parse(
		fmt.Sprintf("`@{{ .Iyashi.AuthTest.User }} %s` で http://%s.tumblr.com/ から画像をランダムで返すよ", command, tumblrId),
	))
	return Command{
		Help: tmpl,
		Func: func(ctx Context, cmd string, args []string) error {
			res1, err := tumblrSearch(token, tumblrId, 0)
			if err != nil {
				return err
			}

			offset := rand.Intn(res1.Response.TotalPosts/20+1) * 20

			res2, err := tumblrSearch(token, tumblrId, offset)
			if err != nil {
				return err
			}

			urls := []string{}
			for _, post := range res2.Response.Posts {
				for _, photo := range post.Photos {
					urls = append(urls, photo.OriginalSize.Url)
				}
			}

			if len(urls) == 0 {
				return fmt.Errorf("見つかんないよ(´・ω・｀)")
			}

			ctx.ReplyWithoutPermalink(urls[rand.Intn(len(urls))])

			return nil
		},
	}
}

// util
func get(baseUrl string, param map[string]string) (*http.Response, error) {
	queries := []string{}
	for k, v := range param {
		queries = append(queries, fmt.Sprintf("%s=%s", k, url.QueryEscape(v)))
	}
	reqUrl := fmt.Sprintf("%s?%s", baseUrl, strings.Join(queries, "&"))
	return http.Get(reqUrl)
}

func merge(m1, m2 map[string]string, mn ...map[string]string) map[string]string {
	ans := map[string]string{}
	for k, v := range m1 {
		ans[k] = v
	}
	for k, v := range m2 {
		ans[k] = v
	}
	for _, m := range mn {
		for k, v := range m {
			ans[k] = v
		}
	}
	return ans
}
