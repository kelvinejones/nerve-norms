function initPlots(data) {
	// Create all of the plots
	const plots = [{
		chart: new RecoveryCycle(data.plots[7].data),
		selector: "#recoveryCycle svg",
	}, {
		chart: new ThresholdElectrotonus(data.plots[2].data, data.plots[3].data, data.plots[4].data, data.plots[5].data),
		selector: "#thresholdElectrotonus svg",
	}, {
		chart: new ChargeDuration(data.plots[1].data),
		selector: "#chargeDuration svg",
	}, {
		chart: new ThresholdIV(data.plots[6].data),
		selector: "#thresholdIV svg",
	}, {
		chart: new StimulusResponse(data.plots[0].data),
		selector: "#stimulusResponse svg",
	}, {
		chart: new StimulusRelative(data.plots[0].data),
		selector: "#stimulusResponseRelative svg",
	}, ]

	// Draw them all
	plots.forEach(pl => {
		pl.chart.draw(d3.select(pl.selector), true)
	})

	let opacity = 0.8,
		red = d3.hsl("red"),
		green = d3.hsl("lightgreen");

	red.opacity = opacity
	green.opacity = opacity

	let interpolate = d3.interpolateHsl(green, red);

	// Now set all excitability variables

	function setExcitabilityVariable(idString, value, score) {
		var row = document.getElementById(idString);
		row.getElementsByClassName("excite-value")[0].innerHTML = value

		percent = score * 100
		color = interpolate(score)

		row.style.background = "linear-gradient(to right, " + color + " " + percent + "%, #ffffff 0%)"
	}

	setExcitabilityVariable("overall-score", data.outlierScore, data.outlierScore)

	data.plots.map(function(pl) { return pl.discreteMeasures; }).flat()
		.concat(data.discreteMeasures)
		.forEach(function(exind) {
			setExcitabilityVariable("qtrac-excite-" + exind.qtracExciteID, exind.value, exind.outlierScore)
		})
}
