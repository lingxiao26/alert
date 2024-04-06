package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const TIME_LAYOUT = "2006-01-02 15:04:05"

func getLocalTime(startTime time.Time) string {
	sh, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Errorf("load location Asia/Shanghai: %v", err)
		return ""
	}
	return startTime.In(sh).Format(TIME_LAYOUT)
}

type Event struct {
	Receiver string `json:"receiver"`
	Status   string `json:"status"`

	Alerts []struct {
		Status string `json:"status"`
		Labels struct {
			Alertname string `json:"alertname"`
			Instance  string `json:"instance"`
		} `json:"labels"`

		Annotations struct {
			Summary string `json:"summary"`
			At      string `json:"at"`
			Wecom   string `json:"wecom"`
		} `json:"annotations"`

		StartsAt     time.Time `json:"startsAt"`
		EndsAt       time.Time `json:"endsAt"`
		GeneratorURL string    `json:"generatorURL"`
		Fingerprint  string    `json:"fingerprint"`
		SilenceURL   string    `json:"silenceURL"`
		DashboardURL string    `json:"dashboardURL"`
		PanelURL     string    `json:"panelURL"`
		Values       any       `json:"values"`
		ValueString  string    `json:"valueString"`
	} `json:"alerts"`

	GroupLabels struct {
	} `json:"groupLabels"`

	CommonLabels struct {
		Alertname string `json:"alertname"`
		Instance  string `json:"instance"`
	} `json:"commonLabels"`

	CommonAnnotations struct {
		Summary string `json:"summary"`
	} `json:"commonAnnotations"`

	ExternalURL     string `json:"externalURL"`
	Version         string `json:"version"`
	GroupKey        string `json:"groupKey"`
	TruncatedAlerts int    `json:"truncatedAlerts"`
	OrgID           int    `json:"orgId"`
	Title           string `json:"title"`
	State           string `json:"state"`
	Message         string `json:"message"`
}

func index(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorf("read request body: %v", err)
		return
	}

	var e Event
	if err := json.Unmarshal(data, &e); err != nil {
		log.Errorf("unmarshal event: %v", err)
		return
	}

	for _, alert := range e.Alerts {
		var sb strings.Builder
		// sb.WriteString("<font color=\"warning\">橙红色</font>" + alert.Labels.Alertname + "\n")
		sb.WriteString(fmt.Sprintf("<font color=\"warning\">%s</font>\n", alert.Labels.Alertname))
		sb.WriteString(">告警内容: " + alert.Annotations.Summary + "\n")
		sb.WriteString(">开始时间: " + getLocalTime(alert.StartsAt) + "\n")
		sb.WriteString(fmt.Sprintf("<@%s>", alert.Annotations.At))
		msg := &Message{
			Msgtype: "markdown",
			Markdown: struct {
				Content string "json:\"content\""
			}{
				Content: sb.String(),
			},
			Webhook: alert.Annotations.Wecom,
		}
		msg.sendAlert()
	}
}

type Message struct {
	Msgtype  string `json:"msgtype"`
	Markdown struct {
		Content string `json:"content"`
	} `json:"markdown"`
	Webhook string `json:"webhook"`
}

func (m *Message) sendAlert() {
	// 序列化请求体
	data, err := json.Marshal(m)
	if err != nil {
		log.Errorf("marshal alert message: %v", err)
		return
	}

	// 发送请求
	response, err := http.Post(m.Webhook, "application/json", bytes.NewReader(data))
	if err != nil {
		log.Errorf("send alert to %s: %v", m.Webhook, err)
		return
	}
	defer response.Body.Close()

	// 处理响应消息
	var respData = make([]byte, 1024)
	_, err = response.Body.Read(respData)
	if err != nil {
		log.Errorf("response body read: %v", err)
	}

	log.Infof("qiwei response: %s", respData)
}
