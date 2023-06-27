# go-nginx-log-import-2pgdb
Simple Go utility code to import Nginx log to Postgres database table


### how to use :

+ create the nginx log table if it doesn't exist by running :


    psql -f create_table_nginx_log.sql go_nginx_log



+ generate the csv file from your nginx access log


    go run cmd/convertNginxLogToCsv/convertNginxLogToCsv.go nginx_access.log 2> import_error.log > /tmp/nginxLogData.csv

+ import the csv in your database table


    COPY public.nginx_log (remote_addr, remote_user, time_local, method, url, protocol, status, body_bytes_sent, http_referer, http_user_agent)
    FROM '/tmp/nginxLogData.csv';


### performance :
I tested this code with a 15GB nginx log file containing  57'198'408 lines and the csv file was ready in ~18 min. Here are the time values :

    real	9m30.115s
    user	7m53.318s
    sys	1m55.410s

the end of my log tells me that i "imported 57'197'005 lines" and  

    tail import_error.log
    ...
    go-nginx-log-import-2pgdb]2023/06/23 17:27:40 convertNginxLogToCsv.go:148: # INFO: 'imported 57197005 lines (36 rejected lines with strange request)   log file : data/tilesmn95.log'

i can verify that with a wc :

    wc -l /tmp/nginxLogData.csv 
    57197005 /tmp/nginxLogData.csv

the final part is to import the csv file in the postgres database :

    time psql -c "COPY public.nginx_log (remote_addr, remote_user, time_local, method, url, protocol, status, body_bytes_sent, http_referer, http_user_agent)                                               
    FROM '/tmp/nginxLogData.csv';"   go_nginx_log
    COPY 57197005
    real	2m26,291s
    user	0m0,041s
    sys	0m0,013s

conclusion : 57 millions rows inserted from the nginx log file in less then 15 min on my Linux computer.   


### Nginx "combined" Log format

    remote_addr - remote_user [time_local] "request" status body_bytes_sent "http_referer" "http_user_agent"

|Field Name |	Description
| ------------- | ------------- |
|remote_addr	|Client IP address
|remote_user	|Client name
|time_local	|Local server time
method	|HTTP request method
|url	|URL
|protocol|	Protocol type
|status	|HTTP request status code
|body_bytes_sent|	Number of bytes sent to client
|http_referer	|Access source page URL
|http_user_agent	|Client browser information

### more info:
+ [Configuring Nginx Logging](https://docs.nginx.com/nginx/admin-guide/monitoring/logging/)
+ [Nginx access_log](https://nginx.org/en/docs/http/ngx_http_log_module.html?&_ga=2.19815224.1030081961.1687533866-1467598998.1687533866#access_log)
