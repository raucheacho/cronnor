function updateCronExpression() {
  var scheduleType = document.getElementById("schedule_type").value;
  var cronInput = document.getElementById("cron_expr");
  var description = document.getElementById("cron_description");

  document.getElementById("preset_options").style.display =
    scheduleType === "preset" ? "block" : "none";
  document.getElementById("simple_builder").style.display =
    scheduleType === "simple" ? "block" : "none";
  document.getElementById("custom_input").style.display =
    scheduleType === "custom" ? "block" : "none";

  var cronExpr = "";
  var desc = "";

  if (scheduleType === "preset") {
    var preset = document.getElementById("preset_select");
    cronExpr = preset.value;
    desc = preset.options[preset.selectedIndex].text;
  } else if (scheduleType === "simple") {
    var interval = document.getElementById("simple_interval").value;
    var unit = document.getElementById("simple_unit").value;
    var timePicker = document.getElementById("time_picker");

    if (unit === "minutes") {
      cronExpr = "0 */" + interval + " * * * *";
      desc = "Every " + interval + " min";
      timePicker.style.display = "none";
    } else if (unit === "hours") {
      cronExpr = "0 0 */" + interval + " * * *";
      desc = "Every " + interval + " hour" + (interval > 1 ? "s" : "");
      timePicker.style.display = "none";
    } else if (unit === "days") {
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
    cronExpr = document.getElementById("custom_cron").value;
    desc = "Custom";
  }

  cronInput.value = cronExpr;
  description.textContent = desc;
}

document.addEventListener("DOMContentLoaded", function () {
  var cronInput = document.getElementById("cron_expr");
  if (cronInput.value) {
    document.getElementById("schedule_type").value = "custom";
    document.getElementById("custom_cron").value = cronInput.value;
    updateCronExpression();
  } else {
    updateCronExpression();
  }
});
