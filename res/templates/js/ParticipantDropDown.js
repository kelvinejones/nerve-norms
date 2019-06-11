class ParticipantDropDown {
	// elementID is the HTML ID.
	// data is an object indexed by the drop-down's values
	// action is a function that receives the name name and data as its arguments
	constructor(elementID, data, filter, callback, initialOptions) {
		this.dt = data

		const updateData = (ev) => {
			this.val = ev.srcElement.value
			const participantData = this.dt[this.val]

			ExVars.clearScores()

			callback(participantData)

			ExVars.updateValues(participantData)
			filter.update(this.val)
		}

		const dropDown = document.getElementById(elementID)
		dropDown.addEventListener("change", updateData)
		initialOptions.forEach(function(opt) {
			dropDown.options[dropDown.options.length] = new Option(opt)
		})
		this.val = dropDown.value

		filter.update(this.val)
	}

	get name() {
		return this.val
	}

	get data() {
		return this.dt[this.val]
	}
}
