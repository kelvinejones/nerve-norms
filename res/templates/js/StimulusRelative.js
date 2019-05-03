class StimulusRelative extends Chart {
	constructor(plots) {
		super([0, 200], [0, 100])
		this.data = this.calculateData(plots.sr.data)
		this.xName = 'x'
		this.yName = 'y'
	}

	calculateData(data) {
		const stimFor50PercentMax = data[24].valueX // Could also be extracted from excitability variables
		return data.map(function callback(d, i) {
			return {
				'x': d.valueX / stimFor50PercentMax * 100, // Normalize each element
				'y': (i + 1) * 2,
			}
		})
	}

	get name() { return "Relative Stimulus Response" }
	get xLabel() { return "Stimulus (% Mean Threshold)" }
	get yLabel() { return "Peak Response (% Max)" }

	drawLines(svg) {
		this.createXYLineWithMean(this.data, "srel")
		this.animateXYLineWithMean(this.data, "srel")
	}
}
