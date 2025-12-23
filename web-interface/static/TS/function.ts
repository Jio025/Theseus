// API to fetch active containers
interface DockerContainer {
    id: string;
    name: string;
    container: string;
    status: string;
}
interface ActiveHostMachine {
    id: string
    ip: string
    status: string
}

async function fetchAllActiveContainers(): Promise<DockerContainer[]> {
    try {
        // Get Request
        const response = await fetch("/api/containers/running");

        if (!response.ok) {
            throw new Error(`HTTP error, status : ${response.status}`);
        }

        // Parsing the response
        const data: DockerContainer[] = await response.json();
        return data
    }
    catch (error){
        console.error('error fetching active docker containers:', (error as Error).message);
        throw error;
    }
}

async function fetchAllActiveHostMachine(): Promise<ActiveHostMachine[]> {
    try {
        const response = await fetch("/api/hostmachine/running");

        if(!response.ok){
            throw new Error(`HTTP error, status : ${response.status}`);
        }

        // Parsing the response
        const data: ActiveHostMachine[] = await response.json();
        return data;
    }
    catch(error) {
        console.error("error fetching active host machines:", (error as Error).message);
        throw error;
    }
}

// Function to display the active containers in html
async function displayActiveDockerContainers() {
    const activeContainersDiv = document.getElementById("activeDockerContainers")
    if (!activeContainersDiv) return;
    const activeContainersCountDiv = document.getElementById("activeContainerCount")
    if (!activeContainersCountDiv) return;

    try {
        const containers = await fetchAllActiveContainers()
        let activeContainerCount: number = 0;
        containers.forEach((container) => {
            if(container.status == 'running'){
                activeContainerCount++;
            }
        });
        activeContainersCountDiv.textContent = `${activeContainerCount} Active`
        activeContainersDiv.innerHTML = containers.map(
            (container) => {
                let badgeClass, statusLabel, actionButton;

                switch (container.status) {
                    case 'running':
                        badgeClass = 'text-bg-success';
                        statusLabel = 'Active';
                        actionButton = `<button class="btn btn-sm btn-outline-danger">Stop</button>`;
                        break;
                    case 'restarting':
                        badgeClass = 'text-bg-warning';
                        statusLabel = 'Restarting';
                        actionButton = ``;
                        break;
                    // Stopped is default
                    default:
                        badgeClass = 'text-bg-danger';
                        statusLabel = 'Stopped';
                        actionButton = `<button class="btn btn-sm btn-outline-success">Start</button>`;
                }
                return `
                        <tr>
                            <td><code>${container.id}</code></td>
                            <td>${container.name}</td>
                            <td><span class="badge rounded-pill ${badgeClass}">${statusLabel}</span></td>
                            <td>
                                ${actionButton}
                                <button class="btn btn-sm btn-outline-secondary">Logs</button>
                            </td>
                        </tr>`;
    }
).join("");
} catch (error) {
        activeContainersDiv.innerHTML = `
            <tr>
                <td colspan="4" class="text-danger">
                    Failed to load containers
                </td>
            </tr>
        `;
    }
}

document.addEventListener("DOMContentLoaded", () => {
    displayActiveDockerContainers();
});