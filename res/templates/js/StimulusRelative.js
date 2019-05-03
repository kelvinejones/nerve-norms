class StimulusRelative extends Chart {
	constructor(plots) {
		super([0, 200], [0, 100])
		this.data = this.calculateData(plots.sr.data)
		this.xName = 'x'
		this.yName = 'y'
		this.xSDName = 'SD'
		this.ySDName = undefined
		this.xMeanName = 'mean'
		this.yMeanName = undefined
		this.limFunc = this.sdAtLoc // Only sd lines are supported
	}

	calculateData(data) {
		const stimFor50PercentMax = data[24].valueX // Could also be extracted from excitability variables
		const meanStimFor50PercentMax = data[24].meanX

		function norm(a, b) { return Math.sqrt(a * a + b * b); }

		return data.map((d, i) => {
			return {
				'y': (i + 1) * 2,
				// Normalize each element
				'x': d.valueX / stimFor50PercentMax * 100,
				'SD': norm(d.SDX, d.SDY) / d.meanX * 100,
				'mean': d.meanX / meanStimFor50PercentMax * 100,
			}
		})
	}

	get name() { return "Relative Stimulus Response" }
	get xLabel() { return "Stimulus (% Mean Threshold)" }
	get yLabel() { return "Peak Response (% Max)" }

	updatePlots(plots) {
		this.data = this.calculateData(plots.sr.data)
		this.animateXYLineWithMean(this.data, "srel")
	}

	drawLines(svg) {
		this.createXYLineWithMean(this.data, "srel")
		this.animateXYLineWithMean(this.data, "srel")
	}
}
