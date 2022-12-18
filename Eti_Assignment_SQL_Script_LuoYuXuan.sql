CREATE database trip_db;

USE trip_db;

SET SQL_MODE="NO_AUTO_VALUE_ON_ZERO";

CREATE TABLE User 
(
  UserID INT(8) NOT NULL AUTO_INCREMENT,
  UserType VARCHAR(20) NOT NULL,
  Email VARCHAR(50) NOT NULL,  
  Password VARCHAR(255) NOT NULL,
  PRIMARY KEY (UserID)
);


CREATE TABLE Passenger 
(
  PassengerID INT(8) NOT NULL,
  FirstName VARCHAR(50) NOT NULL,
  LastName VARCHAR(50) NOT NULL,
  PhoneNo VARCHAR(20) NOT NULL,
  OnTrip INT(4) NOT NULL DEFAULT 0,
  PRIMARY KEY (PassengerID),
  FOREIGN KEY fk_Ps_Passenger(PassengerID) REFERENCES User(UserID)

);


CREATE TABLE Driver 
(
  DriverID INT(8) NOT NULL,
  FirstName VARCHAR(50) NOT NULL,
  LastName VARCHAR(50) NOT NULL,
  PhoneNo VARCHAR(20) NOT NULL,
  IcNO VARCHAR(20) NOT NULL,
  LicenseNo VARCHAR(20) NOT NULL,
  OnTrip INT(4) NOT NULL DEFAULT 0,
  PRIMARY KEY (DriverID),
  FOREIGN KEY fk_Dv_Driver(DriverID) REFERENCES User(UserID)

);




CREATE TABLE Trip
(
  TripID INT(4) NOT NULL AUTO_INCREMENT,
  PassengerID INT(8) NOT NULL,
  DriverID INT(8) DEFAULT NULL,
  StartTime TIMESTAMP DEFAULT NULL,
  EndTime TIMESTAMP DEFAULT NULL,
  StartPostal VARCHAR(6) DEFAULT NULL,
  EndPostal VARCHAR(6) DEFAULT NULL,
  PRIMARY KEY (TripID),
  FOREIGN KEY fk_TP_Passenger(PassengerID) REFERENCES Passenger(PassengerID),
  FOREIGN KEY fk_TP_Driver(DriverID) REFERENCES Driver(DriverID)
 );



insert into user (UserType, Email,Password) values ('passenger','123@123.com','123456');
insert into user (UserType, Email,Password) values ('driver','abc@123.com','123456');


insert into passenger (PassengerID,FirstName, LastName ,PhoneNo) values (1,'yuxuan','luo','11112222');
insert into driver (DriverID,FirstName, LastName ,PhoneNo,IcNO,LicenseNo) values (2,'kk','luo','11113222','12551155','s10002');

