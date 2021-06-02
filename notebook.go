package jupyter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

type OutputData struct {
	TextHTML       []string `json:"text/html,omitempty"`
	ApplicationPDF *string  `json:"application/pdf,omitempty"`
	TextLaTeX      *string  `json:"text/latex,omitempty"`
	ImageSVGXML    []string `json:"image/svg+xml,omitempty"`
	ImagePNG       *string  `json:"image/png,omitempty"`
	ImageJPEG      *string  `json:"image/jpeg,omitempty"`
	TextMarkdown   []string `json:"text/markdown,omitempty"`
	TextPlain      []string `json:"text/plain,omitempty"`
}

type Output struct {
	OutputType     string     `json:"output_type"`
	ExecutionCount *int       `json:"execution_count,omitempty"`
	Text           []string   `json:"text,omitempty"`
	Data           OutputData `json:"data,omitempty"`
	Traceback      []string   `json:"traceback,omitempty"`
}

type Cell struct {
	CellType       string `json:"cell_type"`
	ExecutionCount *int   `json:"execution_count,omitempty"`
	ID             string `json:"id"`
	//Metadata       struct {
	//Trusted bool `json:"trusted"`
	//} `json:"metadata"`
	Metadata interface{}              `json:"metadata,omitempty"`
	Outputs  []map[string]interface{} `json:"outputs,omitempty"`
	Source   interface{}              `json:"source,omitempty"`
}

type Metadata struct {
	Kernelspec struct {
		DisplayName string `json:"display_name"`
		Language    string `json:"language"`
		Name        string `json:"name"`
	} `json:"kernelspec"`
	LanguageInfo struct {
		FileExtension     string `json:"file_extension"`
		Mimetype          string `json:"mimetype"`
		Name              string `json:"name"`
		NbconvertExporter string `json:"nbconvert_exporter"`
		PygmentsLexer     string `json:"pygments_lexer"`
		Version           string `json:"version"`
	} `json:"language_info"`
}

type DocumentContent struct {
	Cells         []Cell   `json:"cells"`
	Metadata      Metadata `json:"metadata,omitempty"`
	Nbformat      int      `json:"nbformat"`
	NbformatMinor int      `json:"nbformat_minor"`
}

type Document struct {
	Name         string          `json:"name"`
	Path         string          `json:"path"`
	LastModified time.Time       `json:"last_modified"`
	Created      time.Time       `json:"created"`
	Content      DocumentContent `json:"content"`
	Format       string          `json:"format"`
	Mimetype     interface{}     `json:"mimetype"`
	Size         int             `json:"size"`
	Writable     bool            `json:"writable"`
	Type         string          `json:"type"`
}

type Kernel struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	LastActivity   time.Time `json:"last_activity"`
	ExecutionState string    `json:"execution_state"`
	Connections    int       `json:"connections"`
}

type MessageContent struct {
	Code   interface{} `json:"code"`
	Silent bool        `json:"silent"`
}

type MessageHeader struct {
	MsgID    string `json:"msg_id"`
	Username string `json:"username,omitempty"`
	Session  string `json:"session"`
	MsgType  string `json:"msg_type"`
	Version  string `json:"version,omitempty"`
}

type Message struct {
	Header       MessageHeader          `json:"header"`
	ParentHeader MessageHeader          `json:"parent_header"`
	Metadata     map[string]interface{} `json:"metadata"`
	Content      MessageContent         `json:"content,omitempty"`
}

type ExecContent struct {
	Text           string `json:"text,omitempty"`
	Ename          string `json:"ename,omitempty"`
	Evalue         string `json:"evalue,omitempty"`
	ExecutionState string `json:"execution_state,omitempty"`
	ExecutionCount *int   `json:"execution_count,omitempty"`
}

type ExecResult struct {
	MsgID    string      `json:"msg_id"`
	Session  string      `json:"session"`
	Username string      `json:"username,omitempty"`
	Date     string      `json:"date,omitempty"`
	MsgType  string      `json:"msg_type"`
	Version  string      `json:"version,omitempty"`
	Content  ExecContent `json:"content,omitempty"`
}

