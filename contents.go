package jupyter

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"
)

type ContentsContent struct {
	Name         string      `json:"name"`
	Path         string      `json:"path"`
	LastModified time.Time   `json:"last_modified"`
	Created      time.Time   `json:"created"`
	Content      interface{} `json:"content"`
	Format       interface{} `json:"format"`
	Mimetype     interface{} `json:"mimetype"`
	Size         int         `json:"size"`
	Writable     bool        `json:"writable"`
	Type         string      `json:"type"`
}

type Contents struct {
	Name         string            `json:"name"`
	Path         string            `json:"path"`
	LastModified time.Time         `json:"last_modified"`
	Created      time.Time         `json:"created"`
	Content      []ContentsContent `json:"content"`
	Format       string            `json:"format"`
	Mimetype     interface{}       `json:"mimetype"`
	Size         interface{}       `json:"size"`
	Writable     bool              `json:"writable"`
	Type         string            `json:"type"`
}

func (c *Client) List(p string) ([]ContentsContent, error) {
	u, err := url.Parse(c.config.Origin)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, "api", "contents", p)
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "token "+c.config.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result Contents
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result.Content, nil
}

func (c *Client) Cat(p string, w io.Writer) error {
	u, err := url.Parse(c.config.Origin)
	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, "api", "contents", p)
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "token "+c.config.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(w, resp.Body)
	return err
}

func (c *Client) Get(p string) ([]byte, error) {
	u, err := url.Parse(c.config.Origin)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, "api", "contents", p)
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "token "+c.config.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (c *Client) Save(p string, r io.Reader) error {
	u, err := url.Parse(c.config.Origin)
	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, "api", "contents", p)
	req, err := http.NewRequest(http.MethodPut, u.String(), r)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "token "+c.config.Token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		var errmsg Error
		err = json.NewDecoder(resp.Body).Decode(&errmsg)
		if err != nil {
			return err
		}
		return errors.New(errmsg.Message)
	}
	return nil
}
