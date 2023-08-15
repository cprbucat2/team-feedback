drop table if exists teams;
drop table if exists users;
drop table if exists submissions;
drop table if exists memberentries;

create table teams (
	id bigint auto_increment not null primary key,
	name varchar(255) not null
);

insert into teams (name) values ("none");

create table users (
	id bigint auto_increment not null primary key,
	name varchar(255) not null,
	team_id bigint not null,
	index (team_id),
	foreign key (team_id) references teams (id)
);

insert into users (name, team_id) values ("nobody", 1);

create table submissions (
	id bigint auto_increment not null primary key,
	author bigint not null,
	improvement varchar(255) not null,
	index (author),
	foreign key (author) references users (id)
);

create table entries (
	id bigint auto_increment not null primary key,
	submission_id bigint not null,
	member bigint not null,
	Participation float not null,
	Collaboration float not null,
	Contribution float not null,
	Attitude float not null,
	Goals float not null,
	Comment varchar(255),
	index (submission_id),
	foreign key (submission_id) references submissions (id),
	index (member),
	foreign key (member) references users (id)
);

insert into teams (name) values ("Team Feedback");
insert into users (name, team_id) values ("Aiden Woodruff", 2);
insert into users (name, team_id) values ("Aidan Hoover", 2);
insert into teams (name) values ("The Usual Suspects");
insert into users (name, team_id) values ("Keaton", 3);
insert into users (name, team_id) values ("MacManus", 3);
insert into users (name, team_id) values ("Fenster", 3);
insert into users (name, team_id) values ("Hockney", 3);
insert into users (name, team_id) values ("Verbal", 3);
