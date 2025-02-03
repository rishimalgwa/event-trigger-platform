document.addEventListener("DOMContentLoaded", function () {
    // Add event listener for the test button
    document.getElementById("testButton").addEventListener("click", function() {
        toggleTestMode();
    });
    document.getElementById("triggersList").addEventListener('click', function(event) {
        const triggerItem = event.target.closest('.trigger-item');
        if (!triggerItem) return;

        const triggerId = triggerItem.dataset.triggerId;

        if (event.target.classList.contains('update-trigger')) {
            updateTrigger(triggerId);
        }

        if (event.target.classList.contains('delete-trigger')) {
            deleteTrigger(triggerId);
        }
    });
    document.getElementById("triggerType").addEventListener("change", toggleFields);
    showTab('createTrigger');
    fetchTriggers();
});

// Show/Hide Fields Based on Trigger Type
function toggleFields() {
    let type = document.getElementById("triggerType").value;
    document.getElementById("scheduledFields").style.display = (type === "scheduled") ? "block" : "none";
    document.getElementById("apiFields").style.display = (type === "api") ? "block" : "none";
}

// Show Tabs
function showTab(tabId) {
    // Hide all tab contents
    document.querySelectorAll(".tab-content").forEach(tab => tab.style.display = "none");
    
    // Show selected tab
    document.getElementById(tabId).style.display = "block";

    // Special handling for specific tabs
    if (tabId === "allTriggers") fetchTriggers();
    if (tabId === "eventLogs") fetchLogs('active');
}

// Switch Log Type
function switchLogType(type) {
    // Update tab button styles
    document.getElementById("triggersList").addEventListener('click', function(event) {
        const triggerItem = event.target.closest('.trigger-item');
        if (!triggerItem) return;

        const triggerId = triggerItem.dataset.triggerId;

        if (event.target.classList.contains('update-trigger')) {
            updateTrigger(triggerId);
        }

        if (event.target.classList.contains('delete-trigger')) {
            deleteTrigger(triggerId);
        }
    });

    // Fetch logs of the selected type
    fetchLogs(type);
}

// Modify the HTML rendering in fetchTriggers() function
function fetchTriggers() {
    fetch("/api/trigger")
        .then(response => response.json())
        .then(data => {
            let triggers = data.body;
            let html = "";
            triggers.forEach(trigger => {
                html += `
                  <div class="trigger-item" 
                         id="trigger-${trigger.ID}" 
                         data-trigger-id="${trigger.ID}"
                         data-trigger-type="${trigger.Type}">
                        <div class="trigger-details">
                            <p><strong>ID:</strong> ${trigger.ID}</p>
                            <p><strong>Type:</strong> ${trigger.Type}</p>
                            
                            <div>
                                <label>Schedule Time:</label>
                                <input type="datetime-local" 
                                    id="scheduleTime-${trigger.ID}" 
                                    value="${formatDateTimeLocal(trigger.ScheduleTime)}"
                                    >
                            </div>

                            ${trigger.Type === 'scheduled' ? `
                                <div id="recurringFields-${trigger.ID}" style="display: ${trigger.IsRecurring ? 'block' : 'none'}">
                                    <div>
                                        <label>Interval (secs):</label>
                                        <input type="number" 
                                            id="intervalSecs-${trigger.ID}" 
                                            value="${trigger.IntervalSecs || ''}" 
                                            min="10"
                                            ${trigger.IsRecurring ? '' : 'disabled'}>
                                    </div>
                                    <div>
                                        <label>Occurrences:</label>
                                        <input type="number" 
                                            id="occurrences-${trigger.ID}" 
                                            value="${trigger.NumberOfOccurrences || ''}" 
                                            min="1"
                                            ${trigger.IsRecurring ? '' : 'disabled'}>
                                    </div>
                                    
                                </div>
                            ` : ''}

                            ${trigger.Type === 'api' ? `
                                <div>
                                    <label>API URL:</label>
                                    <input type="text" 
                                        id="apiUrl-${trigger.ID}" 
                                        value="${trigger.APIURL || ''}">
                                </div>
                                <div>
                                    <label>API Payload:</label>
                                    <textarea 
                                        id="apiPayload-${trigger.ID}"
                                        rows="4">${trigger.APIPayload ? JSON.stringify(JSON.parse(trigger.APIPayload), null, 2) : ''}</textarea>
                                </div>
                            ` : ''}
                        </div>
                       <div class="trigger-actions">
                            <button class="update-trigger">Update</button>
                            <button class="delete-trigger">Delete</button>
                        </div>
                        </div>
                    </div>
                `;
            });
            document.getElementById("triggersList").innerHTML = html;
        })
        .catch(err => {
            console.error("Error fetching triggers:", err);
            document.getElementById("triggersList").innerHTML = 
                `<div class="error">Unable to load triggers. ${err.message}</div>`;
        });
}

// Toggle Recurring Fields
function toggleRecurringFields(id) {
    const isRecurringCheckbox = document.getElementById(`isRecurring-${id}`);
    const intervalInput = document.getElementById(`intervalSecs-${id}`);
    const occurrencesInput = document.getElementById(`occurrences-${id}`);

    if (isRecurringCheckbox.checked) {
        intervalInput.disabled = false;
        occurrencesInput.disabled = false;
    } else {
        intervalInput.disabled = true;
        occurrencesInput.disabled = true;
    }
}

