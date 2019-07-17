class ThresholdElectrotonus extends Chart {
	constructor(participant, norms) {
		super([0, 200], [-150, 100])
		this.participant = participant.sections.TE
		this.norms = (norms == null) ? undefined : norms.TE
		this.expectedKeys = Object.keys(this.participant) // set the initial keys
	}

	get name() { return "Threshold Electrotonus" }
	get xLabel() { return "Delay (ms)" }
	get yLabel() { return "Threshold Reduction (%)" }

	updateParticipant(participant) {
		this.participant = participant.sections.TE
		this.animateParticipant()
	}

	updateNorms(norms) {
		this.norms = norms.TE
		this.animateUpdatedNorms(this.norms)
	}

	drawLines(svg) {
		const isNull = (this.norms == null)
		const norms = isNull ? this.participant : this.norms
		Object.keys(this.participant).forEach(key => {
			this.createXYLine(this.participant[key].data, key)
			this.createNorms(norms[key].data, key, !isNull)
		})
		this.drawHorizontalLine(this.linesLayer, 0)
		this.animateParticipant()
		this.animateUpdatedNorms(norms, !isNull)
	}

	animateParticipant() {
		this.expectedKeys.forEach(key => {
			if (this.participant == null || this.participant[key] == null || this.participant[key].data == null) {
				this.animateXYLine(undefined, key)
			} else {
				this.animateXYLine(this.participant[key].data, key)
			}
		})
	}

	animateUpdatedNorms(norms, useSD = true) {
		Object.keys(norms).forEach(key => {})

		this.expectedKeys.forEach(key => {
			if (this.norms == null || this.norms[key] == null || this.norms[key].data == null) {
				this.animateNorms(undefined, key, useSD)
			} else {
				this.animateNorms(norms[key].data, key, useSD)
			}
		})
	}
}
