/**
 * Updates the cron expression based on the selected schedule type and parameters
 * Supports three modes: preset, simple builder, and custom
 */
function updateCronExpression() {
  var scheduleType = document.getElementById("schedule_type").value;
  var cronInput = document.getElementById("cron_expr");
  var description = document.getElementById("cron_description");

  // Show/hide relevant input sections based on schedule type
  document.getElementById("preset_options").style.display =
    scheduleType === "preset" ? "block" : "none";
  document.getElementById("simple_builder").style.display =
    scheduleType === "simple" ? "block" : "none";
  document.getElementById("custom_input").style.display =
    scheduleType === "custom" ? "block" : "none";

  var cronExpr = "";
  var desc = "";

  // Generate cron expression based on selected schedule type
  if (scheduleType === "preset") {
    // Use predefined cron expressions from dropdown
    var preset = document.getElementById("preset_select");
    cronExpr = preset.value;
    desc = preset.options[preset.selectedIndex].text;
  } else if (scheduleType === "simple") {
    // Build cron expression from simple interval/unit inputs
    var interval = document.getElementById("simple_interval").value;
    var unit = document.getElementById("simple_unit").value;
    var timePicker = document.getElementById("time_picker");

    if (unit === "minutes") {
      // Format: "0 */N * * * *" (every N minutes)
      cronExpr = "0 */" + interval + " * * * *";
      desc = "Every " + interval + " min";
      timePicker.style.display = "none";
    } else if (unit === "hours") {
      // Format: "0 0 */N * * *" (every N hours)
      cronExpr = "0 0 */" + interval + " * * *";
      desc = "Every " + interval + " hour" + (interval > 1 ? "s" : "");
      timePicker.style.display = "none";
    } else if (unit === "days") {
      // Format: "0 MM HH */N * *" (every N days at specific time)
      timePicker.style.display = "block";
      var time = document.getElementById("simple_time").value.split(":");
      cronExpr = "0 " + time[1] + " " + time[0] + " */" + interval + " * *";
      desc =
        "Every " +
        interval +
        " day" +
        (interval > 1 ? "s" : "") +
        " at " +
        document.getElementById("simple_time").value;
    }
  } else if (scheduleType === "custom") {
    // Use raw cron expression from user input
    cronExpr = document.getElementById("custom_cron").value;
    desc = "Custom";
  }

  // Update the hidden input and description display
  cronInput.value = cronExpr;
  description.textContent = desc;
}

/**
 * Initialize the form on page load
 * If editing an existing job, switch to custom mode to show the current cron expression
 */
document.addEventListener("DOMContentLoaded", function () {
  var cronInput = document.getElementById("cron_expr");
  if (cronInput.value) {
    // Existing job - show in custom mode
    document.getElementById("schedule_type").value = "custom";
    document.getElementById("custom_cron").value = cronInput.value;
    updateCronExpression();
  } else {
    // New job - show default preset
    updateCronExpression();
  }
});
