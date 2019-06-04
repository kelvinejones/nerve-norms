class StimulusResponse extends Chart {
	constructor(participant, norms) {
		super([1, 20], [0.01, 20], Chart.scaleType.LOG, Chart.scaleType.LOG)
		this.participant = this.calculateParticipant(participant)
		this.norms = norms.SR.data
		this.sdFunc = Chart.logSD
		this.xSDName = 4
	}

	calculateParticipant(participant) {
		let peakResponse = 0;
		participant.sections.ExVars.data.forEach(function(exvar) {
			if (exvar[0] == 6) {
				peakResponse = exvar[1];
			}
		})

		return participant.sections.SR.data.data.map((d, i) => {
			return [
				d[1],
				(i + 1) * 2 / 100 * peakResponse,
			]
		})
	}

	get name() { return "Stimulus Response" }
	get xLabel() { return "Stimulus Current (mA)" }
	get yLabel() { return "Peak Response (mV)" }

	updateParticipant(participant) {
		this.participant = this.calculateParticipant(participant)
		this.animateXYLine(this.participant, "sr")
	}

	updateNorms(norms) {
		this.norms = norms.SR.data
		this.animateNorms(this.norms, "sr")
	}

	drawLines(svg) {
		this.createXYLine(this.participant, "sr")
		this.createNorms(this.norms, "sr")
		this.animateXYLine(this.participant, "sr")
		this.animateNorms(this.norms, "sr")
	}
}
