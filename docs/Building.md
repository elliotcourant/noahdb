# Building From Source

To build noahdb from source make sure you have the required software mentioned in the docs readme.

Start with cloning this repository.

```bash
git clone git@github.com:elliotcourant/noahdb.git
cd noahdb
```

Then some of noah's dependencies will need to be setup. Since noahdb uses go mod this can be done
by simply running:

```bash
GO111MODULE=on bash -c 'go mod vendor'
```

This will download all the needed dependencies into the vendor folder in the project's directory.
The `GO111MODULE=on` part will allow go mod to be run from within your GOPATH if that's where
you are building noahdb.

Some parts of noahdb are not included in the repository and should be generated or embedded.
To do this run:

```bash
make generated
``` 

This creates several strings files for enums which improves logging. It also generates go files from
the `.proto` files that are used for most of noahdb's internal data types.
The last thing it does is create an embeddable version of `internal_sql.sql` found in the core/files
directory.

At this point you can either run:

```bash
make build
```

Which will create a `bin` folder in the project's directory and will create the `noahdb` executable
within that folder.

Or you can run:

```bash
make test
```

Which will run all of noahdb's unit tests. Please note that to run tests a local instance of docker
needs to be running as noah uses docker to create separate clean environments easily for each test.

Note: both the build and test make commands will run the generated command as part of them. However
generated should be run if you want to do anything with the code before building.