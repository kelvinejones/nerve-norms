class ChartFactory {
	constructor(participants, norms, outlierScores) {
		this.osAccessor = {
			participant: "",
			normative: "",
			getScores: function() {
				return outlierScores[this.participant][this.normative]
			},
		}

		this.partDropDown = new DataDropDown("select-participant-dropdown", participants, (name, currentParticipant) => {
			this.osAccessor.participant = name
			ExVars.update(this.osAccessor.getScores(), currentParticipant);
			plots.forEach(pl => {
				pl.chart.updateParticipant(currentParticipant.plots)
			})
		})

		this.normDropDown = new DataDropDown("select-normative-dropdown", norms, (name, currentNormative) => {
			this.osAccessor.normative = name
			ExVars.update(this.osAccessor.getScores(), this.partDropDown.data);
			plots.forEach(pl => {
				pl.chart.updateNorms(currentNormative.plots)
			})
		})

		this.osAccessor.participant = this.partDropDown.selection
		this.osAccessor.normative = this.normDropDown.selection

		// Create all of the plots
		const plots = [{
			chart: this.build("rc"),
			selector: "#recoveryCycle svg",
		}, {
			chart: this.build("te"),
			selector: "#thresholdElectrotonus svg",
		}, {
			chart: this.build("cd"),
			selector: "#chargeDuration svg",
		}, {
			chart: this.build("tiv"),
			selector: "#thresholdIV svg",
		}, {
			chart: this.build("sr"),
			selector: "#stimulusResponse svg",
		}, {
			chart: this.build("srel"),
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

	build(typeStr) {
		switch (typeStr) {
			case "rc":
				return new RecoveryCycle(this.partDropDown.data.plots, this.normDropDown.data.plots)
			case "te":
				return new ThresholdElectrotonus(this.partDropDown.data.plots, this.normDropDown.data.plots)
			case "cd":
				return new ChargeDuration(this.partDropDown.data.plots, this.normDropDown.data.plots)
			case "tiv":
				return new ThresholdIV(this.partDropDown.data.plots, this.normDropDown.data.plots)
			case "sr":
				return new StimulusResponse(this.partDropDown.data.plots, this.normDropDown.data.plots)
			case "srel":
				return new StimulusRelative(this.partDropDown.data.plots, this.normDropDown.data.plots)
		}
	}
}
