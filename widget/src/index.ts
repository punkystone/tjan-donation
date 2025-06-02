window.addEventListener("onWidgetLoad", function (data) {
    const format = data.detail.fieldData["format"];
    if (format === undefined) return;
    const months = [
        "Tjanuar",
        "Tjebruar",
        "Tjärz",
        "Tjapril",
        "Tjai",
        "Tjuni",
        "Tjuli",
        "Tjaugust",
        "Tjeptember",
        "Tjoktober",
        "Tjovember",
        "Tjecember",
    ];
    const monthElement = document.querySelector(".month");
    if (monthElement === null) return;
    const updateLoop = (): void => {
        const now = new Date();
        const monthName = months[now.getMonth()];
        if (monthName === undefined) return;
        monthElement.textContent = format.replace("{month}", monthName);
    };
    updateLoop();
    setInterval(updateLoop, 1000);
});
