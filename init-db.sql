CREATE TABLE Book (
  ID serial PRIMARY KEY,
  Title varchar(200) DEFAULT NULL,
  Author varchar(200) DEFAULT NULL,
  Publisher varchar(200) DEFAULT NULL,
  PublishDate timestamp NULL DEFAULT NULL,
  Rating float DEFAULT NULL,
  Status int
);