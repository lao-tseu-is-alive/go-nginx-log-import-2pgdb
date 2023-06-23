package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

const (
	VERSION        = "0.1.1"
	APP            = "go-nginx-log-import-2pgdb"
	defaultLogFile = "data/sample.log"
)

type Month int

const (
	Jan Month = iota + 1
	Feb
	Mar
	Apr
	May
	Jun
	Jul
	Aug
	Sep
	Oct
	Nov
	Dec
)

var (
	MonthMap = map[string]Month{
		"jan": Jan,
		"feb": Feb,
		"mar": Mar,
		"apr": Apr,
		"may": May,
		"jun": Jun,
		"jul": Jul,
		"aug": Aug,
		"sep": Sep,
		"oct": Oct,
		"nov": Nov,
		"dec": Dec,
	}
)

func ConvString2Month(strMonth string) (Month, bool) {
	month, found := MonthMap[strings.ToLower(strMonth)]
	return month, found
}

//goland:noinspection RegExpRedundantEscape
func main() {
	args := os.Args[1:]
	var logPath string
	l := log.New(os.Stderr, fmt.Sprintf("[%s]", APP), log.Ldate|log.Ltime|log.Lshortfile)
	//l := log.New(io.Discard, fmt.Sprintf("[%s]", APP), log.Ldate|log.Ltime|log.Lshortfile)
	//l.Printf("# INFO: 'Starting %s version:%s  num args:%d'\n", APP, VERSION, len(args))
	if len(args) == 1 {
		logPath = os.Args[1]
	} else {
		flag.StringVar(&logPath, "f", defaultLogFile, "Path to your log file")
		flag.Parse()
	}
	l.Printf("# INFO: 'about to open log file : %s'\n", logPath)
	file, err := os.Open(logPath)
	if err != nil {
		l.Fatalf("💥💥 ERROR: 'problem opening log at os.Open(*logPath:%s), got error: %v'\n", logPath, err)
	}
	defer file.Close()

	// NGINX “combined” log format: http://nginx.org/en/docs/http/ngx_http_log_module.html#log_format
	var myNginxRegex = regexp.MustCompile(`^(?P<remote_addr>[^ ]+)\s-\s(?P<remote_user>[^ ]+)\s\[(?P<time_local>[^\]]+)\]\s"(?P<request>[^"]*)"\s(?P<status>\d{1,3})\s(?P<body_bytes_send>\d+)\s"(?P<http_referer>[^"]*)"\s"(?P<http_user_agent>[^"]*)"`)
	var myDateTimeRegex = regexp.MustCompile("^(?P<day>\\d{1,2})\\/(?P<month>\\w{1,3})\\/(?P<year>\\d{2,4}):(?P<hour>\\d{1,2}):(?P<minute>\\d{1,2}):(?P<second>\\d{1,2})")

	l.Printf("# INFO: 'about to read log file : %s'\n", logPath)
	scanner := bufio.NewScanner(file)
	numLine, goodLines, badLines := 0, 0, 0
	for scanner.Scan() {
		// load a line of log
		line := scanner.Text()
		numLine++
		if numLine%5 == 0 {
			l.Printf("# DEBUG: 'handling line number: %d'\n", numLine)
		}
		match := myNginxRegex.FindStringSubmatch(line)
		nginxCombinedFields := make(map[string]string)
		for i, name := range myNginxRegex.SubexpNames() {
			if i != 0 && name != "" {
				nginxCombinedFields[name] = match[i]
			}
		}
		// verb, url, protocol := strings.Split(nginxCombinedFields["request"], " ")
		requestParts := strings.Split(nginxCombinedFields["request"], " ")
		// usually request should be [HTTP_VERB URL PROTOCOL] like in : "GET /index.html HTTP/1.1"
		if len(requestParts) > 2 {
			// let's keep only the http verb for this task
			if requestParts[0] == "GET" {
				matchDate := myDateTimeRegex.FindStringSubmatch(nginxCombinedFields["time_local"])
				nginxDateTimeFields := make(map[string]string)
				for j, name := range myDateTimeRegex.SubexpNames() {
					if j != 0 && name != "" {
						nginxDateTimeFields[name] = matchDate[j]
					}
				}
				monthInNumber, success := ConvString2Month(nginxDateTimeFields["month"])
				if !success {
					l.Printf("## Warning ConvString2Month does not know how to convert month for %s\n", nginxCombinedFields["time_local"])
				}
				var requestUrl string
				if strings.Contains(requestParts[1], "?") {
					urlParts := strings.Split(requestParts[1], "?")
					// let's keep only the url without the query for now
					requestUrl = urlParts[0]
				} else {
					requestUrl = requestParts[1]
				}
				goodLines++
				fmt.Printf("%s\t%s\t%s-%02d-%s %s:%s:%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
					nginxCombinedFields["remote_addr"],
					nginxCombinedFields["remote_user"],
					nginxDateTimeFields["year"],
					monthInNumber,
					nginxDateTimeFields["day"],
					nginxDateTimeFields["hour"],
					nginxDateTimeFields["minute"],
					nginxDateTimeFields["second"],
					requestParts[0],
					requestUrl,
					requestParts[2], //protocol
					nginxCombinedFields["status"],
					nginxCombinedFields["body_bytes_send"],
					nginxCombinedFields["http_referer"],
					nginxCombinedFields["http_user_agent"],
				)
			}
		}
		badLines++
		l.Printf("# 💥💥 WARNING: 'unusual http request found on line %d of : %s'\n", numLine, logPath)
		l.Printf("# 💥💥 DISCARDED_LINE [%d]:\t%s\n", numLine, line)
	}
	l.Printf("# INFO: 'imported %d lines (%d rejected lines with strange request)   log file : %s'\n", goodLines, badLines, logPath)
}
