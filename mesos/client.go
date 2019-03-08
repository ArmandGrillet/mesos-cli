package mesos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/dcos/dcos-cli/pkg/httpclient"
	"github.com/golang/protobuf/proto"
	"github.com/mesos/mesos-go/api/v1/lib/master"
	"github.com/spf13/afero"
)

// Client is a Mesos client for DC/OS.
type Client struct {
	fs   afero.Fs
	http *httpclient.Client
}

// NewClient creates a new Mesos client.
func NewClient(baseClient *httpclient.Client, ctxFs afero.Fs) *Client {
	return &Client{
		fs:   ctxFs,
		http: baseClient,
	}
}

// Debug returns the path of a task's sandbox.
func (c *Client) Debug(agent string, framework string, executor string, container string) (string, error) {
	resp, err := c.http.Get("/agent/" + agent + "/files/debug")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
		debug := make(map[string]string)
		err = json.NewDecoder(resp.Body).Decode(&debug)
		if err != nil {
			return "", err
		}

		for key, value := range debug {
			if strings.Contains(key, "/frameworks/"+framework+"/executors/"+executor+"/runs/"+container) {
				return value, nil
			}
		}

		return "", fmt.Errorf("unable to find task")
	default:
		return "", fmt.Errorf("HTTP %d error", resp.StatusCode)
	}
}

// Download downloads a task's sandbox in a given directory.
func (c *Client) Download(agent string, agentPath string, downloadDir string) error {
	resp, err := c.http.Get("/agent/" + agent + "/files/browse?path=" + agentPath)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
		files := []File{}
		err = json.NewDecoder(resp.Body).Decode(&files)
		if err != nil {
			return err
		}

		if err := c.fs.MkdirAll(downloadDir, 0755); err != nil {
			return err
		}
		for _, file := range files {
			if strings.HasPrefix(file.Mode, "d") {
				err = c.Download(agent, agentPath+"/"+path.Base(file.Path), filepath.Join(downloadDir, path.Base(file.Path)))
				if err != nil {
					return err
				}
			} else {
				resp, err := c.http.Get("/agent/" + agent + "/files/download?path=" + file.Path)
				if err != nil {
					return err
				}
				defer resp.Body.Close()
				switch resp.StatusCode {
				case 200:
					defer resp.Body.Close()
					out, err := c.fs.Create(filepath.Join(downloadDir, path.Base(file.Path)))
					if err != nil {
						return fmt.Errorf("Unable to create file '%s'", path.Base(file.Path))
					}
					defer out.Close()
					io.Copy(out, resp.Body)
					return nil
				default:
					return fmt.Errorf("HTTP %d error", resp.StatusCode)
				}
			}
		}
		return nil
	default:
		return fmt.Errorf("HTTP %d error", resp.StatusCode)
	}
}

// Sandbox returns the sandbox of a task.
func (c *Client) Sandbox(task string) error {
	resp, err := c.http.Get("/mesos/master/state")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
		return nil
	default:
		return fmt.Errorf("HTTP %d error", resp.StatusCode)
	}
}

// State returns the current State of the Mesos master.
func (c *Client) State() (*master.Response_GetState, error) {
	body := master.Call{
		Type: master.Call_GET_STATE,
	}
	reqBody, err := proto.Marshal(&body)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Post("/mesos/api/v1", "application/x-protobuf", bytes.NewBuffer(reqBody),
		httpclient.Header("Accept", "application/x-protobuf"))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		var state master.Response
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		err = proto.Unmarshal(bodyBytes, &state)
		return state.GetState, err
	case 503:
		return nil, fmt.Errorf("could not connect to the leading mesos master")
	default:
		return nil, fmt.Errorf("HTTP %d error", resp.StatusCode)
	}
}
