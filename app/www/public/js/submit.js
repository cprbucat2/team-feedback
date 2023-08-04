/**
 * @file Handles submissions from index.html.
 * @author Aidan Hoover, Aiden Woodruff
 * @copyright 2023 Aidan Hoover and Aiden Woodruff
 * @license BSD-3-Clause
 */

/**
 * Collect scores from a .feedback-data__score-table.
 * @param {HTMLTableElement} table The table to collect data from.
 * @requires table Is a well-formed .feedback-data__score-table.
 * @returns {{name: string, scores: number[]}[]} A list of member names and score lists.
 */
function collect_scores(table) {
	const data = [];
	for (const row of table.rows) {
		if (!row.classList.contains("feedback-data__categories") &&
		!row.classList.contains("feedback-data__colavg")) {
			let name;
			const datarow = [];
			for (const cell of row.cells) {
				if (cell.classList.contains("feedback-data__row-name")) {
					name = cell.innerText;
				} else if (!cell.classList.contains("feedback-data__row-avg")) {
					datarow.push(parseFloat(cell.firstChild.value));
				}
			}
			data.push({name, scores: datarow});
		}
	}
	return data;
}

/**
 * Submit form and print status beneath button or reject.
 * @listens MouseEvent
 */
function submit_form() {
	// Collect feedback score table data.
	const entries = collect_scores(document.querySelector(".feedback-data__score-table"));
	/** @type {string} */
	const improvement = document.querySelector("#self_improvement").value;

	const comments = [];
	document.querySelectorAll(".feedback-comments__member-comments").forEach(e => {
		if (e.id !== "self_improvement") {
			comments.push(e.value);
		}
	});
	if (entries.length !== comments.length) {
		console.error("Score table and comment rows do not match.");
		/** @todo Make this a fatal error once pages are generated. */
	}
	for (let i = 0; i < entries.length; ++i) {
		entries[i].comment = comments[i];
	}

	// Create Submission object.
	const membersubmission = {
		author: entries[0].name,
		entries,
		improvement
	};

	// Submit and report status.
	fetch("/api/submit", {
		method: "POST",
		body: JSON.stringify(membersubmission),
		headers: {
			"Content-type": "application/json; charset=UTF-8"
		}
	}).then(res => {
		if (res.ok) {
			document.getElementById("successful_submit").innerText = "Form submitted successfully.";
		} else {
			document.getElementById("successful_submit").innerText = "Form submission error.";
		}
	}).catch(() => {
		document.getElementById("successful_submit").innerText = "Form submission error.";
	});
}

window.addEventListener("load", () => {
	document.querySelector("#submit").addEventListener("click", submit_form);
});
