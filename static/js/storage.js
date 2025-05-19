const initialFiles = [
    { name: "АНАПА_2007.zip", color: "bg-blue-500" },
    { name: "ВАЖНЫЙ_ФАЙЛ.docx", color: "bg-blue-900" },
    { name: "что-то.txt", color: "bg-gray-500" },
    { name: "Презентация.ptpx", color: "bg-orange-400" },
    { name: "Таблица.xlsx", color: "bg-green-500" }
];

let deletedFiles = JSON.parse(localStorage.getItem("deletedFiles")) || [];

function saveTrash() {
    localStorage.setItem("deletedFiles", JSON.stringify(deletedFiles));
}

function isTrashPage() {
    return window.location.pathname.includes("trash.html");
}

function renderFiles() {
    const container = document.getElementById(isTrashPage() ? "trash-list" : "file-list");
    if (!container) return;
    container.innerHTML = "";

    const filesToShow = isTrashPage()
        ? initialFiles.filter(f => deletedFiles.includes(f.name))
        : initialFiles.filter(f => !deletedFiles.includes(f.name));

    filesToShow.forEach(file => {
        const wrapper = document.createElement("div");
        wrapper.className = "flex justify-between items-center bg-white p-3 rounded border shadow";

        const left = document.createElement("div");
        left.className = "flex items-center gap-3";
        left.innerHTML = `<div class="w-3 h-3 ${file.color} rounded-sm"></div><span class="font-semibold">${file.name}</span>`;

        const right = document.createElement("div");
        right.className = "flex gap-2";

        if (isTrashPage()) {
            const restoreBtn = document.createElement("button");
            restoreBtn.textContent = "↩ Восстановить";
            restoreBtn.className = "bg-green-500 text-white px-3 py-1 rounded shadow text-sm";
            restoreBtn.onclick = () => {
                deletedFiles = deletedFiles.filter(f => f !== file.name);
                saveTrash();
                renderFiles();
            };
            right.appendChild(restoreBtn);
        } else {
            const downloadBtn = document.createElement("button");
            downloadBtn.textContent = "Загрузить";
            downloadBtn.className = "bg-white border px-3 py-1 rounded shadow text-sm";

            const deleteBtn = document.createElement("button");
            deleteBtn.textContent = "🗑 Удалить";
            deleteBtn.className = "bg-red-400 text-white px-3 py-1 rounded shadow text-sm";
            deleteBtn.onclick = () => {
                showConfirm(file.name, () => {
                    deletedFiles.push(file.name);
                    saveTrash();
                    renderFiles();
                });
            };

            right.append(downloadBtn, deleteBtn);
        }

        wrapper.append(left, right);
        container.appendChild(wrapper);
    });
}

function showConfirm(fileName, onConfirm) {
    const modal = document.getElementById("confirm-modal");
    const text = document.getElementById("confirm-text");
    const yes = document.getElementById("confirm-yes");
    const no = document.getElementById("confirm-no");

    if (!modal) return;

    text.textContent = `Вы точно хотите удалить файл «${fileName}»?`;
    modal.classList.remove("hidden");

    // Удаляем старые обработчики
    const newYes = yes.cloneNode(true);
    const newNo = no.cloneNode(true);
    yes.parentNode.replaceChild(newYes, yes);
    no.parentNode.replaceChild(newNo, no);

    newYes.onclick = () => {
        modal.classList.add("hidden");
        onConfirm();
    };

    newNo.onclick = () => {
        modal.classList.add("hidden");
    };
}

document.addEventListener("DOMContentLoaded", renderFiles);
