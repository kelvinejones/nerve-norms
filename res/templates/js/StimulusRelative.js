class StimulusRelative extends Chart {
	constructor(participant) {
		super([0, 200], [0, 100])
		this.participant = this.calculateParticipant(participant.sections.SR.data.data)
		this.norms = this.participant
		this.ySDName = undefined
		this.yMeanName = 3
		this.xSDName = 1
		this.xMeanName = 0
	}

	calculateParticipant(data) {
		const stimFor50PercentMax = data[24][1] // Could also be extracted from excitability variables
		return data.map((d, i) => {
			return [
				d[1] / stimFor50PercentMax * 100,
				d[0],
			]
		})
	}

	get name() { return "Relative Stimulus Response" }
	get xLabel() { return "Stimulus (% Mean Threshold)" }
	get yLabel() { return "Peak Response (% Max)" }

	updateParticipant(participant) {
		this.participant = this.calculateParticipant(participant.sections.SR.data.data)
		this.animateXYLine(this.participant, "srel")
	}

	updateNorms(norms) {
		this.norms = norms.SRel.data
		this.animateNorms(this.norms, "srel")
	}

	drawLines(svg) {
		this.createXYLine(this.participant, "srel")
		this.createNorms(this.norms, "srel", false)
		this.animateXYLine(this.participant, "srel")
		this.animateNorms(this.norms, "srel", false)
	}
}
