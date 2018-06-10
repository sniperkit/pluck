package metric

import (
	"fmt"
	"net/url"
	"time"
)

type Response struct {
	ParentURL       url.URL   `json:'parent_url' yaml:'parent_url' toml:'-' xml:'parentURL' ini:'-' csv:'parentURL'`
	URL             url.URL   `json:'url' yaml:'url' toml:'-' xml:'url' ini:'-' csv:'url'`
	StatusCode      int       `json:'status_code' yaml:'status_code' toml:'-' xml:'statusCode' ini:'-' csv:'statusCode'`
	NumberOfWorkers int       `json:'workers' yaml:'workers' toml:'-' xml:'workers' csv:'workers' ini:'-'`
	WorkerID        int       `json:'worker_id' yaml:'-' toml:'-' xml:'workerID' csv:'workerID' ini:'-'`
	ResponseSize    int       `json:'response_size' yaml:'response_size' toml:'response_size' xml:'responseSize' csv:'responseSize' ini:'-'`
	ContentType     string    `json:'content_type' yaml:'content_type' toml:'content_type' xml:'contentType' ini:'contentType' ini:'-'`
	StartTime       time.Time `json:'start_time' yaml:'start_time' toml:'' xml:'startTime' csv:'startTime' ini:'-'`
	EndTime         time.Time `json:'end_time' yaml:'end_time' toml:'' xml:'endTime' csv:'endTime' ini:'-'`
	Err             error     `json:'error' yaml:'error' toml:'' xml:'error' csv:'MessageError' ini:'-'`
}

func (r Response) String() string {
	return fmt.Sprintf("#%03d: %03d %9s %15s %20s",
		r.WorkerID,
		r.StatusCode,
		fmt.Sprintf("%d", r.ResponseSize),
		fmt.Sprintf("%f ms", r.GetResponseTime().Seconds()*1000),
		r.URL.String(),
	)
}

func (r Response) GetError() error {
	return r.Err
}

func (r Response) GetParentURL() url.URL {
	return r.ParentURL
}

func (r Response) GetURL() url.URL {
	return r.URL
}

func (r Response) GetSize() int {
	return r.ResponseSize
}

func (r Response) GetStatusCode() int {
	return r.StatusCode
}

func (r Response) GetStartTime() time.Time {
	return r.StartTime
}

func (r Response) GetEndTime() time.Time {
	return r.EndTime
}

func (r Response) GetResponseTime() time.Duration {
	return r.EndTime.Sub(r.StartTime)
}

func (r Response) GetContentType() string {
	return r.ContentType
}

func (r Response) GetWorkerID() int {
	return r.WorkerID
}

func (r Response) GetNumberOfWorkers() int {
	return r.NumberOfWorkers
}
