class SendEmail {
    static get url() { return "https://us-central1-nervenorms-294404.cloudfunctions.net/" }
    

    // API call to send email
	static ContactEmail(data) {
		const query = SendEmail.url + "contact"
		fetch(query, { method: 'POST', body: data })
			.then(response => {
				if (!response.ok){
					alert("Error sending email, please try again or email directly at kelvin.jones@ualberta.ca")
					throw Error(response.statusText)
				}
				$('#btn-send-email').html("Success!")
				return response
			}).catch(e =>{
				alert("Network error. Check your connection or email directly at kelvin.jones@ualberta.ca")
			});

    }
}

function btnSendEmail(){

	var name = $('#contact_name').val()
	var email = $('#contact_email').val()
	var message = $('#contact_message').val()
	var cc = $('#email_copy').is(":checked")
	// field validation, probably change this to proper bootstrap in future
	var isValid = true
	if (name == ""){
		$('#contact_name').get(0).setCustomValidity("Empty Field")
		isValid = false
	} else{
		$('#contact_name').get(0).setCustomValidity("")
	}

	const regex = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
    if (!regex.test(String(email).toLowerCase())){
		$('#contact_email').get(0).setCustomValidity("Bad Email")
		isValid = false
	} else{
		$('#contact_email').get(0).setCustomValidity("")
	}

	if (message == ""){
		$('#contact_message').get(0).setCustomValidity("Empty Field")
		isValid = false
	} else{
		$('#contact_message').get(0).setCustomValidity("")
	}

	if (!isValid){
		return
	}
	var btn = $(this);
	btn.prop('disabled', true);
           setTimeout(function(){btn.prop('disabled', false); }, 3000);
	var package = {
		"Name": name,
		"Sender": email,
		"Subject": "Contact from " + name,
		"Message": message,
		"CarbonCopy": cc,
	}
	var data = JSON.stringify(package)

	SendEmail.ContactEmail(data)
}