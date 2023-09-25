# consul-ext
## use
1. You can use consul-ext backup your consul agent services and kvs to mysql. You can also quickly restore Consul kvs and services from MySQL data
2. GitOpt. you can commit msg to a bound repo, it will update kv to consul.In this way,you can log who and when update kv.
## begin
1. create config in /etc
```yaml
database: # mysql
    host: your mysql host IP
    port: your mysql port
    username: your mysql username
    password: your mysql password
    dbname: your mysql dbname
# you can use gitea or gitlab, only one of them.
# I recommend you use gitlab, because it is more stable.
# gitea is free,but its restful api is not abundant.It only can push file one by one.So it's slower.
# however, use all kv sync api is Low frequency,you can even use it only once,it's OK.
gitlab:
    host: your repo host IP
    port: you repo port
    # if use url ,host and port will be ignored
    url: your gitlab url # eg: http://www.mygitlab.net
    repo: your_group/your_repo
    branch: master
    token: your gitlab token
gitea:
    host: your repo host IP
    port: you repo port
    # if use url ,host and port will be ignored
    url: your gitlab url # eg: http://www.mygitea.net
    owner: your_repo_user
    repo: your_repo
    branch: main
    token: your gitea token
consul: # the consul agents you want backup
    - host: your consul host IP
      port: you repo port
      svc_name: # if not set, it will be all service backup
        - node_exporter
        - wmi_exporter
        - prometheus
backup_time: 7 # change this to control delete time
kv_prefix: # consul kv pefix,if not set,consul-ext get all kv
```
2. create table in you mysql
```sql
create table consul_svcs
(
    svc_id           varchar(256) default ''                not null,
    svc_name         varchar(64)  default ''                not null,
    consul_address   varchar(64)  default ''                not null comment 'consul agent address',
    svc_catalog_json json                                   not null comment 'svc details',
    backup_time      datetime     default CURRENT_TIMESTAMP not null comment 'backup time',
    primary key (svc_id, backup_time)
);
```
3. services date in mysql will delete over 7 days,you can change this in config, eg
4. api list
```bash
PUT # restore services from mysql to consul
/api/v1/consul-ext/svc/restore
GET # sync all consul kv to git repo
/api/v1/consul-ext/kv
PUT # sync git repo file to consul
/api/v1/consul-ext/path/file
POST # webhook sync git repo change to consul
/api/v1/consul-ext/:repoType/webhook
```
5. deploy this project,use docker or k8s,then start sync
