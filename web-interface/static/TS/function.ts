interface HostMachine {
    id: string
    ip: string
    status: string
}
interface PortBinding {
    internal: number;
    external: number;
}
interface DockerContainer {
    id?: string;
    type?: string;
    desktopEnv?: string;
    name: string;
    container: string;
    hostmachine: HostMachine;
    restartpolicy: string;
    ports: PortBinding[];
    environmentvariables: Record<string, string>;
    volumemounts: Record<string, string>;
    shmsize?: string;
    status?: string;
}

// API to fetch active containers
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
    catch (error) {
        console.error('error fetching active docker containers:', (error as Error).message);
        throw error;
    }
}

async function fetchAllActiveHostMachine(): Promise<HostMachine[]> {
    try {
        const response = await fetch("/api/hostmachine/running");

        if (!response.ok) {
            throw new Error(`HTTP error, status : ${response.status}`);
        }

        // Parsing the response
        const data: HostMachine[] = await response.json();
        return data;
    }
    catch (error) {
        console.error("error fetching active host machines:", (error as Error).message);
        throw error;
    }
}

function fetchContainerInfo(): DockerContainer {
    // 1. Fetch Basic Inputs
    const containerName = (document.getElementById("containerName") as HTMLInputElement).value;
    const imageName = (document.getElementById("imageName") as HTMLInputElement).value;
    // TODO check for host machines
    const hostMachineId = (document.getElementById("hostMachine") as HTMLSelectElement).value;
    const restartPolicy = (document.getElementById("restartPolicy") as HTMLSelectElement).value;

    // 2. Fetch Radio Button (Container Type)
    const containerType = (document.querySelector('input[name="containerType"]:checked') as HTMLInputElement).value;

    // 3. Fetch Optional Webtop Data
    const desktopEnv = (document.getElementById("desktopEnv") as HTMLSelectElement).value;

    // 4. Helper function to fetch dynamic list data (Ports, Envs, Volumes)
    const getListData = (containerId: string) => {
        const rows = document.querySelectorAll(`#${containerId} .row`);
        return Array.from(rows).map(row => {
            const inputs = row.querySelectorAll('input');
            return {
                val1: inputs[0]?.value,
                val2: inputs[1]?.value
            };
        }).filter(item => item.val1 || item.val2); // Filter out empty rows
    };

    const ports: PortBinding[] = getListData("portMappings").map(p => ({
        external: parseInt(p.val1 || "0"),
        internal: parseInt(p.val2 || "0")
    }));

    const environmentvariables: Record<string, string> = {};
    getListData("envVariables").forEach(e => {
        if (e.val1) environmentvariables[e.val1] = e.val2 || "";
    });

    const volumemounts: Record<string, string> = {};
    getListData("volumeMounts").forEach(v => {
        if (v.val1) volumemounts[v.val1] = v.val2 || "";
    });

    const finalData: DockerContainer = {
        name: containerName,
        type: containerType,
        container: containerType === 'webtop' ? `lscr.io/linuxserver/webtop:${desktopEnv}` : imageName,
        desktopEnv: containerType === 'webtop' ? desktopEnv : undefined,
        hostmachine: { id: hostMachineId, ip: "", status: "" },
        restartpolicy: restartPolicy,
        ports,
        environmentvariables,
        volumemounts,
        shmsize: containerType === 'webtop' ? "1gb" : undefined
    };

    console.log("Deployment Data:", finalData);
    return finalData;
}


