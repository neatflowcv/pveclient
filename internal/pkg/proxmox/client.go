package proxmox

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

type Client struct {
	baseURL   string
	apiToken  string
	requester *Requester
}

type VersionResponse struct {
	Data struct {
		Version string `json:"version"`
		Release string `json:"release"`
		RepoBit string `json:"repoid"`
	} `json:"data"`
}

func NewClient(baseURL string, apiToken string, opts ...ConfigOption) *Client {
	var config Config
	for _, opt := range opts {
		opt(&config)
	}

	return &Client{
		baseURL:   baseURL,
		apiToken:  apiToken,
		requester: NewRequester(config.insecureSkipTLS),
	}
}

var ErrInvalidStatusCode = errors.New("invalid status code")

func (c *Client) Version(ctx context.Context) (string, error) {
	endpoint, err := url.JoinPath(c.baseURL, "/api2/json/version")
	if err != nil {
		return "", fmt.Errorf("failed to construct URL: %w", err)
	}

	req := NewGetRequest(ctx, endpoint, map[string][]string{
		"Authorization": {c.apiToken},
	})

	statusCode, content, err := c.requester.Call(req)
	if err != nil {
		return "", err
	}

	if statusCode != StatusCodeOK {
		return "", fmt.Errorf("%w: %d", ErrInvalidStatusCode, statusCode)
	}

	var versionResp VersionResponse

	err = json.Unmarshal(content, &versionResp)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return versionResp.Data.Version, nil
}

type ListNodesResponse struct {
	Data []struct {
		Status         string  `json:"status"`
		Maxcpu         int     `json:"maxcpu"`
		Mem            int64   `json:"mem"`
		CPU            float64 `json:"cpu"`
		Level          string  `json:"level"`
		Maxdisk        int64   `json:"maxdisk"`
		ID             string  `json:"id"`
		Maxmem         int64   `json:"maxmem"`
		Disk           int64   `json:"disk"`
		Type           string  `json:"type"`
		SslFingerprint string  `json:"ssl_fingerprint"`
		Node           string  `json:"node"`
		Uptime         int     `json:"uptime"`
	} `json:"data"`
}

func (c *Client) ListNodes(ctx context.Context) (*ListNodesResponse, error) {
	endpoint, err := url.JoinPath(c.baseURL, "/api2/json/nodes")
	if err != nil {
		return nil, fmt.Errorf("failed to construct URL: %w", err)
	}

	req := NewGetRequest(ctx, endpoint, map[string][]string{
		"Authorization": {c.apiToken},
	})

	statusCode, content, err := c.requester.Call(req)
	if err != nil {
		return nil, err
	}

	if statusCode != StatusCodeOK {
		return nil, fmt.Errorf("%w: %d", ErrInvalidStatusCode, statusCode)
	}

	var ret ListNodesResponse

	err = json.Unmarshal(content, &ret)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &ret, nil
}

type DiskType string

const (
	DiskTypeHDD DiskType = "hdd"
	DiskTypeSSD DiskType = "ssd"
)

type ListDisksResponse struct {
	Data []struct {
		Size         int64       `json:"size"`
		OsdidList    interface{} `json:"osdid-list"`
		Osdid        IntValue    `json:"osdid"`
		Used         string      `json:"used,omitempty"`
		Wwn          string      `json:"wwn"`
		Health       string      `json:"health"`
		Rpm          IntValue    `json:"rpm"`
		Gpt          int         `json:"gpt"`
		Type         DiskType    `json:"type"`
		ByIDLink     string      `json:"by_id_link"`
		Serial       string      `json:"serial"`
		Devpath      string      `json:"devpath"`
		Wearout      Wearout     `json:"wearout"`
		Model        string      `json:"model"`
		Vendor       string      `json:"vendor"`
		Db           int         `json:"db,omitempty"`
		Bluestore    int         `json:"bluestore,omitempty"`
		Osdencrypted int         `json:"osdencrypted,omitempty"`
	} `json:"data"`
}

func (c *Client) ListDisks(ctx context.Context, node string) (*ListDisksResponse, error) {
	endpoint, err := url.JoinPath(c.baseURL, "/api2/json/nodes", node, "disks/list")
	if err != nil {
		return nil, fmt.Errorf("failed to construct URL: %w", err)
	}

	req := NewGetRequest(ctx, endpoint, map[string][]string{
		"Authorization": {c.apiToken},
	})

	statusCode, content, err := c.requester.Call(req)
	if err != nil {
		return nil, err
	}

	if statusCode != StatusCodeOK {
		return nil, fmt.Errorf("%w: %d", ErrInvalidStatusCode, statusCode)
	}

	var ret ListDisksResponse

	err = json.Unmarshal(content, &ret)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &ret, nil
}

