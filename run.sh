#!/usr/bin/env bash

/node/bin/postgres -D /data/db &
# test
/node/bin/noahdb start -dA -s /data/db/cluster