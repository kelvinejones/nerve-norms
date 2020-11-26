class DataManager {
	// data is an object indexed by the drop-down's values
	// The dataUsers is expected to provide a list of objects that implement 'updateParticipant', 'updateNorms', and 'updateScore'
	constructor(data, dataUsers) {
		this.dt = data
		this.dataUsers = dataUsers
		this.uploadCount = 0
		this.normCache = {}
		this.outlierCache = {}

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
			this._setFilter()
			ExVars.clearScores()
			this._updateParticipant()
			this._fetchUpdates()
		})
		this._updateDropDownOptions()


		this.uploadMemInput = document.getElementById('uploadMEM')
		uploadMemInput.onchange = e => {
			const file = e.target.files[0]

			const reader = new FileReader()
			reader.readAsText(file, 'UTF-8')

			reader.onload = readerEvent => {
				// Randomizing name client side
				var unAnonContent = readerEvent.target.result // this is the content!
				var randName = this._getRandomName(6)
				const content = unAnonContent.replace(/([\n\r].*Name:\s*)([^\n\r]*)/i, "$1"+randName );

				// Something in here is forcing the page to scroll to the top, didn't have time to fix
				Fetch.MEM(this.queryString, content, convertedMem => {
					if (convertedMem.error != null) {
						console.log("Conversion error", convertedMem.error)
						alert("The MEM could not be converted. Please email it to jbell1@ualberta.ca for troubleshooting.")
						return
					}

					this.uploadCount = this.uploadCount + 1
					const name = "Upload " + this.uploadCount + ": " + convertedMem.participant.header.name
					this.participants[this.participants.length] = new Participant(convertedMem.participant, name, false)
					this._updateDropDownOptions()

					this.dropDown.selectedIndex = this.dropDown.options.length - 1
					this.participantIndex = this.dropDown.selectedIndex

					this.outlierCache[this._cacheString(this.participantIndex)] = convertedMem.outlierScores
					this._updateParticipant()
					this._setFilter()
					$("#apply-filter").trigger("click")
					ExVars.updateScores(convertedMem.outlierScores)
				})
			}
		}
		this._fetchUpdates()
	}

	// Randomized alphanumeric name
	// First two will always be letters
	_getRandomName(length) {
		var result           = '';
		var characters       = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';
		var charactersLength = characters.length;
		for ( var i = 0; i < Math.min(2,length); i++ ) {
		   result += characters.charAt(Math.floor(Math.random() * charactersLength));
		}
		if (length > 2){
			result += "-" 
			for ( var i = 2; i < length; i++ ) {
				result += Math.floor(Math.random() * 10).toString();
			}
		}
		return result
	}

	_fetchUpdates() {
		const lastQuery = this.queryString
		this.queryString = Filter.queryString
		const normChanged = (lastQuery != this.queryString)
		if (normChanged) {
			this._fetchNorms()
		}

		const nameChanged = (this.participantIndex != this.dropDown.selectedIndex)
		if (normChanged || nameChanged) {
			this.participantIndex = this.dropDown.selectedIndex
			const participant = this.participants[this.participantIndex]

			ExVars.clearScores()

			this._fetchOutliers(participant)
		}
	}

	_fetchNorms() {
		const query = this.queryString
		const norms = this.normCache[query]
		if (norms != null) {
			Object.values(this.dataUsers()).forEach(pl => {
				pl.updateNorms(norms)
				this._updateToolTips(norms.ExVars.data)
			})
		} else {
			Fetch.Norms(query, (norms) => {
				this.normCache[query] = norms
				if (query == this.queryString) {
					// An update not has occurred since we requested this data, so update the display!
					Object.values(this.dataUsers()).forEach(pl => {
						pl.updateNorms(norms)
						this._updateToolTips(norms.ExVars.data)
					})
				}
			})
		}
	}

	// Update the tooltips to display normalized mean and SD
	_updateToolTips(data){
		for (var i = 0; i < data.length; i++){
			var eleID = "qtrac-excite-" + data[i][3].toString();
			var target = document.getElementById(eleID)
			if (target != null){
				if (target.title != undefined){
					let mean = data[i][0].toFixed(2);
					let standard_dev = data[i][1].toFixed(2);
					
					//template strings are ugly, bad indent is intended
					var newTitle = `Normalized Mean: ${mean};
Standard Deveation: ${standard_dev}`;
					$("#" + eleID).attr('title', newTitle)
						.attr('data-original-title', newTitle)
						.tooltip('update');
				}
			}
		}
	}

	_cacheString(ind) {
		if (ind == null) {
			ind = this.dropDown.selectedIndex
		}
		return this.queryString + "&id=" + ind
	}

	_fetchOutliers(participant) {
		const cacheString = this._cacheString() // Save this string because it's where we want to save the data
		const scores = this.outlierCache[cacheString]
		if (scores != null) {
			this._updateOutliers(scores)
		} else {
			const updateAction = (scores) => {
				this.outlierCache[cacheString] = scores
				if (cacheString == this._cacheString()) {
					// An update not has occurred since we requested this data, so update the display!
					this._updateOutliers(scores)
				}
			}

			if (participant.dataIsLocal) {
				Fetch.OutliersFromName(this.queryString, participant.name, updateAction)
			} else {
				Fetch.OutliersFromJSON(this.queryString, participant.data, updateAction)
			}
		}
	}

	_updateOutliers(scores) {
		ExVars.updateScores(scores)
		Object.values(this.dataUsers()).forEach(pl => {
			pl.updateScore(scores)
		})
	}

	_updateDropDownOptions() {
		const selection = this.dropDown.selectedIndex

		let index = 0
		this.participants.forEach(opt => {
			this.dropDown.options[index++] = new Option(opt.name)
		})

		if (selection >= 0) {
			this.dropDown.selectedIndex = selection
		}
	}

	_updateParticipant() {
		const data = this.participantData
		Object.values(this.dataUsers()).forEach(pl => {
			pl.updateParticipant(data)
		})

		ExVars.updateValues(data)
	}

	_setFilter(){
		//convert to hash table in long run if more options added, specifically countries
		if (document.getElementById("autoFilter").checked == true){
			const data = this.participantData
			switch(data.header.sex){
				case(0):
					$("#sex-option2").prop("checked", true).trigger("click");
					break;
				case(1):
					$("#sex-option3").prop("checked", true).trigger("click");
					break;
			}

			switch(true){
				case data.header.age < 30:
					$("#age-option2").prop("checked", true).trigger("click");
					break;
				case data.header.age <= 40:
					$("#age-option3").prop("checked", true).trigger("click");
					break;
				case data.header.age <= 50:
					$("#age-option4").prop("checked", true).trigger("click");
					break;
				case data.header.age <= 60:
					$("#age-option5").prop("checked", true).trigger("click");
					break;
				case data.header.age > 60:
					$("#age-option6").prop("checked", true).trigger("click");
					break;
			}

			switch(data.header.country){
				case ("CA"):
					$("#country-option2").prop("checked", true).trigger("click");
					break;
				case ("JP"):
					$("#country-option3").prop("checked", true).trigger("click");
					break;
				case ("PT"):
					$("#country-option4").prop("checked", true).trigger("click");
					break;
				default:
					$("#country-option1").prop("checked", true).trigger("click");
			}
		}
	}

	get norms() {
		return this.normCache[this.queryString]
	}

	get outliers() {
		return this.outlierCache[this._cacheString()]
	}

	get participantName() {
		return this.participants[this.dropDown.selectedIndex].name
	}

	get participantData() {
		return this.participants[this.dropDown.selectedIndex].data
	}
}
