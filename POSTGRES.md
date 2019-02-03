PostgreSQL

---
Dif postgres utilities

utilities are used outside of the postgres intance, once issued, you can then log into postgres to see the changes

createuser: creates a user (use --createdb etc to add permissions at the time of creation)
createdb: creates a database
dropuser: deletes a user
dropdb: deletes a database
postgres: executes the SQL server itself (we saw that one above when we checked our Postgres version!)
pg_dump: dumps the contents of a single database to a file
pg_dumpall: dumps all databases to a file
psql: we recognize that one!
---


first we install postgres with brew

brew install postgresql

once installed follow the instructions in your terminal to start the postgres service 
OR
pg_ctl -D /usr/local/var/postgres start && brew services start postgresql

now well connecgt to the local postgres instance so we can configure it using psql cli tool

psql postgres will connect to the local postgres instance running 

lets list the users  using \du

postgres=# \du
                                     List of roles
  Role name   |                         Attributes                         | Member of
--------------+------------------------------------------------------------+-----------
 marvin.matos | Superuser, Create role, Create DB, Replication, Bypass RLS | {}

postgres does not manage users and groups like most permissions models do. postgres directly manages roles

Using the default user is bad practice because this is the superuser to the db and can things like delete dbs etc... you should create a role with minimal permissions 
in postgres there are no users groups there are roles

to create roles there are two ways, by using the sql CREATE ROLE or using the wrapper utility createuser which uses the CREATE ROLE query

To change the password of a role you can use the following:
\password <role name>
\password marvin.matos

you will be prompted to enter your password and confirm.

lets create a more restrictive user

Remember there are two ways to do this, lets create the first way

CREATE ROLE marvin WITH LOGIN PASSWORD 'test';

\du
                                     List of roles
  Role name   |                         Attributes                         | Member of
--------------+------------------------------------------------------------+-----------
 marvin       |                                                            | {}
 marvin.matos | Superuser, Create role, Create DB, Replication, Bypass RLS | {}

notice the new role witout any attributes, this is because we did not configure them yet
postgres will create this user without admin privs the roles will only have access to tables dbs or rows that it has permissions for

lets add the CREATEDB permission to our role with the ALTER ROLE directive


ALTER ROLE marvin CREATEDB;

postgres=# \du
                                     List of roles
  Role name   |                         Attributes                         | Member of
--------------+------------------------------------------------------------+-----------
 marvin       | Create DB                                                  | {}
 marvin.matos | Superuser, Create role, Create DB, Replication, Bypass RLS | {}


now with the create user utility.
TODO


CREATE DATABASE

using psql we will create a new database

well sign in with our new account 
psql postgres -U marvin

CREATE DATABASE api_server;

once you create the db you will need to have at least one user who has permissions to access the db, lets grand privs to our new restricted user

GRANT ALL PRIVILEGES ON DATABASE api_server TO marvin;

listing the current databases
postgres=> \list
                                          List of databases
    Name    |    Owner     | Encoding |   Collate   |    Ctype    |         Access privileges
------------+--------------+----------+-------------+-------------+-----------------------------------
 api_server | marvin       | UTF8     | en_US.UTF-8 | en_US.UTF-8 | =Tc/marvin                       +
            |              |          |             |             | marvin=CTc/marvin
 postgres   | marvin.matos | UTF8     | en_US.UTF-8 | en_US.UTF-8 |
 template0  | marvin.matos | UTF8     | en_US.UTF-8 | en_US.UTF-8 | =c/"marvin.matos"                +
            |              |          |             |             | "marvin.matos"=CTc/"marvin.matos"
 template1  | marvin.matos | UTF8     | en_US.UTF-8 | en_US.UTF-8 | =c/"marvin.matos"                +
            |              |          |             |             | "marvin.matos"=CTc/"marvin.matos"
(4 rows)

now well connect to the new database

postgres=> \connect api_server
You are now connected to database "api_server" as user "marvin".
api_server=>



nycMMATOSmbp:autobahn marvin.matos$ createuser marvin_svc --createdb
nycMMATOSmbp:autobahn marvin.matos$ psql postgres -U marvin
psql (11.1)
Type "help" for help.

postgres=> \du
                                     List of roles
  Role name   |                         Attributes                         | Member of
--------------+------------------------------------------------------------+-----------
 marvin       | Create DB                                                  | {}
 marvin.matos | Superuser, Create role, Create DB, Replication, Bypass RLS | {}
 marvin_svc   | Create DB                                                  | {}



nycMMATOSmbp:autobahn marvin.matos$ createdb api_server_v2 -U marvin
nycMMATOSmbp:autobahn marvin.matos$ psql postgres -U marvin
psql (11.1)
Type "help" for help.

postgres=> \list
                                            List of databases
     Name      |    Owner     | Encoding |   Collate   |    Ctype    |         Access privileges
---------------+--------------+----------+-------------+-------------+-----------------------------------
 api_server    | marvin       | UTF8     | en_US.UTF-8 | en_US.UTF-8 | =Tc/marvin                       +
               |              |          |             |             | marvin=CTc/marvin
 api_server_v2 | marvin       | UTF8     | en_US.UTF-8 | en_US.UTF-8 |
....
(5 rows)

postgres=> GRANT ALL PRIVILEGES ON DATABASE api_server_v2 TO marvin;
GRANT
postgres=> \list
                                            List of databases
     Name      |    Owner     | Encoding |   Collate   |    Ctype    |         Access privileges
---------------+--------------+----------+-------------+-------------+-----------------------------------
 api_server    | marvin       | UTF8     | en_US.UTF-8 | en_US.UTF-8 | =Tc/marvin                       +
               |              |          |             |             | marvin=CTc/marvin
 api_server_v2 | marvin       | UTF8     | en_US.UTF-8 | en_US.UTF-8 | =Tc/marvin                       +
               |              |          |             |             | marvin=CTc/marvin
 postgres      | marvin.matos | UTF8     | en_US.UTF-8 | en_US.UTF-8 |
 template0     | marvin.matos | UTF8     | en_US.UTF-8 | en_US.UTF-8 | =c/"marvin.matos"                +
               |              |          |             |             | "marvin.matos"=CTc/"marvin.matos"
 template1     | marvin.matos | UTF8     | en_US.UTF-8 | en_US.UTF-8 | =c/"marvin.matos"                +
               |              |          |             |             | "marvin.matos"=CTc/"marvin.matos"
(5 rows)

postgres=> \q

