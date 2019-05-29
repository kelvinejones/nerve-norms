class StimulusRelative extends Chart {
	constructor(participant, norms) {
		super([0, 200], [0, 100])
		this.participant = participant.sections.SR.data
		this.norms = this.calculateNorms(norms.sections.SR.data)
		this.xName = 'valueX'
		this.yName = 'valueY'
		this.xSDName = 'SD'
		this.ySDName = undefined
		this.xMeanName = 'mean'
		this.yMeanName = 'y'
	}

	calculateNorms(data) {
		const meanStimFor50PercentMax = data[24].meanX
		return data.map((d, i) => {
			return {
				'y': (i + 1) * 2,
				'SD': d.SDX / d.meanX * 100,
				'mean': d.meanX / meanStimFor50PercentMax * 100,
			}
		})
	}

	get name() { return "Relative Stimulus Response" }
	get xLabel() { return "Stimulus (% Mean Threshold)" }
	get yLabel() { return "Peak Response (% Max)" }

	updateParticipant(participant) {
		this.participant = participant.sections.SR.data
		this.animateXYLine(this.participant, "srel")
	}

	updateNorms(norms) {
		this.norms = this.calculateNorms(norms.sections.SR.data)
		this.animateNorms(this.norms, "srel")
	}

	drawLines(svg) {
		this.createXYLine(this.participant, "srel")
		this.createNorms(this.norms, "srel")
		this.animateXYLine(this.participant, "srel")
		this.animateNorms(this.norms, "srel")
	}
}
