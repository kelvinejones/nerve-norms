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
			ExVars.setScoresToZero()
			Object.values(this.plots).forEach(pl => {
				pl.updateParticipant(currentParticipant)
			})
			ExVars.updateValues(currentParticipant)
			this.updateOutliers(this.osAccessor.participant)
		}, ["CA-CR21S", "CA-WI20S", "Rat on Drugs", ])

		this.normDropDown = new DataDropDown("select-normative-dropdown", norms, (name, currentNormative) => {
			this.osAccessor.normative = name
			ExVars.setScoresToZero()
			Object.values(this.plots).forEach(pl => {
				pl.updateNorms(currentNormative)
			})
			ExVars.updateValues(this.partDropDown.data);
			this.updateOutliers(this.osAccessor.participant)
		}, ["Human Norms", "M30 Norms", ])

		this.osAccessor.participant = this.partDropDown.selection
		this.osAccessor.normative = this.normDropDown.selection

		this.applyFilter = (event) => {
			ExVars.setScoresToZero()

			const queryString = Filter.asQueryString()
			this.updateNorms(queryString)
			this.updateOutliers(this.osAccessor.participant, queryString)

			event.preventDefault()
		}

		document.querySelector("form").addEventListener("submit", this.applyFilter)

		this.plots = {
			"recoveryCycle": null,
			"thresholdElectrotonus": null,
			"chargeDuration": null,
			"thresholdIV": null,
			"stimulusResponse": null,
			"stimulusResponseRelative": null,
		}

		Object.keys(this.plots).forEach(key => {
			this.plots[key] = this.build(key)
			this.plots[key].draw(d3.select("#" + key + " svg"), true)
			this.plots[key].setDelayTime(Chart.fastDelay) // After initial setup, remove the delay
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
				return new RecoveryCycle(this.partDropDown.data)
			case "thresholdElectrotonus":
				return new ThresholdElectrotonus(this.partDropDown.data)
			case "chargeDuration":
				return new ChargeDuration(this.partDropDown.data)
			case "thresholdIV":
				return new ThresholdIV(this.partDropDown.data)
			case "stimulusResponse":
				return new StimulusResponse(this.partDropDown.data)
			case "stimulusResponseRelative":
				return new StimulusRelative(this.partDropDown.data)
		}
	}

	updateOutliers(name, queryString) {
		if (queryString === undefined) {
			queryString = Filter.asQueryString()
		}

		fetch(this.url + "outliers" + queryString + "&name=" + name)
			.then(response => {
				return response.json()
			})
			.then(scores => {
				ExVars.updateScores(scores)
			})
	}

	updateNorms(queryString) {
		if (queryString === undefined) {
			queryString = Filter.asQueryString()
		}

		fetch(this.url + "norms" + queryString)
			.then(response => {
				return response.json()
			})
			.then(norms => {
				Object.values(this.plots).forEach(pl => {
					pl.updateNorms(norms)
				})
			})
	}
}
