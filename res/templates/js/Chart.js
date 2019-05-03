class Chart {
	constructor(xRange, yRange, xScaleType = Chart.scaleType.LINEAR, yScaleType = Chart.scaleType.LINEAR) {
		this.xRange = xRange
		this.yRange = yRange

		this.margin = { top: 20, right: 20, bottom: 50, left: 50 };
		this.width = 600 - this.margin.left - this.margin.right;
		this.height = 300 - this.margin.top - this.margin.bottom;

		// Set the default scaling
		this.xscale = this.makeScale(xScaleType).range([0, this.width]).domain(xRange);
		this.yscale = this.makeScale(yScaleType).range([this.height, 0]).domain(yRange);

		// Set default config values
		this.xName = 'delay'
		this.yName = 'value'
		this.ySDName = 'SD'
		this.yMeanName = 'mean'

		this.yAnimStart = this.animationStartValue(this.yRange)

		this.delayTime = Chart.slowDelay;
		this.transitionTime = Chart.slowTransition;

		this.group = {}

		this.limFunc = this.limAtLoc // For SD: (dpt, loc) => { return this.sdAtLoc(dpt, loc, 1) }
	}

	makeScale(name) {
		let scale
		switch (name) {
			case "LINEAR":
				scale = d3.scaleLinear()
				scale.scaleType = Chart.scaleType.LINEAR
				break
			case "LOG":
				scale = d3.scaleLog()
				scale.scaleType = Chart.scaleType.LOG
				break
		}
		return scale
	}

	animationStartValue(range) {
		// The animation start value should be at zero unless that's outside of range
		if (range[0] < 0 && range[1] >= 0) {
			return 0
		} else {
			return range[0]
		}
	}

	updatePlots(plots) { throw new Error("A Chart must implement updatePlots(plots)") }
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
			svg.attr("transform", "scale(0.9) translate(110, 0)")
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

	// sdAtLoc calculates limits based on the standard deviations
	sdAtLoc(dpt, loc, numSD = 2) {
		const xmn = this.xMeanName || this.xName
		const ymn = this.yMeanName || this.yName
		const xsd = this.xSDName
		const ysd = this.ySDName

		function stPoint(xSign, ySign) {
			return {
				x: dpt[xmn] + xSign * numSD * (dpt[xsd] || 0),
				y: dpt[ymn] + ySign * numSD * (dpt[ysd] || 0),
			}
		}

		switch (loc) {
			case Chart.limLoc.UPPER_LEFT:
				return stPoint(-1, 1)
			case Chart.limLoc.UPPER_RIGHT:
				return stPoint(1, 1)
			case Chart.limLoc.LOWER_LEFT:
				return stPoint(-1, -1)
			case Chart.limLoc.LOWER_RIGHT:
				return stPoint(1, -1)
		}
	}

	normativeLimits(data) {
		if (this.limFunc === undefined) {
			return []
		}

		// this.limFunc can be changed if a different calculation is more appropriate.
		return this.scaleArrayWithinRange((Array.from(data)
				.map(d => { return this.limFunc(d, Chart.limLoc.UPPER_LEFT) }))
			.concat(this.limFunc(data[data.length - 1], Chart.limLoc.UPPER_RIGHT))
			.concat(Array.from(data).reverse().map(d => { return this.limFunc(d, Chart.limLoc.LOWER_RIGHT) }))
			.concat(this.limFunc(data[0], Chart.limLoc.LOWER_LEFT)))
	}

	// limAtLoc extracts the calculated limits from the dataset, which describes the range in which a healthy measure is expected
	limAtLoc(dpt, loc) {
		const xmn = this.xMeanName || this.xName
		const ymn = this.yMeanName || this.yName
		switch (loc) {
			case Chart.limLoc.UPPER_LEFT:
				return { x: dpt['leftLimit'] || dpt[xmn], y: dpt['upperLimit'] || dpt[ymn] }
				break
			case Chart.limLoc.UPPER_RIGHT:
				return { x: dpt['rightLimit'] || dpt[xmn], y: dpt['upperLimit'] || dpt[ymn] }
				break
			case Chart.limLoc.LOWER_LEFT:
				return { x: dpt['leftLimit'] || dpt[xmn], y: dpt['lowerLimit'] || dpt[ymn] }
				break
			case Chart.limLoc.LOWER_RIGHT:
				return { x: dpt['rightLimit'] || dpt[xmn], y: dpt['lowerLimit'] || dpt[ymn] }
				break
		}
	}

	scaleArrayWithinRange(ar) {
		// With a log scale, values can't be plotted at or below zero.
		if (this.xscale.scaleType == Chart.scaleType.LOG) {
			ar = this.raiseZeroValues(ar, 'x', this.xRange[0])
		}
		if (this.yscale.scaleType == Chart.scaleType.LOG) {
			ar = this.raiseZeroValues(ar, 'y', this.yRange[0])
		}
		return ar
	}

	raiseZeroValues(ar, axis, min) {
		return ar.map(function(d) {
			if (d[axis] < min) {
				d[axis] = min
			}
			return d
		})
	}

	dataAsXY(data, xName, yName) {
		if (data[0][xName] === undefined || data[0][yName] === undefined) {
			return []
		}
		return data.map(function(d) { return { x: d[xName], y: d[yName] } })
	}

	xZeroPath() {
		return d3.line()
			.x(d => this.xscale(d.x))
			.y(d => this.yscale(this.yAnimStart))
	}

	xyPath() {
		return d3.line()
			.x(d => this.xscale(d.x))
			.y(d => this.yscale(d.y));
	}

	drawHorizontalLine(svg, yVal) {
		svg.append("path")
			.data([
				[{ x: this.xRange[0], y: yVal }, { x: this.xRange[1], y: yVal }]
			])
			.attr("class", "meanline")
			.attr("d", this.xyPath());
	}

	drawVerticalLine(svg, xVal) {
		svg.append("path")
			.data([
				[{ y: this.yRange[0], x: xVal }, { y: this.yRange[1], x: xVal }]
			])
			.attr("class", "meanline")
			.attr("d", this.xyPath());
	}

	createGroup(svg, typeString, name) {
		svg = svg.append("g")
		this.group[typeString + "-" + name] = svg
		return svg
	}

	animateGroup(typeString, newData, name) {
		return this.group[typeString + "-" + name].selectAll(typeString)
			.data(newData)
			.transition()
			.delay(this.delayTime)
			.duration(this.transitionTime)
	}

	createPath(svg, path, groupName, className) {
		this.createGroup(svg, "path", groupName + "-" + className)
			.append("path")
			.data(path)
			.attr("class", className)
			.attr("d", this.xZeroPath())
	}

	animatePath(path, groupName, className) {
		this.animateGroup("path", path, groupName + "-" + className)
			.attr("d", this.xyPath());
	}

	createCircles(svg, circleLocations, name) {
		// create circle locations at init for each name, with right amount and position of circles
		// Add circles into a separate SVG group
		this.createGroup(svg, "circle", name)
			.selectAll("circle")
			.data(circleLocations)
			.enter()
			.append("circle")
			.attr("cx", d => this.xscale(d[this.xName]))
			.attr("cy", this.yscale(this.yAnimStart))
			.attr("r", d => 3)
			.style("fill", d => "black");
	}

	animateCircles(circleLocations, name) {
		this.animateGroup("circle", circleLocations, name)
			.attr("r", d => d.wasImputed ? 3 : 5)
			.style("fill", d => d.wasImputed ? "red" : "black")
			.attr("cy", d => this.yscale(d[this.yName]))
			.attr("cx", d => this.xscale(d[this.xName]))
	}

	createXYLineWithMean(lineData, name) {
		this.createPath(this.ciLayer, [this.normativeLimits(lineData)], name, "confidenceinterval")
		this.createPath(this.meanLayer, [this.dataAsXY(lineData, this.xMeanName || this.xName, this.yMeanName)], name, "meanline")
		this.createPath(this.valueLayer, [this.dataAsXY(lineData, this.xName, this.yName)], name, "line")
		this.createCircles(this.circlesLayer, lineData, name)
	}

	animateXYLineWithMean(lineData, name) {
		this.animatePath([this.normativeLimits(lineData)], name, "confidenceinterval")
		this.animatePath([this.dataAsXY(lineData, this.xMeanName || this.xName, this.yMeanName)], name, "meanline")
		this.animatePath([this.dataAsXY(lineData, this.xName, this.yName)], name, "line")
		this.animateCircles(lineData, name)
	}

	setDelayTime(dt) {
		// Return 'this' for chaining
		this.delayTime = dt
		return this
	}

	setTransitionTime(tt) {
		// Return 'this' for chaining
		this.transitionTime = tt
		return this
	}
}

// Set some constants for the class
Object.defineProperty(Chart, 'slowDelay', {
	value: 750,
	enumerable: true,
})
Object.defineProperty(Chart, 'fastDelay', {
	value: 0,
	enumerable: true,
})
Object.defineProperty(Chart, 'slowTransition', {
	value: 2500,
	enumerable: true,
})
Object.defineProperty(Chart, 'fastTransition', {
	value: 1000,
	enumerable: true,
})
Object.defineProperty(Chart, 'scaleType', {
	value: {
		LINEAR: "LINEAR",
		LOG: "LOG",
	},
	enumerable: true,
})
Object.defineProperty(Chart, 'limLoc', {
	value: {
		UPPER_LEFT: "UPPER_LEFT",
		UPPER_RIGHT: "UPPER_RIGHT",
		LOWER_LEFT: "LOWER_LEFT",
		LOWER_RIGHT: "LOWER_RIGHT",
	},
	enumerable: true,
})
