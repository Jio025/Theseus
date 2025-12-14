// API to fetch active containers
interface DockerContainer {
    id: string;
    name: string;
    container: string;
}

async function fetchAllActiveContainers(): Promise<DockerContainer[]> {
    try {
        // Get Request
        const response = await fetch("/api/containers/running");

        if (!response.ok) {
            throw new Error(`HTTP error, Status : ${response.status}`);
        }

        // Parsing the response
        const data: DockerContainer[] = await response.json();
        console.log(data);
        return data
    }
    catch (error){
        console.error('error fetching active docker containers:', (error as Error).message);
        throw error;
    }
}

// Function to display the active containers in html
async function displayActiveDockerContainers() {
    const activeContainersDiv = document.getElementById("activeDockerContainers")
    if (!activeContainersDiv) return;

    try {
        const containers = await fetchAllActiveContainers()
        activeContainersDiv.innerHTML = containers.map(
            (container) => `<tr>
                                <td><code>${container.id}</code></td>
                                <td>${container.name}</td>
                                <td><span class="badge rounded-pill text-bg-success">Running</span></td>
                                <td>
                                <button class="btn btn-sm btn-outline-danger">Stop</button>
                                <button class="btn btn-sm btn-outline-secondary">Logs</button>
                                </td>
                            </tr>`
        )
        .join("");
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