create table `users` (
  `id` VARCHAR(255) not null,
  `email` varchar(255) not null,
  `password` varchar(255) not null,
  `created_at` datetime not null default CURRENT_TIMESTAMP,
  primary key (`id`)
)

create table `clients` (
  `id` VARCHAR(255) not null,
  `name` varchar(255) not null,
  `secret` varchar(255) not null,
  `created_at` datetime not null default CURRENT_TIMESTAMP,
  `redirect_uri` varchar(255) not null,
  primary key (`id`)
)