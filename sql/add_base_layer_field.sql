SELECT url ,base_layer,
       array_length(regexp_split_to_array(url, E'\\/'), 1),
       regexp_split_to_array(url, E'\\/'),
       split_part(url, '/',4)
FROM nginx_log
WHERE id<10;

ALTER TABLE nginx_log ADD COLUMN path_level integer;
ALTER TABLE nginx_log ADD COLUMN base_layer text;

UPDATE nginx_log SET path_level = array_length(regexp_split_to_array(url, E'\\/'), 1);
UPDATE nginx_log SET base_layer = split_part(url, '/', 4)
WHERE status <300 AND path_level >6;

SELECT COUNT(*) FROM nginx_log WHERE base_layer is not  null;

EXPLAIN SELECT COUNT(*),status,path_level, remote_addr, min(time_local),max(time_local)
FROM nginx_log
WHERE status > 304 AND path_level < 7
GROUP BY status,path_level, remote_addr
ORDER BY 1 DESC