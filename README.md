# containerd-compose (WIP)

**containerd-compose** is a tool for defining and running multi-container applications. 
With **containerd-compose**, you use a YAML file (e.g. the docker-compose.yml) to configure your application's services. 
Then, using a single command, you create and start all the services from your configuration.

This program is designed to be compatible with docker-compose files, but we have not yet implemented all the features. 
Please check the [Supported Parameters section](https://github.com/mc256/containerd-compose#supported-parameters).

### Install
Please install [containerd](https://github.com/containerd/containerd) first.


```
git clone https://github.com/mc256/containerd-compose
cd ./containerd-compose
go install
```


### Supported Parameters

- services
    - image
    - volumes (short syntax only)
    - volumes_from
    - environment
    


### Example
**containerd-compose** also supports `.env` file. You can write `$ENVIRONMENT_VARIABLES` in your compose file.

Example for [nextcloud](https://github.com/nextcloud/docker):
```yaml
version: '2'

volumes:
  nextcloud:
  db:
  dbrun:

services:
  db:
    image: mariadb
    command: --transaction-isolation=READ-COMMITTED --binlog-format=ROW
    restart: always
    volumes:
      - db:/var/lib/mysql:rw
      - dbrun:/var/run:rw
    environment:
      - MYSQL_ROOT_PASSWORD=$PASSWORD
      - MYSQL_PASSWORD=$PASSWORD
      - MYSQL_DATABASE=nextcloud
      - MYSQL_USER=nextcloud
  app:
    image: nextcloud:fpm
    links:
      - db
    volumes:
      - nextcloud:/var/www/html

  web:
    image: nginx
    ports:
      - 8080:80
    links:
      - app
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
    volumes_from:
      - app
```

