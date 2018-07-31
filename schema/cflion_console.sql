CREATE DATABASE IF NOT EXISTS cflion_console DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;

use cflion_console;

create table app (
  id bigint(20) not null auto_increment,
  name varchar(45) not null comment 'app name',
  env varchar(45) not null comment 'app env',
  outdated tinyint(2) default 1 comment 'whether it is outdated, 1=yes, 0=no',
  ctime datetime DEFAULT NULL,
  utime timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  primary key (id)
)  ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

# create table config (
#   id bigint(20) not null auto_increment,
#   name varchar(256) not null,
#   value text not null,
#   comment varchar(256) default null,
#   ctime datetime DEFAULT NULL,
#   utime timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
#   primary key (id)
# ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
# insert into config (name, value, comment, ctime, utime) values ('dev.manager.url', 'http://127.0.0.1:8080', '', now(), now());
