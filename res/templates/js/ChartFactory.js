class ChartFactory {
	constructor(participants, norms, outlierScores) {
		this.osAccessor = {
			participant: "",
			normative: "",
			getScores: function() {
				return outlierScores[this.participant][this.normative]
			},
		}
		this.url = "https://us-central1-nervenorms.cloudfunctions.net/"

		this.partDropDown = new DataDropDown("select-participant-dropdown", participants, (name, currentParticipant) => {
			this.osAccessor.participant = name
			ExVars.setToZero()
			Object.values(plots).forEach(pl => {
				pl.updateParticipant(currentParticipant)
			})
			ExVars.updateValues(currentParticipant)

			fetch(this.url + "outliers?name=" + this.osAccessor.participant)
				.then(function(response) {
					return response.json()
				})
				.then(function(myJson) {
					ExVars.updateScores(myJson)
				})
		}, ["CA-CR21S", "CA-WI20S", "Rat on Drugs", ])

		this.normDropDown = new DataDropDown("select-normative-dropdown", norms, (name, currentNormative) => {
			this.osAccessor.normative = name
			ExVars.setToZero()
			Object.values(plots).forEach(pl => {
				pl.updateNorms(currentNormative)
			})
			ExVars.updateValues(this.partDropDown.data);

			fetch("https://us-central1-nervenorms.cloudfunctions.net/outliers?name=" + this.osAccessor.participant)
				.then(function(response) {
					return response.json();
				})
				.then(function(myJson) {
					ExVars.updateScores(myJson);
				});
		}, ["Human Norms", "M30 Norms", ])

		this.osAccessor.participant = this.partDropDown.selection
		this.osAccessor.normative = this.normDropDown.selection

		this.applyFilter = (event) => {
			fetch(this.url + "norms" + Filter.asQueryString())
				.then(function(response) {
					return response.json()
				})
				.then(function(norms) {
					Object.values(plots).forEach(pl => {
						pl.updateNorms(norms)
					})
				})
			event.preventDefault()
		}

		document.querySelector("form").addEventListener("submit", this.applyFilter)

		const plots = {
			"recoveryCycle": null,
			"thresholdElectrotonus": null,
			"chargeDuration": null,
			"thresholdIV": null,
			"stimulusResponse": null,
			"stimulusResponseRelative": null,
		}

		Object.keys(plots).forEach(key => {
			plots[key] = this.build(key)
			plots[key].draw(d3.select("#" + key + " svg"), true)
			plots[key].setDelayTime(Chart.fastDelay) // After initial setup, remove the delay
		})

		// Now set all excitability variables
		ExVars.updateScores(this.osAccessor.getScores())
		ExVars.updateValues(this.partDropDown.data)
	}

	drawModal(typeStr) {
		const chart = this.build(typeStr)
		chart.setDelayTime(Chart.fastDelay).setTransitionTime(Chart.fastTransition)

		document.getElementById('modal-title').innerHTML = chart.name
		d3.selectAll("#modal svg > *").remove()

		chart.draw(d3.select('#modal svg'))
		$('#modal').modal('toggle')
	}

	build(typeStr) {
		switch (typeStr) {
			case "recoveryCycle":
				return new RecoveryCycle(this.partDropDown.data, this.normDropDown.data)
			case "thresholdElectrotonus":
				return new ThresholdElectrotonus(this.partDropDown.data, this.normDropDown.data)
			case "chargeDuration":
				return new ChargeDuration(this.partDropDown.data, this.normDropDown.data)
			case "thresholdIV":
				return new ThresholdIV(this.partDropDown.data, this.normDropDown.data)
			case "stimulusResponse":
				return new StimulusResponse(this.partDropDown.data, this.normDropDown.data)
			case "stimulusResponseRelative":
				return new StimulusRelative(this.partDropDown.data, this.normDropDown.data)
		}
	}
}
