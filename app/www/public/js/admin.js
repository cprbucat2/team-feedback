/**
 * @file Handle UI elements on admin pages.
 * @author Aiden Woodruff
 * @copyright 2023 Aidan Hoover and Aiden Woodruff
 * @license BSD-3-Clause
 */

jQuery(function ($) {
	/**
	 * Listen to checkbox clicks and update remove button and head checkbox.
	 * @listens MouseEvent
	 */
	$(".admin-list__checkbox").on("click", () => {
		const boxes = $(".admin-list__checkbox");
		const count = boxes.map((i, e) => e.checked).get()
			.reduce((x, sum) => x + sum);
		if (count === 0) {
			$(".remove-users-btn").prop("disabled", true);
			$(".user-list__head-checkbox").prop("checked", false);
			$(".user-list__head-checkbox").prop("indeterminate", false);
		} else {
			$(".remove-users-btn").prop("disabled", false);
			$(".user-list__head-checkbox").prop("checked", count === boxes.length);
			$(".user-list__head-checkbox")
				.prop("indeterminate", count !== boxes.length);
		}
		$(".admin-remove-btn").prop("value", `Remove (${count} selected)`);
	});

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
			const boxes = $(".admin-list__checkbox");
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
		if (name == "") {
			$(".team-list__add-status").text("Cannot add team with empty name");
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
			$(`<tr class="admin-list__entry team-list__entry" data-id="${team.id}">
			<td>
				<input type="checkbox" class="admin-list__checkbox">
			</td>
			<td class="team-list__name">${team.name}</td>
			<td class="team-list__controls">
				<input type="button" value="Modify" class="team-list__modify">
				<input type="button" value="Remove" class="team-list__remove">
				<a href="/admin/submissions/team/${team.id}"
				class="team-list__submissions">Submissions</a>
			</td>
			<td class="team-list__members" data-members=""></td>
		</tr>`).appendTo(".team-list tbody");
			hide_team_list__add_info();
		}).catch(() => {
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

	$(".team-list__remove").on("click", event => {
		const teamId = $(event.target).parents(".team-list__entry").data("id");
		if (typeof teamId === "undefined") {
			$("#server-error-dialog__msg").text("Failed to remove team.");
			$("#server-error-dialog")[0].showModal();
			return;
		}

		event.target.disabled = true;
		fetch("/api/admin/team/remove", {
			method: "POST",
			body: JSON.stringify({id: teamId}),
			headers: {
				"Content-type": "application/json; charset=UTF-8"
			}
		}).then(res => {
			if (res.ok && res.status == 200) {
				$(event.target).parents("team-list__entry").remove();
			} else {
				$("#server-error-dialog__msg").text("Failed to remove team.");
				$("#server-error-dialog")[0].showModal();
				$("#server-error-dialog").one("close", () => {
					event.target.disabled = false;
				});
			}
		});
	});
});
