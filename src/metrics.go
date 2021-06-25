package main

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/golang-collections/go-datastructures/queue"
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
)

var parseMap = map[string]func(string, *queue.Queue) map[string]interface{}{
	"advmss":         metricOne,
	"app_limited":    metricNone,
	"ato":            metricOne,
	"busy":           metricOne,
	"bytes_acked":    metricOne,
	"bytes_received": metricOne,
	"cubic":          metricNone,
	"cwnd":           metricOne,
	"data_segs_in":   metricOne,
	"data_segs_out":  metricOne,
	"delivery_rate":  metricTwo,
	"lastack":        metricOne,
	"lastrcv":        metricOne,
	"lastsnd":        metricOne,
	"minrtt":         metricOne,
	"mss":            metricOne,
	"pacing_rate":    metricTwo,
	"pmtu":           metricOne,
	"rcvmss":         metricOne,
	"rcv_rtt":        metricOne,
	"rcv_space":      metricOne,
	"rcv_ssthresh":   metricOne,
	"reordering":     metricOne,
	"rto":            metricOne,
	"rtt":            metricRtt,
	"sack":           metricNone,
	"segs_in":        metricOne,
	"segs_out":       metricOne,
	"send":           metricTwo,
	"ssthresh":       metricOne,
	"ts":             metricNone,
	"unacked":        metricOne,
	"wscale":         metricWscale,
}

var typeMap = map[string]metric.SourceType{
	"advmss":         metric.GAUGE,
	"app_limited":    metric.ATTRIBUTE,
	"ato":            metric.GAUGE,
	"busy":           metric.GAUGE,
	"bytes_acked":    metric.GAUGE,
	"bytes_received": metric.GAUGE,
	"cubic":          metric.ATTRIBUTE,
	"cwnd":           metric.GAUGE,
	"data_segs_in":   metric.GAUGE,
	"data_segs_out":  metric.GAUGE,
	"delivery_rate":  metric.GAUGE,
	"lastack":        metric.GAUGE,
	"lastrcv":        metric.GAUGE,
	"lastsnd":        metric.GAUGE,
	"minrtt":         metric.GAUGE,
	"mss":            metric.GAUGE,
	"pacing_rate":    metric.GAUGE,
	"pmtu":           metric.GAUGE,
	"rcvmss":         metric.GAUGE,
	"rcv_rtt":        metric.GAUGE,
	"rcv_space":      metric.GAUGE,
	"rcv_ssthresh":   metric.GAUGE,
	"reordering":     metric.GAUGE,
	"rto":            metric.GAUGE,
	"rtt_average":    metric.GAUGE,
	"rtt_std_dev":    metric.GAUGE,
	"sack":           metric.ATTRIBUTE,
	"segs_in":        metric.GAUGE,
	"segs_out":       metric.GAUGE,
	"send":           metric.GAUGE,
	"ssthresh":       metric.GAUGE,
	"ts":             metric.ATTRIBUTE,
	"unacked":        metric.GAUGE,
	"rcv_wscale":     metric.GAUGE,
	"snd_wscale":     metric.GAUGE,
}

func getMetrics(entity *integration.Entity, args argumentList) error {
	out, err := exec.Command("ss", getCommandArgs(), Args.Filter).Output()
	if err != nil {
		log.Error("Error executing command: %+v", err)
	}

	lines := strings.Split(string(out), "\n")

	for i := 1; i < len(lines); i = i + 2 {
		if lines[i] == "" {
			continue
		}
		ms := entity.NewMetricSet("SocketStatisticsSample")
		for m, v := range getMetric(lines[i], lines[i+1]) {
			ms.SetMetric(m, v, getSourceType(m))
		}
	}
	return nil
}
func getCommandArgs() string {
	args := Args.SSArgs
	if Args.Resolve {
		args += "r"
	}
	return args
}

func getFilter(src string, dst string) string {
	filter := "( "

	for _, s := range strings.Fields(src) {
		filter += ("src " + s + " or ")
	}

	for _, d := range strings.Fields(dst) {
		filter += ("dst " + d + " or ")
	}

	filter += " )"

	if strings.HasSuffix(filter, " or  )") {
		filter = strings.Replace(filter, " or  )", " )", 1)
	}

	if filter == "(  )" {
		filter = ""
	}
	return filter
}

func getSourceType(k string) metric.SourceType {
	st, ok := typeMap[k]
	if !ok {
		log.Debug("SourceType not found: %s", k)
		st = metric.ATTRIBUTE
	}
	return st
}

