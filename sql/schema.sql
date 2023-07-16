drop table if exists teams;
drop table if exists members;
drop table if exists membersubmissions;
drop table if exists submissions;

create table teams (
	ID int auto_increment not null,
	Name varchar(255) not null,
	Members varchar(255) not null,
	primary key (ID)
);

create table users (
	ID int auto_increment not null,
	Name varchar(255) not null,
	TeamID int not null,
	primary key (ID)
);

create table submissions (
	ID int auto_increment not null,
	UUID binary(16) not null, -- A uuid_to_bin encoded uuid.
	Author int not null,
	Entries varchar(255) not null, -- A list of submission IDs.
	Improvement varchar(255) not null,
	primary key (ID)
);

create table memberentries (
	ID int auto_increment not null,
	Member int not null,
	Participation float not null,
	Collaboration float not null,
	Contribution float not null,
	Attitude float not null,
	Goals float not null,
	Comment varchar(255),
	primary key(ID)
);
