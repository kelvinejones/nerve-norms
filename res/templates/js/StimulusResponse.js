class StimulusResponse extends Chart {
	constructor(participant, norms) {
		super([1, 20], [0.01, 20], Chart.scaleType.LOG, Chart.scaleType.LOG)
		this.participant = this.calculateParticipant(participant)
		this.norms = (norms === undefined) ? undefined : norms.SR.data
		this.sdFunc = Chart.logSD
		this.xSDIndex = 4
	}

	calculateParticipant(participant) {
		let peakResponse = 0;
		if (participant.sections.ExVars === undefined || participant.sections.ExVars.data == undefined) {
			return undefined
		}
		participant.sections.ExVars.data.forEach(function(exvar) {
			if (exvar[0] == 6) {
				peakResponse = exvar[1];
			}
		})
		if (peakResponse == 0) {
			return undefined
		}

		if (participant.sections.SR === undefined || participant.sections.SR.data == undefined || participant.sections.SR.data.data == undefined) {
			return undefined
		}
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
		const useSD = (this.norms !== undefined)
		const norms = (this.norms === undefined) ? this.participant : this.norms
		this.createXYLine(this.participant, "sr")
		this.createNorms(norms, "sr", useSD)
		this.animateXYLine(this.participant, "sr")
		this.animateNorms(norms, "sr", useSD)
	}
}
