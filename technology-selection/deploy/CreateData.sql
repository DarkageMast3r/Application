create database if not exists technology_selection;
use technology_selection;

drop table if exists Category;
create table Category (
	Id int primary key auto_increment,
	`Name` varchar(256),
	`Description` varchar(256),
	GeneratedOn DateTime
);

drop table if exists Need;
create table Need (
	Id int primary key auto_increment,
	`Name` varchar(256),
	`Source` varchar(256)
);

insert into Category (`Name`, `Description`, GeneratedOn) values ('Strains', 'For all strainst', '2025-6-27');
insert into Category (`Name`, `Description`, GeneratedOn) values ('Injuries', 'For all injuries', '2025-6-27');
insert into Category (`Name`, `Description`, GeneratedOn) values ('Bones', 'For all issues regarding bones', '2025-6-27');

insert into Need (`Name`, `Source`) values ('Valdetectie', 'Gebroken been');
insert into Need (`Name`, `Source`) values ('Geheugensteun', 'Dementia');

