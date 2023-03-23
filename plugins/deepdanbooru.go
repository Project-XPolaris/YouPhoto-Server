package plugins

import (
	"github.com/allentom/harukap/plugins/deepdanbooru"
	"github.com/project-xpolaris/youplustoolkit/youlog"
	"os"
)

var DefaultDeepDanbooruPlugin = deepdanbooru.NewPlugin()

var DefaultDeepdanbooruLauncher = NewDeepdanbooruLauncher()

type DeepdanbooruRequest struct {
	FilePath     string
	CompleteChan chan DeepdanbooruResult
}
type DeepdanbooruResult struct {
	Result []deepdanbooru.Predictions
	Error  error
}

func NewDeepdanbooruRequest(filePath string) *DeepdanbooruRequest {
	return &DeepdanbooruRequest{
		FilePath:     filePath,
		CompleteChan: make(chan DeepdanbooruResult),
	}
}
func (request *DeepdanbooruRequest) Wait() ([]deepdanbooru.Predictions, error) {
	result := <-request.CompleteChan
	return result.Result, result.Error
}

type DeepdanbooruLauncher struct {
	In     chan *DeepdanbooruRequest
	Logger *youlog.Scope
}

func NewDeepdanbooruLauncher() *DeepdanbooruLauncher {
	return &DeepdanbooruLauncher{
		In: make(chan *DeepdanbooruRequest),
	}
}

func (launcher *DeepdanbooruLauncher) Launch(filePath string) *DeepdanbooruRequest {
	request := NewDeepdanbooruRequest(filePath)
	launcher.In <- request
	return request
}

func (launcher *DeepdanbooruLauncher) Start() {
	go func() {
		for {
			select {
			case request := <-launcher.In:
				result, err := launcher.process(request)
				request.CompleteChan <- DeepdanbooruResult{
					Result: result,
					Error:  err,
				}
			}
		}
	}()
}

func (launcher *DeepdanbooruLauncher) process(request *DeepdanbooruRequest) ([]deepdanbooru.Predictions, error) {
	file, err := os.Open(request.FilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	predict, err := DefaultDeepDanbooruPlugin.Client.Tagging(file)
	if err != nil {
		return nil, err
	}
	return predict, nil
}
