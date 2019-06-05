class ThresholdElectrotonus extends Chart {
	constructor(participant, norms) {
		super([0, 200], [-150, 100])
		this.participant = participant.sections.TE
		this.norms = (norms === undefined) ? undefined : norms.TE
	}

	get name() { return "Threshold Electrotonus" }
	get xLabel() { return "Threshold Reduction (%)" }
	get yLabel() { return "Delay (ms)" }

	updateParticipant(participant) {
		this.participant = participant.sections.TE
		this.animateParticipant()
	}

	updateNorms(norms) {
		this.norms = norms.TE
		this.animateUpdatedNorms(this.norms)
	}

	drawLines(svg) {
		const useSD = (this.norms !== undefined)
		const norms = (this.norms === undefined) ? this.participant : this.norms
		Object.keys(this.participant).forEach(key => {
			this.createXYLine(this.participant[key].data, key)
			this.createNorms(norms[key].data, key, useSD)
		})
		this.drawHorizontalLine(this.linesLayer, 0)
		this.animateParticipant()
		this.animateUpdatedNorms(norms, useSD)
	}

	animateParticipant() {
		Object.keys(this.participant).forEach(key => {
			this.animateXYLine(this.participant[key].data, key)
		})
	}

	animateUpdatedNorms(norms, useSD = true) {
		Object.keys(norms).forEach(key => {
			this.animateNorms(norms[key].data, key, useSD)
		})
	}
}
