create table user
(
    user_id          varchar(64)  not null
        primary key,
    provided_user_id varchar(64)  not null,
    login            varchar(64)  not null,
    name             varchar(128) not null,
    constraint uq_user_provided_id
        unique (provided_user_id)
)
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_unicode_520_ci;

create table media
(
    media_id     varchar(64)            not null
        primary key,
    user_id      varchar(64)            not null,
    description  varchar(255)           not null,
    content_type varchar(64) default '' not null,
    size         bigint      default 0  not null,
    width        bigint      default 0  not null,
    height       bigint      default 0  not null,
    constraint fk_media_user
        foreign key (user_id) references user (user_id)
);

create table playlist
(
    playlist_id varchar(64)  not null
        primary key,
    user_id     varchar(64)  not null,
    description varchar(255) not null,
    constraint fk_playlist_user
        foreign key (user_id) references user (user_id)
)
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_unicode_520_ci;

create table media_playlist
(
    media_playlist_id varchar(64) not null,
     playlist_id       varchar(64) not null,
    media_id          varchar(64) not null,
    seq_no            bigint      not null,
    constraint media_playlist_uq
        unique (playlist_id, media_id),
    constraint seq_no_uq
        unique (media_playlist_id, seq_no),
    constraint fk_md_pl_media
        foreign key (media_id) references media (media_id),
    constraint fk_md_pl_playlist
        foreign key (playlist_id) references playlist (playlist_id)
)
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_unicode_520_ci;

create index ix_fk_playlist_id
    on media_playlist (playlist_id);

create table session
(
    session_id  varchar(64)            not null
        primary key,
    user_id     varchar(64)            not null,
    token       varchar(128)           not null,
    expired_ind varchar(1) default 'N' not null,
    expiry_dt   datetime               not null,
    constraint fk_session_user
        foreign key (user_id) references user (user_id)
)
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_unicode_520_ci;

create table seqvalues
(
    id varchar(64) primary key,
    seqval bigint not null default 1
)
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_unicode_520_ci;

insert into seqvalues (id, seqval) value ('1', 1);
