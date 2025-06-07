package panel

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/InazumaV/V2bX/conf"
	"github.com/go-resty/resty/v2"
)

// Panel is the interface for different panel's api.
type Panel interface {
	GetUserList() ([]UserInfo, error)
	GetUserAlive() (map[int]int, error)
	ReportUserTraffic([]UserTraffic) error
	ReportNodeOnlineUsers(*map[int][]string) error
	GetNodeInfo() (*NodeInfo, error)
	GetAPIHost() string
	GetNodeType() string
	GetNodeId() int
}

type Client struct {
	client           *resty.Client
	APIHost          string
	Token            string
	NodeType         string
	NodeId           int
	nodeEtag         string
	userEtag         string
	responseBodyHash string
	UserList         *UserListBody
	AliveMap         *AliveMap
}

func New(c *conf.ApiConfig) (*Client, error) {
	client := resty.New()
	client.SetRetryCount(3)
	if c.Timeout > 0 {
		client.SetTimeout(time.Duration(c.Timeout) * time.Second)
	} else {
		client.SetTimeout(5 * time.Second)
	}
	client.OnError(func(req *resty.Request, err error) {
		var v *resty.ResponseError
		if errors.As(err, &v) {
			// v.Response contains the last response from the server
			// v.Err contains the original error
			logrus.Error(v.Err)
		}
	})
	client.SetBaseURL(c.APIHost)
	// Check node type
	c.NodeType = strings.ToLower(c.NodeType)
	switch c.NodeType {
	case "v2ray":
		c.NodeType = "vmess"
	case
		"vmess",
		"trojan",
		"shadowsocks",
		"hysteria",
		"hysteria2",
		"tuic",
		"anytls",
		"vless":
	default:
		return nil, fmt.Errorf("unsupported Node type: %s", c.NodeType)
	}
	// set params
	client.SetQueryParams(map[string]string{
		"node_type": c.NodeType,
		"node_id":   strconv.Itoa(c.NodeID),
		"token":     c.Key,
	})
	return &Client{
		client:   client,
		Token:    c.Key,
		APIHost:  c.APIHost,
		NodeType: c.NodeType,
		NodeId:   c.NodeID,
		UserList: &UserListBody{},
		AliveMap: &AliveMap{},
	}, nil
}

// NewPanel 工厂方法，根据配置创建不同面板实现
func NewPanel(c *conf.ApiConfig) (Panel, error) {
	switch c.PanelType {
	case "v2board":
		return New(c)
	// 预留：后续可添加更多面板类型
	default:
		return nil, fmt.Errorf("unsupported panel type: %s", c.PanelType)
	}
}

func (c *Client) GetAPIHost() string  { return c.APIHost }
func (c *Client) GetNodeType() string { return c.NodeType }
func (c *Client) GetNodeId() int      { return c.NodeId }
