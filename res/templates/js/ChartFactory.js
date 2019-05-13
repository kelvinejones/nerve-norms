class ChartFactory {
	constructor(participants, norms, outlierScores) {
		this.osAccessor = {
			participant: "",
			normative: "",
			getScores: function() {
				return outlierScores[this.participant][this.normative]
			},
		}

		this.partDropDown = new DataDropDown("select-participant-dropdown", participants, function(name, currentParticipant) {
			this.osAccessor.participant = name
			ExVars.update(this.osAccessor.getScores(), currentParticipant);
			plots.forEach(pl => {
				pl.chart.updateParticipant(currentParticipant.plots)
			})
		})

		this.normDropDown = new DataDropDown("select-normative-dropdown", norms, function(name, currentNormative) {
			this.osAccessor.normative = name
			ExVars.update(this.osAccessor.getScores(), this.partDropDown.data);
			plots.forEach(pl => {
				pl.chart.updateNorms(currentNormative.plots)
			})
		})

		this.osAccessor.participant = this.partDropDown.selection
		this.osAccessor.normative this.normDropDown.selection

		// Create all of the plots
		const plots = [{
			chart: new RecoveryCycle(this.partDropDown.data.plotsthis.normDropDown.data.plots),
			selector: "#recoveryCycle svg",
		}, {
			chart: new ThresholdElectrotonus(this.partDropDown.data.plotsthis.normDropDown.data.plots),
			selector: "#thresholdElectrotonus svg",
		}, {
			chart: new ChargeDuration(this.partDropDown.data.plotsthis.normDropDown.data.plots),
			selector: "#chargeDuration svg",
		}, {
			chart: new ThresholdIV(this.partDropDown.data.plotsthis.normDropDown.data.plots),
			selector: "#thresholdIV svg",
		}, {
			chart: new StimulusResponse(this.partDropDown.data.plotsthis.normDropDown.data.plots),
			selector: "#stimulusResponse svg",
		}, {
			chart: new StimulusRelative(this.partDropDown.data.plotsthis.normDropDown.data.plots),
			selector: "#stimulusResponseRelative svg",
		}, ]

		// Draw them all
		plots.forEach(pl => {
			pl.chart.draw(d3.select(pl.selector), true)
			pl.chart.setDelayTime(Chart.fastDelay) // After initial setup, remove the delay
		})

		// Now set all excitability variables
		ExVars.update(this.osAccessor.getScores(), this.partDropDown.data);
	}
}
