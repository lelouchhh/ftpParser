#!/bin/sh
for i in {10..31}
do
  go run cmd/main.go  -day=202207"${i}" i00 -outputXML=./cmd/xmls/ -outputZip=./cmd/zips/

done