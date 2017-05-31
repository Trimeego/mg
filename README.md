# mg
Simple DB migration CLI based on http://github.com/mattes/migrate

## Installation

This installation assume

```
$ go get -u -d github.com/trimeego/mg
$ go build -tags 'mssql' -o /usr/local/bin/mgsql github.com/trimeego/mg
$ go build -tags 'postgres' -o /usr/local/bin/mgpostgres github.com/trimeego/mg
```

## Usage

### Creating New Migration Files

`mg` can create proper migration files.  In order to create a migration file, you supply the intent behind the script in snake case, such as `create_user_table`.

```
$ mg create create_user_table
```

This will create timestamped `.up.sql` and `.down.sql` files.

### Configuration

For convenience, you can place a `.mg.yaml` file in the migration folder with the `env` options that you want to make available.  This is strongly recommended to avoid mistakes.

Adding the following `.mg.yaml` allows provides `dev`, `qa`, and `prod` env options:

```
dev: postgres://icgadmin:password@localhost:5432/vantage
qa: postgres://icgadmin:password@test-postgres.caxctts3oq6h.us-west-2.rds.amazonaws.com:5432/vantage
prod: postgres://icgadmin:password@prod-postgres.ctuifd37tvce.us-west-2.rds.amazonaws.com:5432/vantage
```

These can be used anywhere using the `--env` flag.  Given this config file, the following two commands are equivalent.

```
$ mg version --url postgres://icgadmin:password@localhost:5432/vantage
$ mg version --env dev
```

### Getting the Current Version

To get the current version:

```
$ mg version --env dev
```

 
### Migrating Up

To migrate up to the current level.

```
$ mg up --env dev
```

To migrate up a specific number of steps from the current version:

```
$ mg up 3 --env dev
```


### Migrating Down

To migrate down a specific number of steps from the current version:

```
$ mg down 3 --env dev
```


### Migrating to a Specific Version, regardless of direction

To migrate to a specific version from the current version:

```
$ mg goto 1479382989 --env dev
```


### Forcing a Version

In some cases, particularly during development, one must force the current level to a particular level.  This is particularly true when a new script fails to do a syntax error.  In this case, you will need to force the migration to think that it is at a specific version.

To migrate down a specific number of steps from the current version:

```
$ mg force 1478293889 --env dev
```






