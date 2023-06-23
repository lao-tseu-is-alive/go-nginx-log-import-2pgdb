create table if not exists public.nginx_log
(
    id              serial,
    remote_addr     text                   not null,
    remote_user     text default '-'::text not null,
    time_local      timestamp              not null,
    method          text                   not null,
    url             text,
    protocol        text,
    status          integer                not null,
    body_bytes_sent integer                not null,
    http_referer    text default '-'::text not null,
    http_user_agent text default '-'::text not null
);

comment on table public.nginx_log is 'combined nginx access log';

comment on column public.nginx_log.id is 'id of log entry';

comment on column public.nginx_log.remote_addr is 'Client IP address';

comment on column public.nginx_log.remote_user is 'Client user name';

comment on column public.nginx_log.time_local is 'Local server time of this request';

comment on column public.nginx_log.method is 'HTTP request method';

comment on column public.nginx_log.status is 'HTTP request status code';

comment on column public.nginx_log.body_bytes_sent is 'Number of bytes sent to client';

comment on column public.nginx_log.http_referer is 'Access source page URL';

comment on column public.nginx_log.http_user_agent is 'Client browser information';

alter table public.nginx_log
    owner to go_nginx_log;

