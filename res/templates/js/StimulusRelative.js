class StimulusRelative extends Chart {
	constructor(participant, norms) {
		super([0, 200], [0, 100])
		this.participant = this.calculateParticipant(participant.sections.SR.data.data)
		this.norms = (norms === undefined) ? undefined : norms.SRel.data
		this.ySDIndex = undefined
		this.yMeanIndex = 3
		this.xSDIndex = 1
		this.xMeanIndex = 0
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
		const useSD = (this.norms !== undefined)
		const norms = (this.norms === undefined) ? this.participant : this.norms
		this.createXYLine(this.participant, "srel")
		this.createNorms(norms, "srel", useSD)
		this.animateXYLine(this.participant, "srel")
		this.animateNorms(norms, "srel", useSD)
	}
}
