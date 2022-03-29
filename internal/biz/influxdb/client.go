package influxdb

import (
	"context"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/domain"
	"time"
)

type Client struct {
	influxdb2.Client
	org   string
	orgID string
}

func NewInfluxdbClient(serverUrl, authToken, orgName string) (*Client, error) {
	client := influxdb2.NewClient(serverUrl, authToken)
	organization, err := client.OrganizationsAPI().FindOrganizationByName(
		context.Background(), orgName)
	if err != nil {
		return nil, err
	}
	return &Client{
		Client: client,
		org:    orgName,
		orgID:  *organization.Id,
	}, nil
}

// CreateBucket 为用户创建保存设备状态信息、保存下采样数据、保存警告信息的三个bucket
func (c *Client) CreateBucket(username string) error {
	bucketsAPI := c.Client.BucketsAPI()
	buckets := []string{
		username,
		fmt.Sprintf("%s-warning_detect", username),
		fmt.Sprintf("%s-warnings", username),
	}

	var err error
	defer func() {
		if err != nil {
			c.ClearBucket(username)
		}
	}()

	for _, bucket := range buckets {
		_, err = bucketsAPI.CreateBucketWithNameWithID(
			context.Background(),
			c.orgID,
			bucket,
			domain.RetentionRule{
				EverySeconds: int64(720 * time.Hour.Seconds()),
				Type:         "expire",
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// ClearBucket 删除用户相关的bucket
func (c *Client) ClearBucket(username string) error {
	bucketsAPI := c.Client.BucketsAPI()
	buckets := []string{
		username,
		fmt.Sprintf("%s-warning_detect", username),
		fmt.Sprintf("%s-warnings", username),
	}

	for _, bucket := range buckets {
		b, err := bucketsAPI.FindBucketByName(context.Background(), bucket)
		if err == nil {
			bucketsAPI.DeleteBucket(context.Background(), b)
		}
	}
	return nil
}
