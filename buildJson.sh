#!/bin/bash -e

echo "Combining data..."
go run cmd/combineData/combineData.go

echo "Generating norms..."
go run cmd/filterJson/filterJson.go -norm "res/templates/data/normAll.json"
echo -e "normative[\"Human Norms\"] = \n$(cat res/templates/data/normAll.json)" > res/templates/data/normAll.json
go run cmd/filterJson/filterJson.go -sex M -minAge 25 -maxAge 35 -norm res/templates/data/normM30.json
echo -e "normative[\"M30 Norms\"] = \n$(cat res/templates/data/normM30.json)" > res/templates/data/normM30.json

# echo "Calculating outlier scores..."
go run cmd/outlierScores/outlierScores.go -name "Rat on Drugs" -sex M -minAge 25 -maxAge 35 -output "res/templates/data/outlier-rat-M30.json"
echo -e "outlierScore[\"Rat on Drugs\"] = outlierScore[\"Rat on Drugs\"] || {};\noutlierScore[\"Rat on Drugs\"][\"M30 Norms\"] = \n$(cat res/templates/data/outlier-rat-M30.json)" > res/templates/data/outlier-rat-M30.json
go run cmd/outlierScores/outlierScores.go -name "Rat on Drugs" -output "res/templates/data/outlier-rat-all.json"
echo -e "outlierScore[\"Rat on Drugs\"] = outlierScore[\"Rat on Drugs\"] || {};\noutlierScore[\"Rat on Drugs\"][\"Human Norms\"] = \n$(cat res/templates/data/outlier-rat-all.json)" > res/templates/data/outlier-rat-all.json
go run cmd/outlierScores/outlierScores.go -name "CA-WI20S" -sex M -minAge 25 -maxAge 35 -output "res/templates/data/outlier-WI20S-M30.json"
echo -e "outlierScore[\"WI20S\"] = outlierScore[\"WI20S\"] || {};\noutlierScore[\"WI20S\"][\"M30 Norms\"] = \n$(cat res/templates/data/outlier-WI20S-M30.json)" > res/templates/data/outlier-WI20S-M30.json
go run cmd/outlierScores/outlierScores.go -name "CA-WI20S" -output "res/templates/data/outlier-WI20S-all.json"
echo -e "outlierScore[\"WI20S\"] = outlierScore[\"WI20S\"] || {};\noutlierScore[\"WI20S\"][\"Human Norms\"] = \n$(cat res/templates/data/outlier-WI20S-all.json)" > res/templates/data/outlier-WI20S-all.json
go run cmd/outlierScores/outlierScores.go -name "CA-CR21S" -sex M -minAge 25 -maxAge 35 -output "res/templates/data/outlier-CR21S-M30.json"
echo -e "outlierScore[\"CR21S\"] = outlierScore[\"CR21S\"] || {};\noutlierScore[\"CR21S\"][\"M30 Norms\"] = \n$(cat res/templates/data/outlier-CR21S-M30.json)" > res/templates/data/outlier-CR21S-M30.json
go run cmd/outlierScores/outlierScores.go -name "CA-CR21S" -output "res/templates/data/outlier-CR21S-all.json"
echo -e "outlierScore[\"CR21S\"] = outlierScore[\"CR21S\"] || {};\noutlierScore[\"CR21S\"][\"Human Norms\"] = \n$(cat res/templates/data/outlier-CR21S-all.json)" > res/templates/data/outlier-CR21S-all.json
