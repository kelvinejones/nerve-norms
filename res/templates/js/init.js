function initPlots(participants) {
	const osAccessor = {
		participant: "",
		getScores: function() {
			return participants[this.participant]
		},
	}

	const partDropDown = new DataDropDown("select-participant-dropdown", participants, function(name, currentParticipant) {
		osAccessor.participant = name
		ExVars.update(osAccessor.getScores());
		plots.forEach(pl => {
			pl.chart.updatePlots(currentParticipant.plots)
			pl.chart.updateNorms(currentParticipant.plots)
		})
	})

	osAccessor.participant = partDropDown.selection

	// Create all of the plots
	const plots = [{
		chart: new RecoveryCycle(partDropDown.data.plots),
		selector: "#recoveryCycle svg",
	}, {
		chart: new ThresholdElectrotonus(partDropDown.data.plots),
		selector: "#thresholdElectrotonus svg",
	}, {
		chart: new ChargeDuration(partDropDown.data.plots),
		selector: "#chargeDuration svg",
	}, {
		chart: new ThresholdIV(partDropDown.data.plots),
		selector: "#thresholdIV svg",
	}, {
		chart: new StimulusResponse(partDropDown.data.plots),
		selector: "#stimulusResponse svg",
	}, {
		chart: new StimulusRelative(partDropDown.data.plots),
		selector: "#stimulusResponseRelative svg",
	}, ]

	// Draw them all
	plots.forEach(pl => {
		pl.chart.draw(d3.select(pl.selector), true)
		pl.chart.setDelayTime(Chart.fastDelay) // After initial setup, remove the delay
	})

	// Now set all excitability variables
	ExVars.update(osAccessor.getScores());
}
