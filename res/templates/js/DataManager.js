class DataManager {
	// data is an object indexed by the drop-down's values
	// The dataUsers is expected to provide a list of objects that implement 'updateParticipant' and 'updateNorms'
	constructor(data, dataUsers) {
		this.dt = data
		this.dataUsers = dataUsers
		this.uploadCount = 0

		this.participants = [
			Participant.load("CA-WI20S", data),
			Participant.load("CA-AL27H", data),
			Participant.load("JP-20-1", data),
			Participant.load("JP-70-1", data),
			Participant.load("PO-00d97e84", data),
			Participant.load("PO-017182a5", data),
			Participant.load("CA Mean", data),
			Participant.load("JP Mean", data),
			Participant.load("PO Mean", data),
			Participant.load("Rat Fast Axon", data),
			Participant.load("Rat Slow Axon", data),
			Participant.load("Rat on Drugs", data),
		]

		Filter.setCallback(() => this._fetchUpdates())

		this.dropDown = document.getElementById("select-participant-dropdown")
		this.dropDown.addEventListener("change", (ev) => {
			ExVars.clearScores()
			if (this.dropDown.selectedIndex >= this.participants.length) {
				this._uploadMEM()
			} else {
				this._updateParticipant()
				this._fetchUpdates()
			}
		})
		this._updateDropDownOptions()

		this._fetchUpdates()
	}

	_fetchUpdates() {
		const lastQuery = this.queryString
		this.queryString = Filter.queryString
		const normChanged = (lastQuery != this.queryString)
		if (normChanged) {
			Fetch.Norms(this.queryString, norms => {
				this.normData = norms
				Object.values(this.dataUsers()).forEach(pl => {
					pl.updateNorms(norms)
				})
			})
		}

		const nameChanged = (this.participantIndex != this.dropDown.selectedIndex)
		if (normChanged || nameChanged) {
			this.participantIndex = this.dropDown.selectedIndex
			const participant = this.participants[this.participantIndex]

			ExVars.clearScores()

			if (participant.dataIsLocal) {
				Fetch.OutliersFromName(this.queryString, participant.name, ExVars.updateScores)
			} else {
				Fetch.OutliersFromJSON(this.queryString, participant.data, ExVars.updateScores)
			}
		}
	}

	static get uploadOption() { return "Upload MEM..." }

	_updateDropDownOptions() {
		const selection = this.dropDown.selectedIndex

		let index = 0
		this.participants.forEach(opt => {
			this.dropDown.options[index++] = new Option(opt.name)
		})
		this.dropDown.options[index] = new Option(DataManager.uploadOption)

		if (selection >= 0) {
			this.dropDown.selectedIndex = selection
		}
	}

	_uploadMEM() {
		// This code is modified from https://stackoverflow.com/a/40971885
		var input = document.createElement('input')
		input.type = 'file'

		input.onchange = e => {
			var file = e.target.files[0]

			var reader = new FileReader()
			reader.readAsText(file, 'UTF-8')

			reader.onload = readerEvent => {
				var content = readerEvent.target.result // this is the content!
				Fetch.MEM(this.queryString, content, convertedMem => {
					if (convertedMem.error != null) {
						console.log("Conversion error", convertedMem.error)
						alert("The MEM could not be converted. Please email it to jbell1@ualberta.ca for troubleshooting.")
						this.dropDown.selectedIndex = this.participantIndex
						this._fetchUpdates()
						return
					}

					this.uploadCount = this.uploadCount + 1
					const name = "Upload " + this.uploadCount + ": " + convertedMem.participant.header.name
					this.participants[this.participants.length] = new Participant(convertedMem.participant, name, false)
					this._updateDropDownOptions()

					this._updateParticipant()

					ExVars.updateScores(convertedMem.outlierScores)
				})
			}
		}

		input.click()
	}

	_updateParticipant() {
		const data = this.participantData
		Object.values(this.dataUsers()).forEach(pl => {
			pl.updateParticipant(data)
		})

		ExVars.updateValues(data)
	}

	get norms() {
		return this.normData
	}

	get participantName() {
		return this.participants[this.dropDown.selectedIndex].name
	}

	get participantData() {
		return this.participants[this.dropDown.selectedIndex].data
	}
}
