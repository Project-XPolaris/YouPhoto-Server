package sdw

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

var DefaultSDWClient *SDWClient

type Conf struct {
	Url string
}
type SDWClient struct {
	Client resty.Client
	Conf   *Conf
}

func NewSDWClient(conf *Conf) *SDWClient {
	return &SDWClient{
		Client: *resty.New(),
		Conf:   conf,
	}
}

type TextToImageParam struct {
	BatchSize         int     `json:"batch_size"`
	Prompt            string  `json:"prompt"`
	NegativePrompt    string  `json:"negative_prompt"`
	Width             int     `json:"width"`
	Height            int     `json:"height"`
	Steps             int     `json:"steps"`
	Seed              int     `json:"seed"`
	CfgScale          int     `json:"cfg_scale"`
	Hrscale           float64 `json:"hr_scale"`
	HrUpscaler        string  `json:"hr_upscaler"`
	HrResizeX         int     `json:"hr_resize_x"`
	HrResizeY         int     `json:"hr_resize_y"`
	HrSecondPassSteps int     `json:"hr_second_pass_steps"`
	DenoisingStrength float64 `json:"denoising_strength"`
	EnableHr          bool    `json:"enable_hr"`
	Niter             int     `json:"n_iter"`
}

func (c *SDWClient) TextToImage(param *TextToImageParam) (*Text2ImageResponse, error) {
	var result Text2ImageResponse
	_, err := c.Client.NewRequest().
		SetBody(param).
		SetResult(&result).
		Post(fmt.Sprintf("%s/sdapi/v1/txt2img", c.Conf.Url))
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *SDWClient) UpdateOption(param map[string]interface{}) error {
	_, err := c.Client.NewRequest().
		SetBody(param).
		Post(fmt.Sprintf("%s/sdapi/v1/options", c.Conf.Url))
	if err != nil {
		return err
	}
	return nil
}
func (c *SDWClient) GetOption() (*Options, error) {
	var result Options
	_, err := c.Client.NewRequest().
		SetResult(&result).
		Get(fmt.Sprintf("%s/sdapi/v1/options", c.Conf.Url))
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *SDWClient) GetModels() ([]*ModelInfo, error) {
	var result []*ModelInfo
	_, err := c.Client.NewRequest().
		SetResult(&result).
		Get(fmt.Sprintf("%s/sdapi/v1/sd-models", c.Conf.Url))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *SDWClient) GetSamplers() ([]*Sampler, error) {
	var result []*Sampler
	_, err := c.Client.NewRequest().
		SetResult(&result).
		Get(fmt.Sprintf("%s/sdapi/v1/samplers", c.Conf.Url))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *SDWClient) GetUpscaler() ([]*Upscaler, error) {
	var result []*Upscaler
	_, err := c.Client.NewRequest().
		SetResult(&result).
		Get(fmt.Sprintf("%s/sdapi/v1/upscalers", c.Conf.Url))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *SDWClient) GetProgress() (*Progress, error) {
	var result Progress
	_, err := c.Client.NewRequest().
		SetResult(&result).
		Get(fmt.Sprintf("%s/sdapi/v1/progress", c.Conf.Url))
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *SDWClient) Interrupt() error {
	_, err := c.Client.NewRequest().
		Post(fmt.Sprintf("%s/sdapi/v1/interrupt", c.Conf.Url))
	if err != nil {
		return err
	}
	return nil
}
func (c *SDWClient) Skip() error {
	_, err := c.Client.NewRequest().
		Post(fmt.Sprintf("%s/sdapi/v1/skip", c.Conf.Url))
	if err != nil {
		return err
	}
	return nil
}

const (
	PreprocessTxtActionPrepend = "prepend"
	PreprocessTxtActionAppend  = "append"
	PreprocessTxtActionCopy    = "copy"
)

type PreprocessParam struct {
	ProcessSrc              string `json:"process_src"`
	ProcessDst              string `json:"process_dst"`
	ProcessWidth            int    `json:"process_width"`
	ProcessHeight           int    `json:"process_height"`
	IdTask                  string `json:"id_task"`
	PreprocessTxtAction     string `json:"preprocess_txt_action"`
	ProcessFlip             bool   `json:"process_flip"`
	ProcessSplit            bool   `json:"process_split"`
	ProcessCaption          bool   `json:"process_caption"`
	ProcessCaptionDeepbooru bool   `json:"process_caption_deepbooru"`
}

func (c *SDWClient) Preprocess(param *PreprocessParam) error {
	_, err := c.Client.NewRequest().
		SetBody(param).
		Post(fmt.Sprintf("%s/sdapi/v1/preprocess", c.Conf.Url))
	if err != nil {
		return err
	}
	return nil
}