type NotebookOp struct {
	c *Client
	p string
	d Document
	s string
	k *Kernel
}

func (c *Client) Notebook(p string) (*NotebookOp, error) {
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

	var d Document
	err = json.NewDecoder(resp.Body).Decode(&d)
	if err != nil {
		return nil, err
	}
	return &NotebookOp{c: c, p: p, d: d, s: uuid.NewString(), k: nil}, nil
}

func (c *NotebookOp) Update(in io.Reader) error {
	var d Document
	err := json.NewDecoder(in).Decode(&d)
	if err != nil {
		return err
	}
	c.d = d
	return nil
}

func (c *NotebookOp) Doc() *Document {
	return &c.d
}

func (c *NotebookOp) CodeIDs() []string {
	ids := []string{}
	for _, cell := range c.d.Content.Cells {
		ids = append(ids, cell.ID)
	}
	return ids
}

func (c *NotebookOp) Code(id string) (string, error) {
	for _, cell := range c.d.Content.Cells {
		if cell.ID == id {
			return cell.Source.(string), nil
		}
	}
	return "", nil
}

func (c *NotebookOp) Exec(id string, stdout io.Writer, stderr io.Writer) error {
	u, err := url.Parse(c.c.config.Origin)
	if err != nil {
		return err
	}

	base := u.Path

	if c.k == nil {
		u.Path = path.Join(base, "api", "kernels")
		req, err := http.NewRequest(http.MethodPost, u.String(), nil)
		if err != nil {
			return err
		}
		req.Header.Add("Authorization", "token "+c.c.config.Token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		req, err = http.NewRequest(http.MethodGet, u.String(), nil)
		if err != nil {
			return err
		}
		req.Header.Add("Authorization", "token "+c.c.config.Token)
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var kernels []Kernel
		err = json.NewDecoder(resp.Body).Decode(&kernels)
		if err != nil {
			return err
		}
		for _, kernel := range kernels {
			if kernel.Name == c.d.Content.Metadata.Kernelspec.Name {
				c.k = &kernel
				break
			}
		}
		if c.k == nil {
			return errors.New("kernel not found")
		}
	}

	u.Scheme = "ws"
	u.Path = path.Join(base, "api", "kernels", c.k.ID, "channels")
	config, err := websocket.NewConfig(u.String(), c.c.config.Origin)
	if err != nil {
		return err
	}
	config.Header.Add("Authorization", "token "+c.c.config.Token)
	ws, err := websocket.DialConfig(config)
	if err != nil {
		return err
	}
	defer ws.Close()

	for _, code := range c.d.Content.Cells {
		if code.Source == nil {
			continue
		}
		if s, ok := code.Source.(string); ok && s == "" {
			continue
		}
		if code.ID != id {
			continue
		}
		var header MessageHeader
		header.MsgType = "execute_request"
		header.MsgID = uuid.NewString()
		header.Session = c.s
		payload := Message{
			Header:       header,
			ParentHeader: header,
			Metadata:     make(map[string]interface{}),
			Content: MessageContent{
				Code:   code.Source,
				Silent: false,
			},
		}

		err = websocket.JSON.Send(ws, payload)
		if err != nil {
			return err
		}

		var result ExecResult
		var input bool
		for {
			err := websocket.JSON.Receive(ws, &result)
			if err != nil {
				break
			}
			if result.MsgType == "stream" {
				_, err = fmt.Fprint(stdout, result.Content.Text)
				if err != nil {
					return err
				}
			}
			if result.MsgType == "error" {
				return fmt.Errorf("%s: %s", result.Content.Ename, result.Content.Evalue)
			}
			if result.MsgType == "status" {
				if input && result.Content.ExecutionState == "idle" {
					break
				}
			}
			if result.MsgType == "execute_input" {
				input = true
			}
		}
	}

	return nil
}