// Update Trigger
function updateTrigger(triggerId) {
    // Convert triggerId to string if it's not already
    const id = String(triggerId);
    
    // Find the trigger element by ID
    const triggerElement = document.getElementById(`trigger-${id}`);
    
    if (!triggerElement) {
        console.error(`Trigger element not found for ID: ${id}`);
        alert('Could not find trigger details');
        return;
    }

    // Get type from a data attribute instead of parsing text
    const type = triggerElement.dataset.triggerType;

    if (!type) {
        console.error(`Trigger type not found for ID: ${id}`);
        alert('Could not determine trigger type');
        return;
    }
    
    const body = {
        id: id, // Include ID in the payload
        
        schedule_time: new Date(document.getElementById(`scheduleTime-${id}`).value).toISOString().split('.')[0] + "Z"
    };

    if (type === 'scheduled') {
        const isRecurring = document.getElementById(`isRecurring-${id}`) != null ? document.getElementById(`isRecurring-${id}`).checked : false;
        body.is_recurring = isRecurring;

        if (isRecurring) {
            body.interval_secs = parseInt(document.getElementById(`intervalSecs-${id}`).value);
            body.number_of_occurrences = parseInt(document.getElementById(`occurrences-${id}`).value);
        }
    } else if (type === 'api') {
        body.api_url = document.getElementById(`apiUrl-${id}`).value;
        body.api_payload = document.getElementById(`apiPayload-${id}`).value;
    }

    fetch(`/api/trigger/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(body),
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to update trigger');
        }
        return response.json();
    })
    .then(() => {
        alert("Trigger updated successfully!");
        fetchTriggers();
    })
    .catch(err => {
        console.error("Error updating trigger:", err);
        alert(`Update failed: ${err.message}`);
    });
}
function formatDateTimeLocal(dateString) {
    if (!dateString) return '';

    const date = new Date(dateString);

    // Get local time components
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0'); // Months are 0-based
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');

    // Format as 'YYYY-MM-DDTHH:MM' (local time)
    return `${year}-${month}-${day}T${hours}:${minutes}`;
}
// Fetch Event Logs
function fetchLogs(type) {
    fetch(`/api/eventlog/${type}`)
        .then(response => response.json())
        .then(data => {
            let logs = data.body;
            if(!logs){
                document.getElementById("eventLogsList").innerHTML = 
                    `<div class="no-logs">No ${type} logs found.</div>`;
                return;
            }
            let html = "";
            logs.forEach(log => {
                html += `
                    <div class="log-entry">
                        <p><strong>Trigger ID:</strong> ${log.TriggerID}</p>
                        <p><strong>Status:</strong> ${log.Status}</p>
                        <p><strong>Triggered At:</strong> ${new Date(log.TriggeredAt).toLocaleString()}</p>
                        ${log.APIPayload && log.APIURL ? `
                            <p><strong>API URL:</strong> ${log.APIURL}</p>
                            <p><strong>API Payload:</strong> <pre>${JSON.stringify(JSON.parse(log.APIPayload), null, 2)}</pre></p>
                        ` : ''}
                    </div>
                `;
            });
            
            document.getElementById("eventLogsList").innerHTML = html || 
                `<div class="no-logs">No ${type} logs found.</div>`;
        })
        .catch(err => {
            console.error("Error fetching logs:", err);
            document.getElementById("eventLogsList").innerHTML = 
                `<div class="error">Unable to load ${type} logs. ${err.message}</div>`;
        });
}

// Track whether the "Test" mode is on or off
let isTestMode = false;

// Toggle "Test" mode
function toggleTestMode() {
    isTestMode = !isTestMode;
    document.getElementById("testButton").textContent = isTestMode ? "Test Mode: ON" : "Test Mode: OFF";
}

// Create Trigger function modified
async function createTrigger() {
    let type = document.getElementById("triggerType").value;
    let body = { type };

    if (type === "scheduled") {
        if (!document.getElementById("scheduleTime").value) {
            alert("Please select a schedule time.");
            return;
        }
        let scheduleTimeInput = document.getElementById("scheduleTime").value;
        body.schedule_time = scheduleTimeInput ?
            new Date(scheduleTimeInput).toISOString().split('.')[0] + "Z" : null;
        body.is_recurring = document.getElementById("isRecurring").checked;

        let interval = document.getElementById("intervalSecs").value;
        let occurrences = document.getElementById("occurrences").value;

        if (interval) body.interval_secs = parseInt(interval);
        if (occurrences) body.number_of_occurrences = parseInt(occurrences);
    } else if (type === "api") {
        if (!document.getElementById("apiScheduleTime").value) {
            alert("Please select a schedule time.");
            return;
        }
        let scheduleTimeInput = document.getElementById("apiScheduleTime").value;
        body.schedule_time = scheduleTimeInput ?
            new Date(scheduleTimeInput).toISOString().split('.')[0] + "Z" : null;
        body.api_url = document.getElementById("apiUrl").value;
        body.api_payload = document.getElementById("apiPayload").value;
    }

    // Determine the API endpoint based on Test Mode
    let endpoint = isTestMode ? "/api/trigger/test" : "/api/trigger/";

    try {
        let res = await fetch(endpoint, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(body)
        });

        let data = await res.json();
        alert(data.msg || "Trigger created successfully!");
        fetchTriggers(); // Refresh the triggers list
    } catch (err) {
        console.error("Error creating trigger:", err);
        alert("Error creating trigger!");
    }
}
// Delete Trigger
function deleteTrigger(id) {
    if (!confirm("Are you sure you want to delete this trigger?")) return;

    fetch(`/api/trigger/${id}`, {
        method: "DELETE",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ id: id })
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to delete trigger');
        }
        return response.json();
    })
    .then(() => {
        alert("Trigger deleted!");
        fetchTriggers(); // Refresh the triggers list
    })
    .catch(err => {
        console.error("Error deleting trigger:", err);
        alert(`Deletion failed: ${err.message}`);
    });
}