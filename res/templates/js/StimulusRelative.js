class StimulusRelative extends Chart {
	constructor(participant, norms) {
		super([0, 200], [0, 100])
		this.participant = this.calculateParticipant(participant.sections.SR.data.data)
		this.norms = (norms == null) ? undefined : norms.SRel.data
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
		if (participant.sections.SR == null || participant.sections.SR.data == null) {
			this.participant = this.calculateParticipant(undefined)
		} else {
			this.participant = this.calculateParticipant(participant.sections.SR.data.data)
		}
		this.animateXYLine(this.participant, "srel")
	}

	updateNorms(norms) {
		if (norms.SRel == null) {
			this.norms = undefined
		} else {
			this.norms = norms.SRel.data
		}
		this.animateNorms(this.norms, "srel")
	}

	updateScore(scores) {
		if (scores != null && scores.SRel != null) {
			this.score = scores.SRel.Overall
		}
	}

	drawLines(svg) {
		const isNull = (this.norms == null)
		const norms = isNull ? this.participant : this.norms
		this.createXYLine(this.participant, "srel")
		this.createNorms(norms, "srel", !isNull)
		this.animateXYLine(this.participant, "srel")
		this.animateNorms(norms, "srel", !isNull)
	}
}
