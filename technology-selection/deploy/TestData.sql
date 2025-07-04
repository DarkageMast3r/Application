create database if not exists technology_selection_test;
use technology_selection_test;

drop table if exists TechNeed;
drop table if exists Need;
drop table if exists TechChoice;
drop table if exists Tech;
drop table if exists Category;
drop table if exists `Case`;

create table `Case` (
	Id int auto_increment,
	ClientId varchar(256),
    `Name` varchar(256),
	`Description` varchar(4096),
    primary key (Id)
);

create table Category (
	Id int auto_increment,
	`Name` varchar(256),
	`Description` varchar(256),
    primary key (Id)
);

create table Tech (
	Id int primary key auto_increment,
    `Name` varchar(256),
    CategoryId int,
    Cost decimal(16, 6),
	constraint fk_Category foreign key (CategoryId) references Category(Id)
);

create table TechChoice (
	Id int primary key auto_increment,
    TechId int,
	CaseId int,
    `Status` int,
    `Reasoning` varchar(256),
    foreign key (TechId) references Tech(Id)
);

create table Need (
	Id int primary key auto_increment,
    `Description` varchar(256)
);

create table TechNeed (
	TechId int,
	NeedId int,
    foreign key (TechId) references Tech(Id),
    foreign key (NeedId) references Need(Id)
);

insert into Category (`Name`, `Description`) values ('Strains', 'For all strainst');
insert into Category (`Name`, `Description`) values ('Injuries', 'For all injuries');
insert into Category (`Name`, `Description`) values ('Bones', 'For all issues regarding bones');

insert into Tech (`Name`, `CategoryId`, `Cost`) values ('Rolstoel', 3, 100);
insert into Tech (`Name`, `CategoryId`, `Cost`) values ('Rijstoel', 3, 2000);
insert into Tech (`Name`, `CategoryId`, `Cost`) values ('Pleister', 2, 2500.13);
insert into Tech (`Name`, `CategoryId`, `Cost`) values ('Aspirine', 2, 503);

insert into Need (`Description`) values ('Wandelondersteuning');
insert into Need (`Description`) values ('Comfort');
insert into Need (`Description`) values ('Geheugensteun');
insert into Need (`Description`) values ('Pijnstiller');

insert into TechNeed (`TechId`, `NeedId`) values (1, 1);
insert into TechNeed (`TechId`, `NeedId`) values (2, 1);
insert into TechNeed (`TechId`, `NeedId`) values (1, 2);

