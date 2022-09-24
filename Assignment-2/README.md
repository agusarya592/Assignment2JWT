# Hacktiv8 Assignment 2

## SQL

````bash
CREATE TABLE IF NOT EXISTS `order`(
	order_id INT NOT NULL AUTO_INCREMENT,
	customer_name VARCHAR (255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE NOW(),
    PRIMARY KEY (order_id)
);

CREATE TABLE IF NOT EXISTS `item`(
    item_id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    item_code VARCHAR(25) NOT NULL,
    description TEXT NOT NULL,
    quantity INT NOT NULL,
    order_id INT NOT null,
    FOREIGN KEY (order_id) REFERENCES `order`(order_id)
);

CREATE TABLE `user` (
    userID int not null auto_increment,
    username varchar(255) not null,
    email varchar(255) not null,
    password varchar(255) not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now() on update now(),
    primary key(userID)
);
## How To Run
```bash
go run main.go
````

## Notes

**Important**
Change the _.env.example_ name to _.env_

```

```
