package lora

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
)

var DefaultLoraTrainClient *LoraTrainClient

type Conf struct {
	Url string
}
type LoraTrainClient struct {
	Client resty.Client
	Conf   *Conf
}

func NewLoraTrainClient(conf *Conf) *LoraTrainClient {
	return &LoraTrainClient{
		Client: *resty.New(),
		Conf:   conf,
	}
}

func (c *LoraTrainClient) FetchInfo() (*InfoResponse, error) {
	var result InfoResponse
	_, err := c.Client.NewRequest().
		SetBody(&result).
		SetResult(&result).
		Get(fmt.Sprintf("%s/info", c.Conf.Url))
	if err != nil {
		return nil, err
	}
	return &result, nil

}
func (c *LoraTrainClient) FetchTasks() (*[]TrainTask, error) {
	var result BaseResponse[[]TrainTask]
	_, err := c.Client.NewRequest().
		SetResult(&result).
		Get(fmt.Sprintf("%s/tasks", c.Conf.Url))
	if err != nil {
		return nil, err
	}
	if result.Success == false {
		return nil, errors.New(result.Err)
	}
	return &result.Data, nil
}

func (c *LoraTrainClient) FetchTask(taskId string) (*TrainTask, error) {
	var result BaseResponse[TrainTask]
	_, err := c.Client.NewRequest().
		SetResult(&result).
		Get(fmt.Sprintf("%s/task/%s", c.Conf.Url, taskId))
	if err != nil {
		return nil, err
	}
	if result.Success == false {
		return nil, errors.New(result.Err)
	}
	return &result.Data, nil
}

func (c *LoraTrainClient) InterruptTask(taskId string) error {
	var result BaseResponse[string]
	_, err := c.Client.NewRequest().
		SetResult(&result).
		Get(fmt.Sprintf("%s/task/%s/interrupt", c.Conf.Url, taskId))
	if err != nil {
		return err
	}
	if result.Success == false {
		return errors.New(result.Err)
	}
	return nil
}

func (c *LoraTrainClient) Train(Param *TrainConfigValues) (*TranResponse, error) {
	var result BaseResponse[TranResponse]
	_, err := c.Client.NewRequest().
		SetBody(Param).
		SetResult(&result).
		Post(fmt.Sprintf("%s/train", c.Conf.Url))
	if err != nil {
		return nil, err
	}
	if result.Success == false {
		return nil, errors.New(result.Err)
	}
	return &result.Data, nil
}
