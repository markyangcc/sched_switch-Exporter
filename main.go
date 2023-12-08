package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type SwitchesData struct {
	Pid                   int
	Process               string
	NrSwitches            int
	NrVoluntarySwitches   int
	NrInvoluntarySwitches int
}

var (
	Process = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "Process",
		Help: "",
	}, []string{"Pid"})
	NrSwitches = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "nr_switches",
		Help: "",
	}, []string{"Pid"})

	NrVoluntarySwitches = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "nr_voluntary_switches",
		Help: "",
	}, []string{"Pid"})

	NrInvoluntarySwitches = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "nr_involuntary_switches",
		Help: "",
	}, []string{"Pid"})
)

func queryMetrics() {

	pidList := []int{33868, 33869, 33870, 33871, 33872, 33873, 33874, 33875}

	// 使用range循环遍历PID列表
	for _, pid := range pidList {

		// 在这里执行针对每个PID的操作

		// nr_switches metrics
		// 使用字符串格式化将%d替换为PID值
		shellcmd := fmt.Sprintf("cat /proc/%d/sched | grep nr_switches | awk '{print $3}'", pid)
		cmd := exec.Command("bash", "-c", shellcmd)
		nr_switches, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Run nr_switches command failed: %s\n", err)
			return
		}
		nr_switchesInt, err := strconv.Atoi(strings.TrimSpace(string(nr_switches)))
		if err != nil {
			fmt.Printf("string to int failed")
		}

		// nr_voluntary_switches metrics
		shellcmd = fmt.Sprintf("cat /proc/%d/sched | grep nr_voluntary_switches | awk '{print $3}'", pid)
		cmd = exec.Command("bash", "-c", shellcmd)
		nr_voluntary_switches, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Run nr_voluntary_switches command failed: %s\n", err)
			return
		}
		nr_voluntary_switchesInt, err := strconv.Atoi(strings.TrimSpace(string(nr_voluntary_switches)))
		if err != nil {
			fmt.Printf("string to int failed")
		}

		// nr_involuntary_switches metrics
		shellcmd = fmt.Sprintf("cat /proc/%d/sched | grep nr_involuntary_switches | awk '{print $3}'", pid)
		cmd = exec.Command("bash", "-c", shellcmd)
		nr_involuntary_switches, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Run nr_involuntary_switches command failed: %s\n", err)
			return
		}
		nr_involuntary_switchesInt, err := strconv.Atoi(strings.TrimSpace(string(nr_involuntary_switches)))
		if err != nil {
			fmt.Printf("string to int failed")
		}

		NrSwitches.WithLabelValues(strconv.Itoa(pid)).Set(float64(nr_switchesInt))
		NrVoluntarySwitches.WithLabelValues(strconv.Itoa(pid)).Set(float64(nr_voluntary_switchesInt))
		NrInvoluntarySwitches.WithLabelValues(strconv.Itoa(pid)).Set(float64(nr_involuntary_switchesInt))
	}

}

func init() {

	prometheus.MustRegister(NrSwitches, NrVoluntarySwitches, NrInvoluntarySwitches)
}

func main() {
	// default listen on port 8005
	port := flag.Int("port", 8005, "the port number to listen on")
	flag.Parse()
	portListenOn := strconv.Itoa(*port)

	http.Handle("/metrics", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryMetrics()

		// 直接传递 promhttp.Handler() 作为 handler
		promhttp.Handler().ServeHTTP(w, r)
	}))
	log.Fatal(http.ListenAndServe(":"+portListenOn, nil))
}
