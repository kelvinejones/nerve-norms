class Chart {
	constructor(xRange, yRange, xscale = d3.scaleLinear(), yscale = d3.scaleLinear()) {
		this.xRange = xRange
		this.yRange = yRange

		this.margin = { top: 20, right: 20, bottom: 50, left: 50 };
		this.width = 600 - this.margin.left - this.margin.right;
		this.height = 300 - this.margin.top - this.margin.bottom;

		// Set the default scaling
		this.xscale = xscale.range([0, this.width]).domain(xRange);
		this.yscale = yscale.range([this.height, 0]).domain(yRange);

		// Set default config values
		this.xName = 'delay'
		this.yName = 'value'
		this.ySDName = 'SD'
		this.yMeanName = 'mean'

		this.yAnimStart = this.animationStartValue(this.yRange)
	}

	animationStartValue(range) {
		// The animation start value should be at zero unless that's outside of range
		if (range[0] < 0 && range[1] >= 0) {
			return 0
		} else {
			return range[0]
		}
	}

	get name() { throw new Error("A Chart must implement name()") }
	get xLabel() { throw new Error("A Chart must implement xLabel()") }
	get yLabel() { throw new Error("A Chart must implement yLabel()") }
	drawLines(svg) { throw new Error("A Chart must implement drawLines(svg)") }

	draw(svg, hideLabels) {
		// append the svg object to the body of the page
		// appends a 'group' element to 'svg'
		// moves the 'group' element to the top left margin
		svg = svg
			.append("g")
			.attr("transform",
				"translate(" + this.margin.left + "," + this.margin.top + ")");

		// Add layers for various elements
		this.ciLayer = svg.append("g")
		this.meanLayer = svg.append("g")
		this.linesLayer = svg.append("g")
		this.valueLayer = svg.append("g")
		this.circlesLayer = svg.append("g")

		// Add the X Axis
		var xelements = svg.append("g")
			.attr("transform", "translate(0," + this.height + ")")
			.call(d3.axisBottom(this.xscale).ticks(2)
				.tickFormat(d3.format("")));

		// Add the Y Axis
		var yelements = svg.append("g")
			.call(d3.axisLeft(this.yscale));

		if (!hideLabels) {
			this.labels(svg);
		} else {
			xelements.selectAll("text").remove();
			yelements.selectAll("text").remove();
		}

		this.drawLines(svg)
	}

	labels(svg) {
		// text label for the x axis
		svg.append("text")
			.attr("transform",
				"translate(" + (this.width / 2) + " ," +
				(this.height + this.margin.top + 20) + ")")
			.style("text-anchor", "middle")
			.text(this.xLabel);

		// text label for the y axis
		svg.append("text")
			.attr("transform", "rotate(-90)")
			.attr("y", 0 - this.margin.left)
			.attr("x", 0 - (this.height / 2))
			.attr("dy", "1em")
			.style("text-anchor", "middle")
			.text(this.yLabel);
	}

	normativeRange(data) {
		let xMean = this.xMeanName || this.xName
		let yMean = this.yMeanName || this.yName
		let ySD = this.ySDName
		let xSD = this.xSDName
		let first = data[0]
		let last = data[data.length - 1]
		return (Array.from(data)
				.map(function(d) { return { x: d[xMean] - 2 * (d[xSD] || 0), y: d[yMean] + 2 * (d[ySD] || 0) } }))
			.concat({ x: last[xMean] + 2 * (last[xSD] || 0), y: last[yMean] + 2 * (last[ySD] || 0) })
			.concat(Array.from(data).reverse().map(function(d) { return { x: d[xMean] + 2 * (d[xSD] || 0), y: d[yMean] - 2 * (d[ySD] || 0) } }))
			.concat({ x: first[xMean] - 2 * (first[xSD] || 0), y: first[yMean] - 2 * (first[ySD] || 0) })
	}

	dataAsXY(data, xName, yName) {
		return data.map(function(d) { return { x: d[xName], y: d[yName] } })
	}

	xZeroLine() {
		return d3.line()
			.x(d => this.xscale(d.x))
			.y(d => this.yscale(this.yAnimStart))
	}

	xyLine() {
		return d3.line()
			.x(d => this.xscale(d.x))
			.y(d => this.yscale(d.y));
	}

	animateCI(svg, ciNormRange) {
		svg.append("path")
			.data(ciNormRange)
			.attr("class", "confidenceinterval")
			.attr("d", this.xZeroLine())
			.transition()
			.delay(Chart.delayTime)
			.duration(Chart.transitionTime)
			.attr("d", this.xyLine());
	}

	animateLine(svg, xyLine, className) {
		svg.append("path")
			.data(xyLine)
			.attr("class", className)
			.attr("d", this.xZeroLine())
			.transition()
			.delay(Chart.delayTime)
			.duration(Chart.transitionTime)
			.attr("d", this.xyLine());
	}

	drawHorizontalLine(svg, yVal) {
		svg.append("path")
			.data([
				[{ x: this.xRange[0], y: yVal }, { x: this.xRange[1], y: yVal }]
			])
			.attr("class", "meanline")
			.attr("d", this.xyLine());
	}

	drawVerticalLine(svg, xVal) {
		svg.append("path")
			.data([
				[{ y: this.yRange[0], x: xVal }, { y: this.yRange[1], x: xVal }]
			])
			.attr("class", "meanline")
			.attr("d", this.xyLine());
	}

	animateCircles(svg, circleLocations) {
		// Add circles into a separate SVG group
		svg = svg.append("g")
		const self = this

		const circles = svg.selectAll("circle")
			.data(circleLocations)
			.enter()
			.append("circle");
		circles.attr("cx", d => self.xscale(d[this.xName]))
			.attr("cy", self.yscale(this.yAnimStart))
			.attr("r", d => d.wasImputed ? 3 : 5)
			.style("fill", d => d.wasImputed ? "red" : "black");
		circles
			.transition()
			.delay(Chart.delayTime)
			.duration(Chart.transitionTime)
			.attr("cy", d => self.yscale(d[this.yName]))
	}

	animateXYLineWithMean(lineData) {
		this.animateCI(this.ciLayer, [this.normativeRange(lineData)])
		this.animateLine(this.meanLayer, [this.dataAsXY(lineData, this.xMeanName || this.xName, this.yMeanName)], "meanline")
		this.animateLine(this.valueLayer, [this.dataAsXY(lineData, this.xName, this.yName)], "line")
		this.animateCircles(this.circlesLayer, lineData)
	}
}

// Set some constants for the class
Object.defineProperty(Chart, 'delayTime', {
	value: 750,
	enumerable: true,
})
Object.defineProperty(Chart, 'transitionTime', {
	value: 2500,
	enumerable: true,
})
