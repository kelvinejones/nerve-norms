function initPlots(data) {
	// Create all of the plots
	const plots = [{
		chart: new RecoveryCycle(data.plots),
		selector: "#recoveryCycle svg",
	}, {
		chart: new ThresholdElectrotonus(data.plots),
		selector: "#thresholdElectrotonus svg",
	}, {
		chart: new ChargeDuration(data.plots),
		selector: "#chargeDuration svg",
	}, {
		chart: new ThresholdIV(data.plots),
		selector: "#thresholdIV svg",
	}, {
		chart: new StimulusResponse(data.plots),
		selector: "#stimulusResponse svg",
	}, {
		chart: new StimulusRelative(data.plots),
		selector: "#stimulusResponseRelative svg",
	}, ]

	// Draw them all
	plots.forEach(pl => {
		pl.chart.draw(d3.select(pl.selector), true)
		pl.chart.setDelayTime(Chart.fastDelay) // After initial setup, remove the delay
	})

	// Now set all excitability variables
	ExVars.update(data);

	function changeParticipant(ev) {
		plots.forEach(pl => {
			currentParticipant = participants[ev.srcElement.value] // This is a global
			pl.chart.updatePlots(currentParticipant.plots)
			pl.chart.updateNorms(currentParticipant.plots)
			ExVars.update(currentParticipant)
		})
	}

	document.getElementById("select-participant-dropdown")
		.addEventListener("change", changeParticipant);
}
