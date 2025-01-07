#!/bin/bash
for file in *.sql; do
  echo "Importing $file..."
  psql -U testnet -d testv2 -f "$file"
done