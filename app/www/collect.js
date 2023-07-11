

function submit_form() {
	//collect the data
	const groupData = [];
	const groupMemberNames = document.querySelectorAll(".feedback-data__row-name");
	const comments = document.querySelectorAll(".feedback-comments__member-comments");
	const allScores = document.querySelectorAll(".feedback-data__cell");
	//for each group member, collect and validate specific data
	for (var i = 0; i < 6; i++) {
		let memComment = "";
		if (i == 0) {
			memComment = comments[0].value;
			if (memComment == "") {
				// error - TODO: implement error in data validation
			}
		}
		memComment = memComment.concat("\n", comments[i+1].value);
		if (memComment == "") {
			// error - TODO: implement error in data validation
		}
		const scores = [];
		for (var k = 0; k < 5; k++) {
			if (allScores[i*5+k].value == '') {
				// error - TODO: implement error in data validation
			}
			scores.push(parseFloat(allScores[i*5+k].value));
		}
		const data = {
			name: groupMemberNames[i].innerText,
			scores,
			comment: memComment
		};
		groupData.push(data);
	}
	//write data to object
	const studentInput = {
		studentName: "John",
		studentGroup: "Beatles",
		groupSize: 6,
		groupsData: groupData
	};
	//json serialize it
	let data = JSON.stringify(studentInput);
	//send in a post reqest to server '/api/submit'?

	console.log(data)
	//report successful submission
	fetch("/api/submit", {
		method: "POST",
		body: data,
		headers: {
			"Content-type": "application/json; charset=UTF-8"
		}
	}).then(res => {
		if (res.ok) {
			document.getElementById("successful_submit").innerHTML = "Form submitted successfully.";
		} else {
			document.getElementById("successful_submit").innerHTML = "Form submission error.";
		}
	}).catch(err => {
		document.getElementById("successful_submit").innerHTML = "Form submission error.";
	});
}
