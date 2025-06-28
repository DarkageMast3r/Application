create database technology_selection;
use technology_selection;

create table Category (
	Id int primary key auto_increment,
	`Name` varchar(256),
	`Description` varchar(256),
	GeneratedOn DateTime
);

insert into Category (`Name`, `Description`, GeneratedOn) values ('Strains', 'For all strainst', '2025-6-27');
insert into Category (`Name`, `Description`, GeneratedOn) values ('Injuries', 'For all injuries', '2025-6-27');
insert into Category (`Name`, `Description`, GeneratedOn) values ('Bones', 'For all issues regarding bones', '2025-6-27');
