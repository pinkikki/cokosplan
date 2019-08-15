#!/bin/bash

set -eou pipefail

curl -XPUT 'localhost:9200/_template/product' -H 'Content-type: application/json' -d @product.json
curl -XPUT 'localhost:9200/_template/mapping-expense' -H 'Content-type: application/json' -d @mapping-expense.json
