CREATE DATABASE IF NOT EXISTS cflion DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;

create table app (
  id bigint(20) not null auto_increment,
  name varchar(45) not null comment 'app name',
  outdated tinyint(2) default 1 comment 'whether it is outdated, 1=yes, 0=no',
  ctime datetime DEFAULT NULL,
  utime timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  primary key (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

create table config_file (
  id bigint(20) not null auto_increment,
  name varchar(45) not null comment 'file name',
  app_id bigint(20) not null comment 'app id',
  ctime datetime DEFAULT NULL,
  utime timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  primary key (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
alter table config_file add index appId_INDEX (app_id);

create table association (
  id bigint(20) not null auto_increment,
  app_id bigint(20) not null,
  file_id bigint(20) not null,
  ctime datetime DEFAULT NULL,
  utime timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  primary key (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
alter table association add index appId_INDEX (app_id);

create table config_item (
  id bigint(20) not null auto_increment,
  file_id bigint(20) not null,
  name varchar(256) not null,
  value text not null,
  comment varchar(256) default null,
  ctime datetime DEFAULT NULL,
  utime timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  primary key (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
alter table config_item add index fileId_INDEX (file_id);
