

function submit_form() {
	//collect the data
	const groupData = [];
	const groupMemberNames = document.querySelectorAll(".feedback-data__row-name");
	const comments = document.querySelectorAll(".feedback-comments__member-comments");
	const allScores = document.querySelectorAll(".feedback-data__cell");
	//for each group member, collect and validate specific data
	for (var i = 0; i < 6; i++) {
		var memComment = "";
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
			scores.push(allScores[i*5+k].value);
		}
		const data = [groupMemberNames[i].innerText, scores, memComment];
		groupData.push(data);
	}
	//write data to object
	const studentInput = {
		studentName: "John",
		studentGroup: "Beatles",
		groupSize: 6,
		groupsData: groupData
	};
	console.log(studentInput);
	//json serialize it
	var data = JSON.stringify(studentInput);
	//send in a post reqest to server '/api/submit'?

	//report successful submission
	document.getElementById("successful_submit").innerHTML = "Form submitted successfully.";
}
