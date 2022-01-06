mysql -u root hackthon_go -e "delete from sales;"
mysql -u root hackthon_go -e "ALTER TABLE sales AUTO_INCREMENT = 1;"

mysql -u root hackthon_go -e "delete from invoices;"
mysql -u root hackthon_go -e "ALTER TABLE invoices AUTO_INCREMENT = 1;"

mysql -u root hackthon_go -e "delete from products;"
mysql -u root hackthon_go -e "ALTER TABLE products AUTO_INCREMENT = 1;"

mysql -u root hackthon_go -e "delete from customers;"
mysql -u root hackthon_go -e "ALTER TABLE customers AUTO_INCREMENT = 1;"