func getMetric(header string, data string) map[string]interface{} {
	//log.Warn("header: %s data: %s", header, data)
	metrics := make(map[string]interface{})
	if header == "" || data == "" {
		return metrics
	}

	fields := strings.Fields(header)
	if len(fields) >= 5 {
		//log.Debug("fields: %+v", fields)
		//log.Debug("source: %+v destination: %+v", fields[3], fields[4])
		metrics["source"] = fields[3]
		metrics["destination"] = fields[4]
	} else {
		log.Warn("Host line from ss too few fields: %+v", header)
		return metrics
	}

	values := strings.Fields(data)
	q := queue.New(int64(len(values)))
	defer q.Dispose()
	for _, f := range values {
		q.Put(f)
	}

	for !q.Empty() {
		i, err := q.Get(1)
		if err != nil {
			log.Warn("Queue err: %s", err)
			continue
		}
		if len(i) != 1 {
			log.Warn("Array != 1: %+v", i)
			continue
		}

		s := (i[0]).(string)
		var found bool = false
		for k, v := range parseMap {
			if strings.HasPrefix(s, k) {
				append(metrics, v(s, q))
				found = true
				break
			}
		}
		if !found {
			log.Info("Unknown result value: %s", s)
		}
	}
	return metrics
}

func append(d map[string]interface{}, s map[string]interface{}) {
	for k, v := range s {
		d[k] = v
	}
}

func metricNone(s string, q *queue.Queue) map[string]interface{} {
	m := make(map[string]interface{})
	m[s] = "true"
	return m
}

func metricOne(s string, q *queue.Queue) map[string]interface{} {
	m := make(map[string]interface{})
	a := strings.Split(s, ":")
	if len(a) == 2 {
		m[a[0]] = stringToNumber(a[1])
	} else {
		log.Warn("metricOne: unparsable: %s", s)
	}
	return m
}

func metricTwo(s string, q *queue.Queue) map[string]interface{} {
	m := make(map[string]interface{})
	v, err := q.Get(1)
	if err != nil || len(v) != 1 {
		log.Warn("metricTwo: err: %s len(v): %d", err, len(v))
		return m
	}
	switch x := v[0].(string); {
	case strings.HasSuffix(x, "Kbps"):
		f, err := strconv.ParseFloat(strings.TrimSuffix(x, "Kbps"), 64)
		if err != nil {
			log.Warn("Invalid format: %s: ", x)
			break
		}
		f = f * 1000
		m[s] = f
		return m
	case strings.HasSuffix(v[0].(string), "Mbps"):
		f, err := strconv.ParseFloat(strings.TrimSuffix(x, "Mbps"), 64)
		if err != nil {
			log.Warn("Invalid format: %s: ", x)
			break
		}
		f = f * 1000000
		m[s] = f
		return m

	case strings.HasSuffix(v[0].(string), "Gbps"):
		f, err := strconv.ParseFloat(strings.TrimSuffix(x, "Gbps"), 64)
		if err != nil {
			log.Warn("Invalid format: %s: ", x)
			break
		}
		f = f * 1000000000
		m[s] = f
		return m

	}
	m[s] = stringToNumber(v[0].(string))
	return m
}

func metricRtt(s string, q *queue.Queue) map[string]interface{} {
	m := make(map[string]interface{})
	a := strings.Split(s, ":")
	if len(a) != 2 {
		log.Warn("Invalid rtt format: %s", s)
	} else {
		b := strings.Split(a[1], "/")
		if len(b) != 2 {
			log.Warn("Invalid rtt format: %s", s)
		} else {
			m["rtt_average"] = stringToNumber(b[0])
			m["rtt_std_dev"] = stringToNumber(b[1])
		}
	}
	return m
}

func metricWscale(s string, q *queue.Queue) map[string]interface{} {
	m := make(map[string]interface{})
	a := strings.Split(s, ":")
	if len(a) != 2 {
		log.Warn("Invalid wscale format: %s", s)
	} else {
		b := strings.Split(a[1], ",")
		if len(b) != 2 {
			log.Warn("Invalid wscale format: %s", s)
		} else {
			m["snd_wscale"] = stringToNumber(b[0])
			m["rcv_wscale"] = stringToNumber(b[1])
		}
	}
	return m
}

func stringToNumber(s string) interface{} {
	s = strings.TrimSpace(s)
	//log.Debug("Converting: %s", s)

	i, err := strconv.ParseInt(s, 0, 32)
	if err == nil {
		//log.Debug("Return int: %d", i)
		return i
	}

	f, err := strconv.ParseFloat(s, 64)
	if err == nil {
		//log.Debug("Return float: %f", f)
		return f
	}

	//log.Debug("Return string: %s", s)
	return s
}
