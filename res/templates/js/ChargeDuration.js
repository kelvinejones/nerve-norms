class ChargeDuration extends Chart {
	constructor(participant, norms) {
		super([0, 1], [0, 10])
		this.participant = participant.sections.CD.data
		this.norms = (norms === undefined) ? undefined : norms.CD.data
	}

	get name() { return "Charge Duration" }
	get xLabel() { return "Stimulus Width (ms)" }
	get yLabel() { return "Threshold Change (mA•ms)" }

	updateParticipant(participant) {
		this.participant = participant.sections.CD.data
		this.animateXYLine(this.participant, "cd")
	}

	updateNorms(norms) {
		this.norms = norms.CD.data
		this.animateNorms(this.norms, "cd")
	}

	drawLines(svg) {
		const useSD = (this.norms !== undefined)
		const norms = (this.norms === undefined) ? this.participant : this.norms
		this.createXYLine(this.participant, "cd")
		this.animateXYLine(this.participant, "cd")
		this.createNorms(norms, "cd", useSD)
		this.animateNorms(norms, "cd", useSD)
	}
}
