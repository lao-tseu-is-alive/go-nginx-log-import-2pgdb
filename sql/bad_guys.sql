SELECT COUNT(*),
       status,
       remote_addr,
       min(time_local),
       max(time_local),
       min(url),
       max(url),
       min(method),
       max(method)
FROM nginx_log
WHERE status > 304
  AND path_level < 7
  AND url not in ('/robots.txt', '/sitemap.xml', '/wp-login.php')
  AND time_local > '2023-01-01'
  AND remote_addr not in ('193.200.220.7')
GROUP BY status, remote_addr
HAVING COUNT(*) > 3
ORDER BY 1 DESC;

SELECT COUNT(*), 'ufw insert 1 reject from ' || nginx_log.remote_addr || ' to any port 443 '
FROM nginx_log
WHERE status > 304
  AND path_level < 7
  AND url not in ('/robots.txt', '/sitemap.xml', '/wp-login.php')
  AND time_local > '2023-01-01'
  AND remote_addr not in ('193.200.220.7')
GROUP BY status, remote_addr
HAVING COUNT(*) > 3
ORDER BY 1 DESC;