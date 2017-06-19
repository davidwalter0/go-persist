CREATE USER 'gorilla'@'*' IDENTIFIED BY 'gorilla';
CREATE DATABASE gorilla CHARACTER SET utf8 COLLATE utf8_bin;
USE gorilla;
GRANT ALL ON gorilla.* TO 'gorilla'@'localhost' identified by 'gorilla';
GRANT ALL ON gorilla.* TO 'gorilla'@'%' identified by 'gorilla';
FLUSH PRIVILEGES;
