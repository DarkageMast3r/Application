create database if not exists signaling;
use signaling;

drop table if exists signals;
drop table if exists assessments;
drop table if exists classifications;

create table signals (
	Id int auto_increment,
    client_id varchar(64),
    `type` varchar(256),
	waarde decimal,
    tijdstip datetime,
    bron varchar(256),
    primary key (Id)
);

create table assessments (
	Id int auto_increment,
    client_id varchar(64),
    `conclusie` varchar(4096),
	urgentie int,
    gevalideerd_door int,
    tijdstip datetime,
    primary key (Id)
);

create table classifications (
	Id int auto_increment,
    client_id varchar(64),
    categorie varchar(256),
    ernst varchar(256),
    motivatie varchar(256),
    primary key (Id)
);
