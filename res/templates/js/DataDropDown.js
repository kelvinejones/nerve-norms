class DataDropDown {
	// elementID is the HTML ID.
	// data is an object indexed by the drop-down's values
	// action is a function that receives the selection name and data as its arguments
	constructor(elementID, data, action, initialOptions) {
		this.elementID = elementID
		this.dt = data
		this.action = action

		this.updateData = (ev) => {
			this.val = ev.srcElement.value
			this.action(this.val, this.dt[this.val])
		}

		const dropDown = document.getElementById(this.elementID)
		dropDown.addEventListener("change", this.updateData)
		initialOptions.forEach(function(opt) {
			dropDown.options[dropDown.options.length] = new Option(opt)
		})
		this.val = dropDown.value
	}

	get selection() {
		return this.val
	}

	get data() {
		return this.dt[this.val]
	}
}
