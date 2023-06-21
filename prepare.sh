#!/bin/bash

apt-get update
apt-get install -y mysql-server

service mysql start

mysql -e "ALTER USER 'root'@'localhost' IDENTIFIED BY 'ihackstuff';"

mysql -e "CREATE DATABASE effat;"

mysql -e "USE effat;"

mysql -e "CREATE TABLE divers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    diverEqp VARCHAR(255) NOT NULL
);"

echo "MySQL setup completed."
