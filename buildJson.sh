#!/bin/bash -e

echo "Combining data..."
go run cmd/combineData/combineData.go
echo -e "const participants = \n$(cat json/all.json)" > res/templates/data/participants.json


echo "Generating norms..."
go run cmd/filterJson/filterJson.go -norm "res/templates/data/normAll.json"
echo -e "normative[\"Human Norms\"] = \n$(cat res/templates/data/normAll.json)" > res/templates/data/normAll.json
go run cmd/filterJson/filterJson.go -sex M -minAge 25 -maxAge 35 -norm res/templates/data/normM30.json
echo -e "normative[\"M30 Norms\"] = \n$(cat res/templates/data/normM30.json)" > res/templates/data/normM30.json

echo "Calculating outlier scores..."
go run cmd/outlierScores/outlierScores.go -name "Rat on Drugs" -sex M -minAge 25 -maxAge 35 -output "res/templates/data/outlier-rat-M30.json"
echo -e "outlierScore[\"Rat on Drugs\"] = outlierScore[\"Rat on Drugs\"] || {};\noutlierScore[\"Rat on Drugs\"][\"M30 Norms\"] = \n$(cat res/templates/data/outlier-rat-M30.json)" > res/templates/data/outlier-rat-M30.json
go run cmd/outlierScores/outlierScores.go -name "Rat on Drugs" -output "res/templates/data/outlier-rat-all.json"
echo -e "outlierScore[\"Rat on Drugs\"] = outlierScore[\"Rat on Drugs\"] || {};\noutlierScore[\"Rat on Drugs\"][\"Human Norms\"] = \n$(cat res/templates/data/outlier-rat-all.json)" > res/templates/data/outlier-rat-all.json
go run cmd/outlierScores/outlierScores.go -name "CA-WI20S" -sex M -minAge 25 -maxAge 35 -output "res/templates/data/outlier-WI20S-M30.json"
echo -e "outlierScore[\"CA-WI20S\"] = outlierScore[\"CA-WI20S\"] || {};\noutlierScore[\"CA-WI20S\"][\"M30 Norms\"] = \n$(cat res/templates/data/outlier-WI20S-M30.json)" > res/templates/data/outlier-WI20S-M30.json
go run cmd/outlierScores/outlierScores.go -name "CA-WI20S" -output "res/templates/data/outlier-WI20S-all.json"
echo -e "outlierScore[\"CA-WI20S\"] = outlierScore[\"CA-WI20S\"] || {};\noutlierScore[\"CA-WI20S\"][\"Human Norms\"] = \n$(cat res/templates/data/outlier-WI20S-all.json)" > res/templates/data/outlier-WI20S-all.json
go run cmd/outlierScores/outlierScores.go -name "CA-CR21S" -sex M -minAge 25 -maxAge 35 -output "res/templates/data/outlier-CR21S-M30.json"
echo -e "outlierScore[\"CA-CR21S\"] = outlierScore[\"CA-CR21S\"] || {};\noutlierScore[\"CA-CR21S\"][\"M30 Norms\"] = \n$(cat res/templates/data/outlier-CR21S-M30.json)" > res/templates/data/outlier-CR21S-M30.json
go run cmd/outlierScores/outlierScores.go -name "CA-CR21S" -output "res/templates/data/outlier-CR21S-all.json"
echo -e "outlierScore[\"CA-CR21S\"] = outlierScore[\"CA-CR21S\"] || {};\noutlierScore[\"CA-CR21S\"][\"Human Norms\"] = \n$(cat res/templates/data/outlier-CR21S-all.json)" > res/templates/data/outlier-CR21S-all.json
