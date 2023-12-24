package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"com.levi/project-common/base"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w CustomResponseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func ActionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("action")

		// 接口耗时统计起始
		start := time.Now()
		blw := &CustomResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()
		var (
			requestId  string
			userId     string
			input      map[string]interface{}
			output     map[string]interface{}
			outputList base.Result
		)
		// 用户Id
		if v, exists := c.Keys["userId"]; exists {
			userId = fmt.Sprintf("%v", v)
		}
		// 全局请求ID
		if v, exists := c.Keys["RequestID"]; exists {
			requestId = fmt.Sprintf("%v", v)
		}
		// 输入参数
		if inputList, err := getRequestInputs(c); err == nil {
			input = inputList
		}
		// 输出参数
		json.Unmarshal(blw.body.Bytes(), &outputList)
		if (outputList.Data == nil) {
			outputList.Data = make(map[string]interface{})
		}
		output = outputList.Data.(map[string]interface{})
		
		zap.L().Info(
			"StdLog",
			zap.String("requestId", requestId),
			zap.String("url", c.Request.URL.String()),
			zap.Int("errno", int(outputList.Code)),
			zap.String("status", strconv.Itoa(c.Writer.Status())),
			zap.String("userId", userId),
			zap.Any("input", input),
			zap.Any("output", output),
			zap.String("cost", time.Since(start).String()),
		)
	}
}

func getRequestInputs(c *gin.Context) (map[string]interface{}, error) {
	const defaultMemory = 32 << 20
	contentType := c.ContentType()
	var (
		dataMap  = make(map[string]interface{})
		queryMap = make(map[string]interface{})
		postMap  = make(map[string]interface{})
	)
	// @see gin@v1.7.7/binding/query.go ==> func (queryBinding) Bind(req *http.Request, obj interface{})
	for k := range c.Request.URL.Query() {
		queryMap[k] = c.Query(k)
	}

	if contentType == "application/json" {
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		// @see gin@v1.7.7/binding/json.go ==> func (jsonBinding) Bind(req *http.Request, obj interface{})
		if c.Request != nil && c.Request.Body != nil {
			if err := json.NewDecoder(c.Request.Body).Decode(&postMap); err != nil {
				return nil, err
			}
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	} else if contentType == "multipart/form-data" {
		// @see gin@v1.7.7/binding/form.go ==> func (formMultipartBinding) Bind(req *http.Request, obj interface{})
		if err := c.Request.ParseMultipartForm(defaultMemory); err != nil {
			return nil, err
		}
		for k, v := range c.Request.PostForm {
			if len(v) > 1 {
				postMap[k] = v
			} else if len(v) == 1 {
				postMap[k] = v[0]
			}
		}
	} else {
		// ParseForm 解析 URL 中的查询字符串，并将解析结果更新到 r.Form 字段
		// 对于 POST 或 PUT 请求，ParseForm 还会将 body 当作表单解析，
		// 并将结果既更新到 r.PostForm 也更新到 r.Form。解析结果中，
		// POST 或 PUT 请求主体要优先于 URL 查询字符串（同名变量，主体的值在查询字符串的值前面）
		// @see gin@v1.7.7/binding/form.go ==> func (formBinding) Bind(req *http.Request, obj interface{})
		if err := c.Request.ParseForm(); err != nil {
			return nil, err
		}
		if err := c.Request.ParseMultipartForm(defaultMemory); err != nil {
			if err != http.ErrNotMultipart {
				return nil, err
			}
		}
		for k, v := range c.Request.PostForm {
			if len(v) > 1 {
				postMap[k] = v
			} else if len(v) == 1 {
				postMap[k] = v[0]
			}
		}
	}
	var mu sync.RWMutex
	for k, v := range queryMap {
		mu.Lock()
		dataMap[k] = v
		mu.Unlock()
	}
	for k, v := range postMap {
		mu.Lock()
		dataMap[k] = v
		mu.Unlock()
	}
	return dataMap, nil
}
