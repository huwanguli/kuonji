package dto

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PageData struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

type SeriesInfo struct {
	Name       string `json:"name"`
	Count      int64  `json:"count"`
	LatestSlug string `json:"latest_slug"`
}

type SeriesLink struct {
	Slug  string `json:"slug"`
	Title string `json:"title"`
}

type ArticleDetail struct {
	Article      interface{} `json:"article"`
	PrevInSeries *SeriesLink `json:"prev_in_series,omitempty"`
	NextInSeries *SeriesLink `json:"next_in_series,omitempty"`
}

func Success(data interface{}) *Response {
	return &Response{Code: 200, Message: "success", Data: data}
}

func SuccessMessage(message string) *Response {
	return &Response{Code: 200, Message: message}
}

func Error(code int, message string) *Response {
	return &Response{Code: code, Message: message}
}

func BadRequest(message string) *Response {
	return &Response{Code: 400, Message: message}
}

func Unauthorized(message string) *Response {
	return &Response{Code: 401, Message: message}
}

func NotFound(message string) *Response {
	return &Response{Code: 404, Message: message}
}

func InternalError(message string) *Response {
	return &Response{Code: 500, Message: message}
}

func PageResult(list interface{}, total int64, page int, pageSize int) *Response {
	return Success(&PageData{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}