// Function to create a webtop container
async function createWebtop(container: DockerContainer): Promise<void> {
    try {
        const response = await fetch("/api/webtop/create", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(container),
        });

        if (!response.ok) {
            throw new Error(`HTTP error, status : ${response.status}`);
        }
    } catch (error) {
        console.error("error creating webtop container:", (error as Error).message);
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
            if (container.status == 'running') {
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

function toggleContainerType(type: 'normal' | 'webtop') {
    const imageField = document.getElementById('imageField');
    const webtopDesktopField = document.getElementById('webtopDesktopField');
    const webtopOptions = document.getElementById('webtopOptions');

    if (type === 'webtop') {
        if (imageField) imageField.style.display = 'none';
        if (webtopDesktopField) webtopDesktopField.style.display = 'block';
        if (webtopOptions) webtopOptions.style.display = 'block';
        setupWebtopDefaults();
    } else {
        if (imageField) imageField.style.display = 'block';
        if (webtopDesktopField) webtopDesktopField.style.display = 'none';
        if (webtopOptions) webtopOptions.style.display = 'none';
    }
}

function setupWebtopDefaults() {
    // Set Container Name
    const nameInput = document.getElementById('containerName') as HTMLInputElement;
    if (nameInput && !nameInput.value) nameInput.value = 'webtop';

    // Set Ports
    const portContainer = document.getElementById('portMappings');
    if (portContainer) {
        portContainer.innerHTML = `
        <div class="row g-2 mb-2">
            <div class="col-md-5">
                <input type="number" class="form-control" value="3000">
            </div>
            <div class="col-md-5">
                <input type="number" class="form-control" value="3000">
            </div>
            <div class="col-md-2">
                <button class="btn btn-outline-danger w-100" type="button" onclick="this.closest('.row').remove()"><i class="bi bi-trash"></i></button>
            </div>
        </div>
        <div class="row g-2 mb-2">
            <div class="col-md-5">
                <input type="number" class="form-control" value="3001">
            </div>
            <div class="col-md-5">
                <input type="number" class="form-control" value="3001">
            </div>
             <div class="col-md-2">
                <button class="btn btn-outline-danger w-100" type="button" onclick="this.closest('.row').remove()"><i class="bi bi-trash"></i></button>
            </div>
        </div>`;
    }

    // Set Envs
    const envContainer = document.getElementById('envVariables');
    if (envContainer) {
        envContainer.innerHTML = `
        <div class="row g-2 mb-2">
            <div class="col-md-5"><input type="text" class="form-control" value="PUID"></div>
            <div class="col-md-5"><input type="text" class="form-control" value="1000"></div>
            <div class="col-md-2"><button class="btn btn-outline-danger w-100" onclick="this.closest('.row').remove()"><i class="bi bi-trash"></i></button></div>
        </div>
        <div class="row g-2 mb-2">
            <div class="col-md-5"><input type="text" class="form-control" value="PGID"></div>
            <div class="col-md-5"><input type="text" class="form-control" value="1000"></div>
            <div class="col-md-2"><button class="btn btn-outline-danger w-100" onclick="this.closest('.row').remove()"><i class="bi bi-trash"></i></button></div>
        </div>
        <div class="row g-2 mb-2">
            <div class="col-md-5"><input type="text" class="form-control" value="TZ"></div>
            <div class="col-md-5"><input type="text" class="form-control" value="Etc/UTC"></div>
            <div class="col-md-2"><button class="btn btn-outline-danger w-100" onclick="this.closest('.row').remove()"><i class="bi bi-trash"></i></button></div>
        </div>
    `;
    }

    // Set Restart Policy
    const restartSelect = document.getElementById('restartPolicy') as HTMLSelectElement;
    if (restartSelect) restartSelect.value = 'unless-stopped';

    // Set Volume Mount hint (Optional, just clearing or setting default)
    const volContainer = document.getElementById('volumeMounts');
    if (volContainer) {
        volContainer.innerHTML = `
        <div class="row g-2 mb-2">
            <div class="col-md-5"><input type="text" class="form-control" placeholder="/path/to/data"></div>
            <div class="col-md-5"><input type="text" class="form-control" value="/config"></div>
            <div class="col-md-2"><button class="btn btn-outline-danger w-100" onclick="this.closest('.row').remove()"><i class="bi bi-trash"></i></button></div>
        </div>`;
    }
}

document.addEventListener("DOMContentLoaded", () => {
    displayActiveDockerContainers();

    const deployBtn = document.getElementById("deployBtn");
    if (deployBtn) {
        deployBtn.addEventListener("click", async () => {
            try {
                const data = fetchContainerInfo();
                console.log("Deploying container:", data);

                await createWebtop(data);
                alert("Container deployment initiated successfully!");
                window.location.href = "/";
            } catch (error) {
                console.error("Deployment failed:", error);
                alert("Failed to deploy container: " + (error as Error).message);
            }
        });
    }

    // UI Toggling
    const normalCard = document.getElementById('normalContainerCard');
    const webtopCard = document.getElementById('webtopContainerCard');

    if (normalCard) {
        normalCard.addEventListener('click', () => {
            const radio = document.getElementById('normalContainer') as HTMLInputElement;
            if (radio) radio.checked = true;
            toggleContainerType('normal');
        });
    }

    if (webtopCard) {
        webtopCard.addEventListener('click', () => {
            const radio = document.getElementById('webtopContainer') as HTMLInputElement;
            if (radio) radio.checked = true;
            toggleContainerType('webtop');
        });
    }

    // Add dynamic row buttons handlers (for the initial rows)
    // Note: Inline onclick handlers are used in the generated HTML for simplicity
    // But we should probably add handlers for the "Add" buttons here too if they aren't working
    const addPortBtn = document.getElementById('addPortBtn');
    if (addPortBtn) {
        addPortBtn.addEventListener('click', () => {
            const container = document.getElementById('portMappings');
            if (container) {
                const div = document.createElement('div');
                div.className = 'row g-2 mb-2';
                div.innerHTML = `
                    <div class="col-md-5"><input type="number" class="form-control" placeholder="Host"></div>
                    <div class="col-md-5"><input type="number" class="form-control" placeholder="Container"></div>
                    <div class="col-md-2"><button type="button" class="btn btn-outline-danger w-100"><i class="bi bi-trash"></i></button></div>
                `;
                div.querySelector('button')?.addEventListener('click', () => div.remove());
                container.appendChild(div);
            }
        });
    }

    // Equivalent logic should be applied to Add Env and Add Volume buttons if not already present in the original code
    // The user didn't show the setup for those buttons in previous dumps, so I'll add them here for completeness
    const addEnvBtn = document.getElementById('addEnvBtn');
    if (addEnvBtn) {
        addEnvBtn.addEventListener('click', () => {
            const container = document.getElementById('envVariables');
            if (container) {
                const div = document.createElement('div');
                div.className = 'row g-2 mb-2';
                div.innerHTML = `
                    <div class="col-md-5"><input type="text" class="form-control" placeholder="Variable"></div>
                    <div class="col-md-5"><input type="text" class="form-control" placeholder="Value"></div>
                    <div class="col-md-2"><button type="button" class="btn btn-outline-danger w-100"><i class="bi bi-trash"></i></button></div>
                `;
                div.querySelector('button')?.addEventListener('click', () => div.remove());
                container.appendChild(div);
            }
        });
    }

    const addVolumeBtn = document.getElementById('addVolumeBtn');
    if (addVolumeBtn) {
        addVolumeBtn.addEventListener('click', () => {
            const container = document.getElementById('volumeMounts');
            if (container) {
                const div = document.createElement('div');
                div.className = 'row g-2 mb-2';
                div.innerHTML = `
                     <div class="col-md-5"><input type="text" class="form-control" placeholder="Host Path"></div>
                                <div class="col-md-5"><input type="text" class="form-control" placeholder="Container Path"></div>
                                <div class="col-md-2">
                                    <button class="btn btn-outline-danger w-100" type="button">
                                        <i class="bi bi-trash"></i>
                                    </button>
                                </div>
                `;
                div.querySelector('button')?.addEventListener('click', () => div.remove());
                container.appendChild(div);
            }
        });
    }
});