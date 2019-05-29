class ThresholdElectrotonus extends Chart {
	constructor(participant, norms) {
		super([0, 200], [-150, 100])
		this.participant = participant.sections.TE.data
		this.norms = norms.sections.TE.data
		this.xSDName = undefined // No standard deviation in the x direction.
		this.ySDName = 'SD'
		this.yMeanName = 'mean'
		this.xMeanName = 'delay'
	}

	get name() { return "Threshold Electrotonus" }
	get xLabel() { return "Threshold Reduction (%)" }
	get yLabel() { return "Delay (ms)" }

	updateParticipant(participant) {
		this.participant = participant.sections.TE.data
		this.animateParticipant()
	}

	updateNorms(norms) {
		this.norms = norms.sections.TE.data
		this.animateUpdatedNorms()
	}

	drawLines(svg) {
		Object.keys(this.participant).forEach(key => {
			this.createXYLine(this.participant[key], key)
			this.createNorms(this.norms[key], key)
		})
		this.drawHorizontalLine(this.linesLayer, 0)
		this.animateParticipant()
		this.animateUpdatedNorms()
	}

	animateParticipant() {
		Object.keys(this.participant).forEach(key => {
			this.animateXYLine(this.participant[key], key)
		})
	}

	animateUpdatedNorms() {
		Object.keys(this.norms).forEach(key => {
			this.animateNorms(this.norms[key], key)
		})
	}
}
