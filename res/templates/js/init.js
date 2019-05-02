function initPlots(data) {
	// Create all of the plots

	var rc = new RecoveryCycle(data.plots[7].data)
	rc.draw(d3.select("#recoveryCycle svg"), true)

	var te = new ThresholdElectrotonus(data.plots[2].data, data.plots[3].data, data.plots[4].data, data.plots[5].data)
	te.draw(d3.select("#thresholdElectrotonus svg"), true)

	var cd = new ChargeDuration(data.plots[1].data)
	cd.draw(d3.select("#chargeDuration svg"), true)

	var tiv = new ThresholdIV(data.plots[6].data)
	tiv.draw(d3.select("#thresholdIV svg"), true)

	var sr = new StimulusResponse(data.plots[0].data)
	sr.draw(d3.select("#stimulusResponse svg"), true)

	var srrel = new StimulusRelative(data.plots[0].data)
	srrel.draw(d3.select("#stimulusResponseRelative svg"), true)

	// Now set all excitability variables

	function setExcitabilityVariable(idString, value, score) {
		var row = document.getElementById(idString);
		row.getElementsByClassName("excite-value")[0].innerHTML = value
		// TODO: update visualization using score
	}

	setExcitabilityVariable("overall-score", data.outlierScore, data.outlierScore)

	data.plots.map(function(pl) { return pl.discreteMeasures; }).flat()
		.concat(data.discreteMeasures)
		.forEach(function(exind) {
			setExcitabilityVariable("qtrac-excite-" + exind.qtracExciteID, exind.value, exind.outlierScore)
		})
}
