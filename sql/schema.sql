create table teams (
  ID int auto_increment not null,
	TeamName varchar(255) not null,
  Members varchar(255) not null,
  primary key (ID)
);

create table members (
	ID int auto_increment not null,
	MemberName varchar(255) not null,
	primary key (ID)
);

create table membersubmissions (
	ID int auto_increment not null,
	UUID binary(16) not null, -- A uuid_to_bin encoded uuid.
	Author int not null,
	Submissions varchar(255) not null, -- A list of submission IDs.
	Improvement varchar(255) not null,
	primary key (ID)
);

create table submissions (
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
