/**
 * @file Handle UI elements on admin pages.
 * @author Aiden Woodruff
 * @copyright 2023 Aidan Hoover and Aiden Woodruff
 * @license BSD-3-Clause
 */

jQuery(function ($) {
	/**
	 * Update remove button and head checkbox.
	 */
	function recountCheckboxes() {
		const boxes = $(".admin-list__checkbox")
			.not(".admin-list__entry--template *");
		const count = boxes.map((i, e) => e.checked).get()
			.reduce((x, sum) => x + sum);
		if (count === 0) {
			$(".admin-remove-btn").prop("disabled", true);
			$(".admin-list__head-checkbox").prop("checked", false);
			$(".admin-list__head-checkbox").prop("indeterminate", false);
		} else {
			$(".admin-remove-btn").prop("disabled", false);
			$(".admin-list__head-checkbox").prop("checked", count === boxes.length);
			$(".admin-list__head-checkbox")
				.prop("indeterminate", count !== boxes.length);
		}
		$(".admin-remove-btn").prop("value", `Remove (${count} selected)`);
	}

	/**
	 * Listen to checkbox clicks and update remove button and head checkbox.
	 * @listens MouseEvent
	 */
	$(".admin-list__checkbox").on("click", recountCheckboxes);

	/**
	 * Update checkboxes based off the header checkbox.
	 * @listens MouseEvent
	 */
	$(".admin-list__head-checkbox").on("click", event => {
		if (!event.target.checked) {
			$(".admin-list__checkbox").prop("checked", false);
			$(".admin-remove-btn").prop("disabled", true);
			$(".admin-remove-btn").prop("value", `Remove (0 selected)`);
		} else {
			const boxes = $(".admin-list__checkbox")
				.not(".admin-list__entry--template *");
			boxes.prop("checked", true);
			$(".admin-remove-btn").prop("disabled", false);
			$(".admin-remove-btn").prop("value", `Remove (${boxes.length} selected)`);
		}
	});

	/**
	 * Control the add team button. On click it makes the add-info row visible.
	 * @listens MouseEvent
	 */
	$(".team-add-btn").on("click", event => {
		event.target.disabled = true;
		$(".team-list__add-info").removeClass("team-list__add-info--hidden");
	});

	function hide_team_list__add_info () {
		$("#team-add-name").val("");
		$(".team-list__add-status").text("");
		$(".team-list__add-info").addClass("team-list__add-info--hidden");
		$(".team-add-btn").prop("disabled", false);
	}

	/**
	 * Control the confirm button when adding a team. Uploads new team name to the
	 * server and updates the team list on success.
	 * @listens MouseEvent
	 */
	$(".team-list__add-confirm").on("click", () => {
		const name = $("#team-add-name").val();

		// Do not upload empty names.
		if (typeof name === "undefined" || name === "") {
			$(".team-list__add-status").addClass("team-list__add-status--error");
			$(".team-list__add-status").text("Cannot add team with empty name");
			return;
		}

		fetch("/api/admin/team/add", {
			method: "POST",
			body: JSON.stringify({name}),
			headers: {
				"Content-type": "application/json; charset=UTF-8"
			}
		}).then(res => {
			if (res.ok && res.status == 201) {
				return res.json();
			} else throw new Error("Not ok.");
		}).then(team => {
			const row = $(".team-list__entry.admin-list__entry--template")
				.clone(true);
			row.removeClass("admin-list__entry--template");
			row.data("id", team.id);
			row.find(".team-list__name").text(team.name);
			row.find(".team-list__submissions")
				.attr("href", "/admin/submissions/team/" + team.id);
			$(".team-list tbody").append(row);
			recountCheckboxes();
			hide_team_list__add_info();
		}).catch(() => {
			$(".team-list__add-status").addClass("team-list__add-status--error");
			$(".team-list__add-status").text("Failed to add team.");
		});
	});

	/**
	 * Control the cancel button when adding a team. On click it clears the team
	 * name and response status, hides the add-info row, and enables the add team
	 * button.
	 * @listens MouseEvent
	 */
	$(".team-list__add-cancel").on("click", hide_team_list__add_info);

	/**
	 * Execute a DELETE query to remove teams.
	 * @param {string[]} teamIDs An array of team ids to remove.
	 * @param {function(): void} success Executes on success after removing rows.
	 * @param {function(): void} failure Executes on failure after the error
	 * dialog is closed.
	 */
	function removeTeams(teamIDs, success, failure) {
		fetch("/api/admin/team/remove", {
			method: "DELETE",
			body: JSON.stringify(teamIDs),
			headers: {
				"Content-type": "application/json; charset=UTF-8"
			}
		}).then(res => {
			if (res.ok && res.status == 200) {
				$(".team-list__entry").filter((_, e) =>
					typeof $(e).data("id") !== "undefined" &&
					teamIDs.includes($(e).data("id").toString())
				).remove();
				if (typeof success === "function") {
					success();
				}
			}
			return res.json();
		}).then(json => {
			if (typeof json.message === "string" && json.message === "success")
				return;
			let errorMsg = "Failed to remove team";
			if (typeof json.id === "number") {
				const teamName = $(".team-list__entry").filter(
					(i, e) => $(e).data("id") &&
					$(e).data("id").toString() === json.id.toString()
				).find(".team-list__name").text();
				errorMsg += " " + (teamName === "" ? json.id : teamName);
			}
			if (typeof json.message === "string") {
				errorMsg += ": " + json.message;
			} else {
				errorMsg += ".";
			}
			$("#server-error-dialog__msg").text(errorMsg);
			$("#server-error-dialog")[0].showModal();
			$("#server-error-dialog").one("close", () => {
				if (typeof failure === "function") {
					failure();
				}
			});
		});
	}

	$(".team-list__remove").on("click", event => {
		const teamId = $(event.target).parents(".team-list__entry").data("id");
		if (typeof teamId === "undefined") {
			$("#server-error-dialog__msg").text("Failed to remove team.");
			$("#server-error-dialog")[0].showModal();
			return;
		}

		event.target.disabled = true;
		removeTeams([teamId.toString()], recountCheckboxes, () => {
			event.target.disabled = false;
		});
	});

	$(".teams-remove-btn").on("click", () => {
		const ids = $(".admin-list__checkbox:checked")
			.parents(".admin-list__entry")
			.map((_, e) => $(e).data("id").toString()).get();
		removeTeams(ids, () => {
			$(".teams-remove-btn").prop("value", "Remove (0 selected)");
			$(".teams-remove-btn").prop("disabled", true);
			$(".admin-list__head-checkbox").prop("checked", false);
			$(".admin-list__head-checkbox").prop("indeterminate", false);
		});
	});
});
