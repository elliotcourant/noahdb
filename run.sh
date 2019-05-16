#!/usr/bin/env bash

/node/bin/postgres -D /node/pgdata &

/node/bin/noahdb start -a