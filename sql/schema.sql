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
	Author int not null,
	Submissions varchar(2400) not null, -- A list of up to 100 base64 encoded UUIDs to submissions.
	Improvement varchar(255) not null,
	primary key (ID)
);

create table submissions (
	ID int auto_increment not null,
	UUID char(24) not null, -- A base64 encoded UUID.
	Participation float not null,
	Collaboration float not null,
	Contribution float not null,
	Attitude float not null,
	Goals float not null,
	Comment varchar(255),
	primary key(ID)
);
