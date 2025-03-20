## MariaDB

> And normally when I'm working with Go, my default database is Postgres because I like Postgres a lot, but I already have several courses up on Udemy with Postgres as the default database. So I thought this time around I'd use Mariadb and Mariadb is effectively a drop in replacement for MySQL. But if you don't Mariadb is, in my humble opinion, a better version of MySQL.

Mac: https://mariadb.com/kb/en/installing-mariadb-on-macos-using-homebrew/

Win: https://community.chocolatey.org/packages/mariadb

Docker: https://mariadb.com/kb/en/installing-and-using-mariadb-via-docker/

## DB client app 

Mac app: Sequel Ace - free 

Win app: 

* HeidiSQL - a lot of people recommend it  https://www.heidisql.com/
* DBeaver

### Maria DB connection

```
go get github.com/go-sql-driver/mysql
```

Login via terminal

```
mysql -u root -p
```

Permissions

```
grant all on widgets.* to 'name'@'%' identified by 'secret';
```

Эта команда используется в **MariaDB** и **MySQL** (в MySQL начиная с версии 8.0 синтаксис немного изменился).

#### Разбор параметров:

- **`GRANT ALL`** — предоставляет все возможные привилегии (например, `SELECT`, `INSERT`, `UPDATE`, `DELETE`, `CREATE`, `DROP` и т. д.) пользователю.
- **`ON widgets.\*`** — означает, что привилегии выдаются на все таблицы (`*`) в базе данных `widgets`.
- **`TO 'name@%'`** — создает или изменяет пользователя с именем `name`, который может подключаться с любого хоста (`%` — это wildcard, означающий «любой хост»).
- **`IDENTIFIED BY 'secret'`** — устанавливает пароль `secret` для пользователя `name`.

#### Только ли в MariaDB?

1. В **MariaDB** этот синтаксис работает без изменений.

2. В **MySQL** до версии 8.0 тоже работало аналогично.

3. В MySQL 8.0+ использование `IDENTIFIED BY` в `GRANT` устарело. Теперь нужно сначала создать пользователя с `CREATE USER`, а затем выдать ему права:

   ```
   sqlCopyEditCREATE USER 'name'@'%' IDENTIFIED BY 'secret';
   GRANT ALL ON widgets.* TO 'name'@'%';
   ```

Если используешь **MySQL 8+**, лучше писать так. В MariaDB оба варианта работают.



### DSN connection string (for .env)

```
DSN=username:passwrod@tcp(127.0.0.1:3306)/widgets?parseTime=true&tls=false
```

### Coalesce 

```sql
select id, name, description, inventory_level, price, created_at, updated_at, coalesce(image, '')
from widgets where id = ?
```

Is there may be situations where we don't have an image and what I could do is take advantage of the SQL null types that are built into our database environment or built into go. But what I'll do instead is just say **coalesce**. It just saves some programming.

