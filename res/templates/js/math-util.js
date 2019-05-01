function normativeRange(data, xName = 'delay', yName = 'mean', ySDName = 'SD') {
	return (Array.from(data)
			.map(function(d) { return { x: d[xName], y: d[yName] + 2 * (d[ySDName] || 0) } }))
		.concat(Array.from(data).reverse().map(function(d) { return { x: d[xName], y: d[yName] - 2 * (d[ySDName] || 0) } }))
}
