---
persist is an abstraction on top of db components

The following environment variables are used to configure the database

Only configured for and tested with postgres, mysql.

```
    SQL_DATABASE=gorilla
    SQL_DRIVER=mysql
    SQL_HOST=localhost
    SQL_PASSWORD=gorilla
    SQL_PORT=5432
    SQL_USER=gorilla
```


```
    #!/bin/bash
    # bash configuration might look like:
    
    if [[ ${SQL_DRIVER} == postgres ]]; then
        export SQL_DRIVER=postgres
        export SQL_PORT=5432
        export SQL_HOST=localhost
    else
        export SQL_DRIVER=mysql
        export SQL_HOST=localhost
        export SQL_PORT=3306
    fi
    export SQL_DATABASE=gorilla
    export SQL_USER=gorilla
    export SQL_PASSWORD=gorilla
```
