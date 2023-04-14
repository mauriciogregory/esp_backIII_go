create schema if not exists dental_clinic;
use dental_clinic;

create table patients (
    id int not null auto_increment,
    surname varchar(50) not null,
    name varchar(25) not null,
    identity_number varchar(10) not null unique,
    created_at datetime not null,
    primary key (id)
);

create table dentists (
    id int not null auto_increment,
    surname varchar(50) not null,
    name varchar(25) not null,
    license_number varchar(10) not null unique,
    primary key (id)
);

create table appointments (
    id int not null auto_increment,
    description varchar(250) not null,
    date_and_time datetime not null,
    dentist_license varchar(10) not null,
    patient_identity varchar(10) not null,
    primary key (id),
    constraint fk_dentist foreign key (dentist_license) references dentists(license_number),
    constraint fk_patient foreign key (patient_identity) references patients(identity_number)
);