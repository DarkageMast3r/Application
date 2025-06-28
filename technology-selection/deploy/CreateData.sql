create table Category (
	Id int primary key identity(1,1),
	Name varchar(256),
	Description varchar(256),
	GeneratedOn DateTime
)

insert into Category ([Name], [Description], GeneratedOn) values ('Strains', 'For all strainst', '6-27-2025')
insert into Category ([Name], [Description], GeneratedOn) values ('Injuries', 'For all injuries', '6-27-2025')
insert into Category ([Name], [Description], GeneratedOn) values ('Bones', 'For all issues regarding bones', '6-27-2025')
