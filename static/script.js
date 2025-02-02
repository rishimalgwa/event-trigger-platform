document.getElementById("triggerType").addEventListener("change", function() {
    let type = this.value;
    document.getElementById("scheduledFields").style.display = type === "scheduled" ? "block" : "none";
    document.getElementById("apiFields").style.display = type === "api" ? "block" : "none";
});

// Function to create trigger
async function createTrigger() {
    let type = document.getElementById("triggerType").value;
    let body = { type };
    if (type === "scheduled") {
        let scheduleTimeInput = document.getElementById("scheduleTime").value;
        let scheduleTimeRFC3339 = scheduleTimeInput ? new Date(scheduleTimeInput).toISOString().split('.')[0] + "Z" : null;
        body.schedule_time =scheduleTimeRFC3339;
        body.is_recurring = document.getElementById("isRecurring").checked ;
        let interval = document.getElementById("intervalSecs").value;
        let occurrences = document.getElementById("occurrences").value;
        if (interval) body.interval_secs = parseInt(interval) ;
        if (occurrences) body.number_of_occurrences = parseInt(occurrences);
    } else if (type === "api") {
        let scheduleTimeInput = document.getElementById("apiScheduleTime").value;
        let scheduleTimeRFC3339 = scheduleTimeInput ? new Date(scheduleTimeInput).toISOString().split('.')[0] + "Z" : null;
        body.api_url = document.getElementById("apiUrl").value;
        body.api_payload = JSON.stringify(document.getElementById("apiPayload").value);
        body.schedule_time =scheduleTimeRFC3339;
    }

    try {
        let res = await fetch("/api/trigger/", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(body)
        });

        let data = await res.json();
        alert(data.msg || "Trigger created successfully!");
    } catch (err) {
        alert("Error creating trigger!",err);
    }
}

async function fetchLogs(type) {
    try {
        let res = await fetch(`/api/eventlog/${type}`);

        if (!res.ok) {
            throw new Error(`HTTP error! Status: ${res.status}`);
        }

        let data = await res.json(); // Fails if response is not valid JSON
        let logContainer = document.getElementById("eventLogs");
        logContainer.innerHTML = "";
        if (data.body === null || data.body.length === 0) {
            logContainer.innerHTML = "<p>No logs found.</p>";
            return;
        }

        let logs = data.body;
        logs.forEach(log => {
            let div = document.createElement("div");
            div.classList.add("log-entry");
            div.innerHTML = `<strong>TriggerID:</strong> ${log.TriggerID} <br> 
                            <strong>TriggeredAt:</strong> ${log.TriggeredAt} <br> 
                            <strong>Status:</strong> ${log.Status} <br> 
                            `;
            if(log.APIPayload && log.APIURL){ 
                div.innerHTML += `<strong>API Payload:</strong> ${log.APIPayload} <br>`;
                div.innerHTML += `<strong>API URL:</strong> ${log.APIURL} <br>`;
            }
            logContainer.appendChild(div);
        });
    } catch (err) {
        console.error("Error fetching logs:", err);
        alert("Error fetching logs! Check console for details.");
    }
}