func (c *Client) IssueTicket(ctx context.Context, realm, username, password string) (*IssueTicketResponse, error) {
	endpoint, err := url.JoinPath(c.baseURL, "/api2/json/access/ticket")
	if err != nil {
		return nil, fmt.Errorf("failed to construct URL: %w", err)
	}

	req := NewPostRequest(ctx, endpoint, map[string][]string{
		"Content-Type": {"application/x-www-form-urlencoded"},
	}, []byte(fmt.Sprintf("username=%s@%s&password=%s", username, realm, password)))

	statusCode, content, err := c.requester.Call(req)
	if err != nil {
		return nil, err
	}

	if statusCode != StatusCodeOK {
		return nil, fmt.Errorf("%w: %d", ErrInvalidStatusCode, statusCode)
	}

	var ret IssueTicketResponse

	err = json.Unmarshal(content, &ret)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &ret, nil
}

type IssueTicketResponse struct {
	Data struct {
		Ticket   string `json:"ticket"`
		Username string `json:"username"`
		Cap      struct {
			Mapping struct {
				MappingAudit      int `json:"Mapping.Audit"`
				MappingUse        int `json:"Mapping.Use"`
				MappingModify     int `json:"Mapping.Modify"`
				PermissionsModify int `json:"Permissions.Modify"`
			} `json:"mapping"`
			Nodes struct {
				SysPowerMgmt      int `json:"Sys.PowerMgmt"`
				SysModify         int `json:"Sys.Modify"`
				SysAudit          int `json:"Sys.Audit"`
				PermissionsModify int `json:"Permissions.Modify"`
				SysSyslog         int `json:"Sys.Syslog"`
				SysAccessNetwork  int `json:"Sys.AccessNetwork"`
				SysConsole        int `json:"Sys.Console"`
				SysIncoming       int `json:"Sys.Incoming"`
			} `json:"nodes"`
			Storage struct {
				PermissionsModify         int `json:"Permissions.Modify"`
				DatastoreAllocate         int `json:"Datastore.Allocate"`
				DatastoreAllocateTemplate int `json:"Datastore.AllocateTemplate"`
				DatastoreAudit            int `json:"Datastore.Audit"`
				DatastoreAllocateSpace    int `json:"Datastore.AllocateSpace"`
			} `json:"storage"`
			Sdn struct {
				PermissionsModify int `json:"Permissions.Modify"`
				SDNUse            int `json:"SDN.Use"`
				SDNAudit          int `json:"SDN.Audit"`
				SDNAllocate       int `json:"SDN.Allocate"`
			} `json:"sdn"`
			Dc struct {
				SysAudit    int `json:"Sys.Audit"`
				SDNAudit    int `json:"SDN.Audit"`
				SDNUse      int `json:"SDN.Use"`
				SysModify   int `json:"Sys.Modify"`
				SDNAllocate int `json:"SDN.Allocate"`
			} `json:"dc"`
			Access struct {
				UserModify        int `json:"User.Modify"`
				PermissionsModify int `json:"Permissions.Modify"`
				GroupAllocate     int `json:"Group.Allocate"`
			} `json:"access"`
			Vms struct {
				VMSnapshot         int `json:"VM.Snapshot"`
				PermissionsModify  int `json:"Permissions.Modify"`
				VMAllocate         int `json:"VM.Allocate"`
				VMConsole          int `json:"VM.Console"`
				VMConfigCPU        int `json:"VM.Config.CPU"`
				VMPowerMgmt        int `json:"VM.PowerMgmt"`
				VMConfigCloudinit  int `json:"VM.Config.Cloudinit"`
				VMConfigNetwork    int `json:"VM.Config.Network"`
				VMConfigDisk       int `json:"VM.Config.Disk"`
				VMConfigMemory     int `json:"VM.Config.Memory"`
				VMSnapshotRollback int `json:"VM.Snapshot.Rollback"`
				VMConfigCDROM      int `json:"VM.Config.CDROM"`
				VMConfigHWType     int `json:"VM.Config.HWType"`
				VMBackup           int `json:"VM.Backup"`
				VMMonitor          int `json:"VM.Monitor"`
				VMConfigOptions    int `json:"VM.Config.Options"`
				VMAudit            int `json:"VM.Audit"`
				VMClone            int `json:"VM.Clone"`
				VMMigrate          int `json:"VM.Migrate"`
			} `json:"vms"`
		} `json:"cap"`
		CSRFPreventionToken string `json:"CSRFPreventionToken"`
		Clustername         string `json:"clustername"`
	} `json:"data"`
}

type Wearout struct {
	IsAvailable bool
	Value       int
}

func (w *Wearout) UnmarshalJSON(data []byte) error {
	if string(data) == "\"N/A\"" {
		w.IsAvailable = false
		w.Value = 0

		return nil
	}

	value, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse wearout: %w", err)
	}

	w.IsAvailable = true
	w.Value = int(value)

	return nil
}

type IntValue struct {
	Value int
}

func (i *IntValue) UnmarshalJSON(data []byte) error {
	// -1인 경우 처리
	if string(data) == "-1" {
		i.Value = -1

		return nil
	}

	// 따옴표로 둘러싸인 문자열인 경우 처리
	if len(data) >= 2 && data[0] == '"' && data[len(data)-1] == '"' {
		// 따옴표 제거
		str := string(data[1 : len(data)-1])

		value, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse osdid from quoted string: %w", err)
		}

		i.Value = int(value)

		return nil
	}

	// 일반 숫자인 경우 처리
	value, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse osdid: %w", err)
	}

	i.Value = int(value)

	return nil
}
